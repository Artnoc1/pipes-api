package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/toggl/pipes-api/pkg/integrations"
	"github.com/toggl/pipes-api/pkg/pipe/service"
)

type Controller struct {
	stResolver ServiceTypeResolver
	ptResolver PipeTypeResolver

	pipesSvc *service.Service
}

func NewController(pipes *service.Service) *Controller {
	return &Controller{
		pipesSvc: pipes,

		stResolver: pipes,
		ptResolver: pipes,
	}
}

func (c *Controller) GetIntegrationsHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	resp, err := c.pipesSvc.WorkspaceIntegrations(workspaceID)
	if err != nil {
		return internalServerError(err.Error())
	}
	return ok(resp)
}

func (c *Controller) GetIntegrationPipeHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID, err := c.getIntegrationParams(req)
	if err != nil {
		return badRequest(err.Error())
	}
	p, err := c.pipesSvc.GetIntegrationPipe(workspaceID, serviceID, pipeID)
	if err != nil {
		return internalServerError(err.Error())
	}
	return ok(p)
}

func (c *Controller) PostPipeSetupHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID, err := c.getIntegrationParams(req)
	if err != nil {
		return badRequest(err.Error())
	}
	err = c.pipesSvc.CreatePipe(workspaceID, serviceID, pipeID, req.body)
	if err != nil {
		if errors.As(err, &service.SetParamsError{}) {
			return badRequest(err.Error())
		}
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) PutPipeSetupHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID, err := c.getIntegrationParams(req)
	if err != nil {
		return badRequest(err.Error())
	}
	if len(req.body) == 0 {
		return badRequest("Missing payload")
	}

	err = c.pipesSvc.UpdatePipe(workspaceID, serviceID, pipeID, req.body)
	if err != nil {
		if errors.Is(err, service.ErrPipeNotConfigured) {
			return badRequest(err.Error())
		}
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) DeletePipeSetupHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID, err := c.getIntegrationParams(req)
	if err != nil {
		return badRequest(err.Error())
	}
	err = c.pipesSvc.DeletePipe(workspaceID, serviceID, pipeID)
	if err != nil {
		if errors.Is(err, service.ErrPipeNotConfigured) {
			return badRequest(err.Error())
		}
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) GetAuthURLHandler(req Request) Response {
	serviceID, err := c.getServiceId(req)
	if err != nil {
		return badRequest(err.Error())
	}
	accountName := req.r.FormValue("account_name")
	if accountName == "" {
		return badRequest("Missing or invalid account_name")
	}
	callbackURL := req.r.FormValue("callback_url")
	if callbackURL == "" {
		return badRequest("Missing or invalid callback_url")
	}

	url, err := c.pipesSvc.GetAuthURL(serviceID, accountName, callbackURL)
	if err != nil {
		if errors.Is(err, &service.LoadError{}) {
			return badRequest(err.Error())
		}
		return internalServerError(err.Error())
	}
	return ok(struct {
		AuthURL string `json:"auth_url"`
	}{url})
}

func (c *Controller) PostAuthorizationHandler(req Request) Response {
	currentToken := currentWorkspaceToken(req.r)
	workspaceID := currentWorkspaceID(req.r)
	serviceID, err := c.getServiceId(req)
	if err != nil {
		return badRequest(err.Error())
	}
	if len(req.body) == 0 {
		return badRequest("Missing payload")
	}

	err = c.pipesSvc.CreateAuthorization(workspaceID, serviceID, currentToken, req.body)
	if err != nil {
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) DeleteAuthorizationHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, err := c.getServiceId(req)
	if err != nil {
		return badRequest(err.Error())
	}
	err = c.pipesSvc.DeleteAuthorization(workspaceID, serviceID)
	if err != nil {
		return internalServerError(err.Error())
	}

	return ok(nil)
}

func (c *Controller) GetServiceAccountsHandler(req Request) Response {
	forceImport := req.r.FormValue("force")
	workspaceID := currentWorkspaceID(req.r)
	serviceID, err := c.getServiceId(req)
	if err != nil {
		return badRequest(err.Error())
	}
	fi, err := strconv.ParseBool(forceImport)
	if err != nil {
		return badRequest(err.Error())
	}

	accountsResponse, err := c.pipesSvc.GetServiceAccounts(workspaceID, serviceID, fi)
	if err != nil {
		if errors.Is(err, &service.LoadError{}) {
			return badRequest(err.Error())
		}
		if errors.Is(err, &service.RefreshError{}) {
			return badRequest(err.Error())
		}

		return internalServerError(err.Error())
	}

	return ok(accountsResponse)
}

func (c *Controller) GetServiceUsersHandler(req Request) Response {
	forceImport := req.r.FormValue("force")
	workspaceID := currentWorkspaceID(req.r)
	serviceID, err := c.getServiceId(req)
	if err != nil {
		return badRequest(err.Error())
	}

	fi, err := strconv.ParseBool(forceImport)
	if err != nil {
		return badRequest(err.Error())
	}

	usersResponse, err := c.pipesSvc.GetServiceUsers(workspaceID, serviceID, fi)
	if err != nil {
		if errors.Is(err, service.ErrNoContent) {
			return noContent()
		}

		if errors.Is(err, &service.LoadError{}) {
			return badRequest("No authorizations for " + serviceID)
		}

		if errors.Is(err, service.ErrPipeNotConfigured) {
			return badRequest(err.Error())
		}

		if errors.Is(err, &service.SetParamsError{}) {
			return badRequest(err.Error())
		}
		return internalServerError(err.Error())
	}
	return ok(usersResponse)
}

func (c *Controller) GetServicePipeLogHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID := currentServicePipeID(req.r)
	pipesLog, err := c.pipesSvc.GetServicePipeLog(workspaceID, serviceID, pipeID)
	if err != nil {
		if errors.Is(err, service.ErrNoContent) {
			return noContent()
		}
		return internalServerError("Unable to get log from DB")
	}
	return Response{http.StatusOK, pipesLog, "text/plain"}
}

func (c *Controller) PostServicePipeClearConnectionsHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID := currentServicePipeID(req.r)

	err := c.pipesSvc.ClearPipeConnections(workspaceID, serviceID, pipeID)
	if err != nil {
		if errors.Is(err, service.ErrPipeNotConfigured) {
			return badRequest(err.Error())
		}
		return internalServerError("Unable to get clear connections: " + err.Error())
	}

	return noContent()
}

func (c *Controller) PostPipeRunHandler(req Request) Response {
	workspaceID := currentWorkspaceID(req.r)
	serviceID, pipeID := currentServicePipeID(req.r)

	err := c.pipesSvc.RunPipe(workspaceID, serviceID, pipeID, req.body)
	if err != nil {
		if errors.Is(err, service.ErrPipeNotConfigured) {
			return badRequest(err.Error())
		}
		if errors.Is(err, &service.SetParamsError{}) {
			return badRequest(err.Error())
		}
		return internalServerError(err.Error())
	}
	return ok(nil)
}

func (c *Controller) GetStatusHandler(Request) Response {
	resp := &struct {
		Reasons []string `json:"reasons"`
	}{}

	errs := c.pipesSvc.Ready()
	for _, err := range errs {
		resp.Reasons = append(resp.Reasons, err.Error())
	}

	if len(resp.Reasons) > 0 {
		return serviceUnavailable(resp)
	}
	return ok(map[string]string{"status": "OK"})
}

func (c *Controller) getIntegrationParams(req Request) (integrations.ExternalServiceID, integrations.PipeID, error) {
	serviceID := integrations.ExternalServiceID(mux.Vars(req.r)["service"])
	if !c.stResolver.AvailableServiceType(serviceID) {
		return "", "", errors.New("missing or invalid service")
	}
	pipeID := integrations.PipeID(mux.Vars(req.r)["pipe"])
	if !c.ptResolver.AvailablePipeType(pipeID) {
		return "", "", errors.New("Missing or invalid pipe")
	}
	return serviceID, pipeID, nil
}

func (c *Controller) getServiceId(req Request) (integrations.ExternalServiceID, error) {
	serviceID := integrations.ExternalServiceID(mux.Vars(req.r)["service"])
	if !c.stResolver.AvailableServiceType(serviceID) {
		return "", errors.New("missing or invalid service")
	}
	return serviceID, nil
}
