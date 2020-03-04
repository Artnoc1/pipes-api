package integrations

import (
	"fmt"
	"time"
)

type Pipe struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	Automatic       bool    `json:"automatic,omitempty"`
	AutomaticOption bool    `json:"automatic_option"`
	Configured      bool    `json:"configured"`
	Premium         bool    `json:"premium"`
	ServiceParams   []byte  `json:"service_params,omitempty"`
	PipeStatus      *Status `json:"pipe_status,omitempty"`

	WorkspaceID int        `json:"-"`
	ServiceID   string     `json:"-"`
	Key         string     `json:"-"`
	Payload     []byte     `json:"-"`
	LastSync    *time.Time `json:"-"`
}

func NewPipe(workspaceID int, serviceID, pipeID string) *Pipe {
	return &Pipe{
		ID:          pipeID,
		Key:         PipesKey(serviceID, pipeID),
		ServiceID:   serviceID,
		WorkspaceID: workspaceID,
	}
}

func (p *Pipe) ValidatePayload(payload []byte) string {
	if p.ID == "users" && len(payload) == 0 {
		return "Missing request payload"
	}
	p.Payload = payload
	return ""
}

func PipesKey(serviceID, pipeID string) string {
	return fmt.Sprintf("%s:%s", serviceID, pipeID)
}
