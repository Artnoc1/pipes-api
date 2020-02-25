package environment

import (
	"errors"
	"fmt"
	"time"
)

type PipeConfig struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description,omitempty"`
	Automatic       bool              `json:"automatic,omitempty"`
	AutomaticOption bool              `json:"automatic_option"`
	Configured      bool              `json:"configured"`
	Premium         bool              `json:"premium"`
	ServiceParams   []byte            `json:"service_params,omitempty"`
	PipeStatus      *PipeStatusConfig `json:"pipe_status,omitempty"`

	WorkspaceID int        `json:"-"`
	ServiceID   string     `json:"-"`
	Key         string     `json:"-"`
	Payload     []byte     `json:"-"`
	LastSync    *time.Time `json:"-"`
}

func NewPipe(workspaceID int, serviceID, pipeID string) *PipeConfig {
	return &PipeConfig{
		ID:          pipeID,
		Key:         PipesKey(serviceID, pipeID),
		ServiceID:   serviceID,
		WorkspaceID: workspaceID,
	}
}

func (p *PipeConfig) ValidateServiceConfig(payload []byte) string {
	service := Create(p.ServiceID, p.WorkspaceID)
	err := service.SetParams(payload)
	if err != nil {
		return err.Error()
	}
	p.ServiceParams = payload
	return ""
}

func (p *PipeConfig) ValidatePayload(payload []byte) string {
	if p.ID == "users" && len(payload) == 0 {
		return "Missing request payload"
	}
	p.Payload = payload
	return ""
}

func PipesKey(serviceID, pipeID string) string {
	return fmt.Sprintf("%s:%s", serviceID, pipeID)
}

var (
	// ErrJSONParsing hides json marshalling errors from users
	ErrJSONParsing = errors.New("Failed to parse response from service, please contact support")
)
