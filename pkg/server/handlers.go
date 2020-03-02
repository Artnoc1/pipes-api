package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/tambet/oauthplain"

	"github.com/toggl/pipes-api/pkg/authorization"
	"github.com/toggl/pipes-api/pkg/environment"
	"github.com/toggl/pipes-api/pkg/integrations"
	"github.com/toggl/pipes-api/pkg/toggl"
)

// mutex to prevent multiple of postPipeRun on same workspace run at same time
var postPipeRunWorkspaceLock = map[int]*sync.Mutex{}
var postPipeRunLock sync.Mutex

type Controller struct {
	stResolver ServiceTypeResolver
	ptResolver PipeTypeResolver

	pipesSvc   *integrations.Service
	pipesStore *integrations.Storage
	authStore  *authorization.Storage
	env        *environment.Environment
	api        *toggl.ApiClient
}

func NewController(env *environment.Environment, pipes *integrations.Service, pipesStore *integrations.Storage, authStore *authorization.Storage, api *toggl.ApiClient) *Controller {
	return &Controller{
		pipesSvc:   pipes,
		pipesStore: pipesStore,
		authStore:  authStore,
		env:        env,
		api:        api,

		stResolver: pipes,
		ptResolver: pipes,
	}
}

func (c *Controller) GetIntegrations(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	resp, err := c.pipesSvc.WorkspaceIntegrations(workspaceID)
	if err != nil {
		return internalServerError(err.Error())
	}
	return ok(resp)
}

func (c *Controller) GetIntegrationPipe(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	pipeID := mux.Vars(req.r)["pipe"]
	if !c.ptResolver.AvailablePipeType(pipeID) {
		return badRequest("Missing or invalid pipe")
	}

	pipe, err := c.pipesStore.LoadPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	if pipe == nil {
		pipe = integrations.NewPipe(workspaceID, serviceID, pipeID)
	}

	pipe.PipeStatus, err = c.pipesStore.LoadPipeStatus(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}

	return ok(pipe)
}

func (c *Controller) PostPipeSetup(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	pipeID := mux.Vars(req.r)["pipe"]
	if !c.ptResolver.AvailablePipeType(pipeID) {
		return badRequest("Missing or invalid pipe")
	}

	pipe := integrations.NewPipe(workspaceID, serviceID, pipeID)
	errorMsg := pipe.ValidateServiceConfig(req.body)
	if errorMsg != "" {
		return badRequest(errorMsg)
	}

	if err := c.pipesStore.Save(pipe); err != nil {
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) PutPipeSetup(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	pipeID := mux.Vars(req.r)["pipe"]
	if !c.ptResolver.AvailablePipeType(pipeID) {
		return badRequest("Missing or invalid pipe")
	}
	if len(req.body) == 0 {
		return badRequest("Missing payload")
	}
	pipe, err := c.pipesStore.LoadPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	if pipe == nil {
		return badRequest("Pipe is not configured")
	}
	if err := json.Unmarshal(req.body, &pipe); err != nil {
		return internalServerError(err.Error())
	}
	if err := c.pipesStore.Save(pipe); err != nil {
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) DeletePipeSetup(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	pipeID := mux.Vars(req.r)["pipe"]
	if !c.ptResolver.AvailablePipeType(pipeID) {
		return badRequest("Missing or invalid pipe")
	}
	pipe, err := c.pipesStore.LoadPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	if pipe == nil {
		return badRequest("Pipe is not configured")
	}
	if err := c.pipesStore.Destroy(pipe, workspaceID); err != nil {
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) GetAuthURL(req Request) Response {
	serviceID := mux.Vars(req.r)["service"]
	accountName := req.r.FormValue("account_name")
	callbackURL := req.r.FormValue("callback_url")

	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	if accountName == "" {
		return badRequest("Missing or invalid account_name")
	}
	if callbackURL == "" {
		return badRequest("Missing or invalid callback_url")
	}
	config, found := c.env.GetOAuth1Configs(serviceID)
	if !found {
		return badRequest("env OAuth config not found")
	}
	transport := &oauthplain.Transport{
		Config: config.UpdateURLs(accountName),
	}
	token, err := transport.AuthCodeURL(callbackURL)
	if err != nil {
		return internalServerError(err.Error())
	}
	return ok(struct {
		AuthURL string `json:"auth_url"`
	}{
		token.AuthorizeUrl,
	})
}

func (c *Controller) PostAuthorization(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	if len(req.body) == 0 {
		return badRequest("Missing payload")
	}

	var payload map[string]interface{}
	err := json.Unmarshal(req.body, &payload)
	if err != nil {
		return internalServerError(err.Error())
	}

	auth := authorization.New(workspaceID, serviceID)
	auth.WorkspaceToken = currentWorkspaceToken(req.r)

	switch c.authStore.GetAvailableAuthorizations(serviceID) {
	case authorization.TypeOauth1:
		auth.Data, err = c.env.OAuth1Exchange(serviceID, payload)
	case authorization.TypeOauth2:
		auth.Data, err = c.env.OAuth2Exchange(serviceID, payload)
	}
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := c.authStore.Save(auth); err != nil {
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) DeleteAuthorization(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	service := integrations.Create(serviceID, workspaceID)
	auth, err := c.authStore.LoadAuth(service.GetWorkspaceID(), service.Name())
	if err != nil {
		return internalServerError(err.Error())
	}
	if err := service.SetAuthData(auth.Data); err != nil {
		return internalServerError(err.Error())
	}

	if err := c.authStore.Destroy(service.GetWorkspaceID(), service.Name()); err != nil {
		return internalServerError(err.Error())
	}
	if err := c.pipesStore.DeletePipeByWorkspaceIDServiceID(workspaceID, serviceID); err != nil {
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) GetServiceAccounts(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	service := integrations.Create(serviceID, workspaceID)
	auth, err := c.authStore.LoadAuth(service.GetWorkspaceID(), service.Name())
	if err != nil {
		return badRequest("No authorizations for " + serviceID)
	}
	if err := service.SetAuthData(auth.Data); err != nil {
		return internalServerError(err.Error())
	}

	if err := c.authStore.Refresh(auth); err != nil {
		return badRequest("oAuth refresh failed!")
	}
	forceImport := req.r.FormValue("force")
	if forceImport == "true" {
		if err := c.pipesStore.ClearImportFor(service, "accounts"); err != nil {
			return internalServerError(err.Error())
		}
	}
	accountsResponse, err := c.pipesStore.GetAccounts(service)
	if err != nil {
		return internalServerError("Unable to get accounts from DB")
	}
	if accountsResponse == nil {
		go func() {
			if err := c.pipesStore.FetchAccounts(service); err != nil {
				log.Print(err.Error())
			}
		}()
		return noContent()
	}
	return ok(accountsResponse)
}

func (c *Controller) GetServiceUsers(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)

	serviceID := mux.Vars(req.r)["service"]
	if !c.stResolver.AvailableServiceType(serviceID) {
		return badRequest("Missing or invalid service")
	}
	service := integrations.Create(serviceID, workspaceID)
	auth, err := c.authStore.LoadAuth(service.GetWorkspaceID(), service.Name())
	if err != nil {
		return badRequest("No authorizations for " + serviceID)
	}
	if err := service.SetAuthData(auth.Data); err != nil {
		return internalServerError(err.Error())
	}

	pipeID := "users"
	pipe, err := c.pipesStore.LoadPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	if pipe == nil {
		return badRequest("Pipe is not configured")
	}
	if err := service.SetParams(pipe.ServiceParams); err != nil {
		return badRequest(err.Error())
	}

	forceImport := req.r.FormValue("force")
	if forceImport == "true" {
		if err := c.pipesStore.ClearImportFor(service, pipeID); err != nil {
			return internalServerError(err.Error())
		}
	}

	usersResponse, err := c.pipesSvc.GetUsers(service)
	if err != nil {
		return internalServerError("Unable to get users from DB")
	}
	if usersResponse == nil {
		if forceImport == "true" {
			go func() {
				if err := c.pipesSvc.FetchObjects(pipe, false); err != nil {
					log.Print(err.Error())
				}
			}()
		}
		return noContent()
	}
	return ok(usersResponse)
}

func (c *Controller) GetServicePipeLog(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID := currentServicePipeID(req.r)

	pipeStatus, err := c.pipesStore.LoadPipeStatus(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError("Unable to get log from DB")
	}
	if pipeStatus == nil {
		return noContent()
	}
	return Response{http.StatusOK, pipeStatus.GenerateLog(), "text/plain"}
}

func (c *Controller) PostServicePipeClearConnections(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID := currentServicePipeID(req.r)

	pipe, err := c.pipesStore.LoadPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	if pipe == nil {
		return badRequest("Pipe is not configured")
	}

	err = c.pipesSvc.ClearPipeConnections(pipe)
	if err != nil {
		return internalServerError("Unable to get clear connections")
	}
	return noContent()
}

func (c *Controller) PostPipeRun(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)

	// make sure no race condition on fetching workspace lock
	postPipeRunLock.Lock()
	workspaceLock, exists := postPipeRunWorkspaceLock[workspaceID]
	if !exists {
		workspaceLock = &sync.Mutex{}
		postPipeRunWorkspaceLock[workspaceID] = workspaceLock
	}
	postPipeRunLock.Unlock()

	serviceID, pipeID := currentServicePipeID(req.r)

	pipe, err := c.pipesStore.LoadPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	if pipe == nil {
		return badRequest("Pipe is not configured")
	}
	if msg := pipe.ValidatePayload(req.body); msg != "" {
		return badRequest(msg)
	}
	if pipe.ID == "users" {
		go func() {
			workspaceLock.Lock()
			c.pipesSvc.Run(pipe)
			workspaceLock.Unlock()
		}()
		time.Sleep(500 * time.Millisecond) // TODO: Is that synchronization ? :D
	} else {
		if err := c.pipesSvc.QueuePipeAsFirst(pipe); err != nil {
			return internalServerError(err.Error())
		}
	}
	return ok(nil)
}

func (c *Controller) GetStatus(req Request) Response {
	resp := &struct {
		Reasons []string `json:"reasons"`
	}{}

	if c.pipesStore.IsDown() {
		resp := &struct {
			Reasons []string `json:"reasons"`
		}{
			[]string{"Database is down"},
		}
		return serviceUnavailable(resp)
	}

	if err := c.api.PingTogglApi(); err != nil {
		resp.Reasons = append(resp.Reasons, err.Error())
	}

	if len(resp.Reasons) > 0 {
		return serviceUnavailable(resp)
	}
	return ok(map[string]string{"status": "OK"})
}
