package pipe

import (
	"github.com/toggl/pipes-api/pkg/integrations"
	"github.com/toggl/pipes-api/pkg/toggl"
)

//go:generate mockery -name Storage -case underscore -inpkg
type Storage interface {
	Queue

	IsDown() bool
	QueuePipeAsFirst(pipe *Pipe) error

	GetAccounts(s integrations.ExternalService) (*toggl.AccountsResponse, error)
	SaveAccounts(s integrations.ExternalService) error
	ClearImportFor(s integrations.ExternalService, pid integrations.PipeID) error

	LoadPipe(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID) (*Pipe, error)
	LoadPipeStatus(workspaceID int, sid integrations.ExternalServiceID, pid integrations.PipeID) (*Status, error)
	LoadAuthorization(workspaceID int, sid integrations.ExternalServiceID) (*Authorization, error)
	LoadConnection(workspaceID int, key string) (*Connection, error)
	LoadReversedConnection(workspaceID int, key string) (*ReversedConnection, error)
	LoadPipes(workspaceID int) (map[string]*Pipe, error)
	LoadLastSync(p *Pipe)
	LoadPipeStatuses(workspaceID int) (map[string]*Status, error)
	LoadWorkspaceAuthorizations(workspaceID int) (map[integrations.ExternalServiceID]bool, error)

	DeletePipeByWorkspaceIDServiceID(workspaceID int, sid integrations.ExternalServiceID) error
	DeletePipeConnections(workspaceID int, pipeConnectionKey, pipeStatusKey string) (err error)

	Destroy(p *Pipe, workspaceID int) error
	DestroyAuthorization(workspaceID int, externalServiceID integrations.ExternalServiceID) error

	Save(p *Pipe) error
	SaveConnection(c *Connection) error
	SavePipeStatus(p *Status) error
	SaveAuthorization(a *Authorization) error

	GetObject(s integrations.ExternalService, pid integrations.PipeID) ([]byte, error)
	SaveObject(workspaceID int, objKey string, obj interface{}) error
}
