package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/toggl/pipes-api/pkg/integrations"
	"github.com/toggl/pipes-api/pkg/pipe"
)

type IntegrationsFileStorage struct {
	availablePipeType     *regexp.Regexp
	availableServiceType  *regexp.Regexp
	availableIntegrations []*pipe.Integration
	// Stores available authorization types for each service
	// Map format: map[externalServiceID]authType
	availableAuthTypes map[integrations.ExternalServiceID]string
	mx                 sync.RWMutex
}

func NewIntegrationsFileStorage(integrationsConfigPath string) *IntegrationsFileStorage {
	svc := &IntegrationsFileStorage{
		availableIntegrations: []*pipe.Integration{},
		availableAuthTypes:    map[integrations.ExternalServiceID]string{},
	}
	svc.loadIntegrations(integrationsConfigPath).fillAvailableServices().fillAvailablePipeTypes()
	svc.mx.RLock()
	for _, integration := range svc.availableIntegrations {
		svc.availableAuthTypes[integration.ID] = integration.AuthType
	}
	svc.mx.RUnlock()
	return svc
}

func (fis *IntegrationsFileStorage) IsValidPipe(pipeID integrations.PipeID) bool {
	return fis.availablePipeType.MatchString(string(pipeID))
}

func (fis *IntegrationsFileStorage) IsValidService(serviceID integrations.ExternalServiceID) bool {
	return fis.availableServiceType.MatchString(string(serviceID))
}

func (fis *IntegrationsFileStorage) LoadAuthorizationType(serviceID integrations.ExternalServiceID) (string, error) {
	fis.mx.RLock()
	defer fis.mx.RUnlock()
	return fis.availableAuthTypes[serviceID], nil
}

func (fis *IntegrationsFileStorage) LoadIntegrations() ([]*pipe.Integration, error) {
	fis.mx.RLock()
	defer fis.mx.RUnlock()
	return fis.availableIntegrations, nil
}

func (fis *IntegrationsFileStorage) SaveAuthorizationType(serviceID integrations.ExternalServiceID, authType string) error {
	fis.mx.Lock()
	defer fis.mx.Unlock()
	fis.availableAuthTypes[serviceID] = authType
	return nil
}

func (fis *IntegrationsFileStorage) loadIntegrations(integrationsConfigPath string) *IntegrationsFileStorage {
	fis.mx.Lock()
	defer fis.mx.Unlock()
	b, err := ioutil.ReadFile(integrationsConfigPath)
	if err != nil {
		log.Fatalf("Could not read integrations.json, reason: %v", err)
	}
	if err := json.Unmarshal(b, &fis.availableIntegrations); err != nil {
		log.Fatalf("Could not parse integrations.json, reason: %v", err)
	}
	return fis
}

func (fis *IntegrationsFileStorage) fillAvailableServices() *IntegrationsFileStorage {
	ids := make([]string, len(fis.availableIntegrations))
	for i := range fis.availableIntegrations {
		ids = append(ids, string(fis.availableIntegrations[i].ID))
	}
	fis.availableServiceType = regexp.MustCompile(strings.Join(ids, "|"))
	return fis
}

func (fis *IntegrationsFileStorage) fillAvailablePipeTypes() *IntegrationsFileStorage {
	fis.mx.Lock()
	defer fis.mx.Unlock()
	str := fmt.Sprintf("%s|%s|%s|%s|%s|%s", integrations.UsersPipe, integrations.ProjectsPipe, integrations.TodoListsPipe, integrations.TodosPipe, integrations.TasksPipe, integrations.TimeEntriesPipe)
	fis.availablePipeType = regexp.MustCompile(str)
	return fis
}