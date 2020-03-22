package storage

import (
	"database/sql"
	"encoding/json"

	"github.com/bugsnag/bugsnag-go"

	"github.com/toggl/pipes-api/pkg/integrations"
	"github.com/toggl/pipes-api/pkg/toggl"
)

const (
	loadImportsSQL = `
	SELECT data FROM imports
	WHERE workspace_id = $1 AND Key = $2
	ORDER by created_at DESC
	LIMIT 1
	`
	saveImportsSQL = `
	INSERT INTO imports(workspace_id, Key, data, created_at)
    VALUES($1, $2, $3, NOW())
	`

	clearImportsSQL = `
	    DELETE FROM imports
	    WHERE workspace_id = $1 AND Key = $2
	`

	truncateImportsSQL = `TRUNCATE TABLE imports`
)

type ImportsPostgresStorage struct {
	db *sql.DB
}

func NewImportsPostgresStorage(db *sql.DB) *ImportsPostgresStorage {
	return &ImportsPostgresStorage{db: db}
}

func (pis *ImportsPostgresStorage) SaveAccountsFor(s integrations.ExternalService, res toggl.AccountsResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}
	return pis.saveObject(s, integrations.AccountsPipe, b)
}

func (pis *ImportsPostgresStorage) SaveUsersFor(s integrations.ExternalService, res toggl.UsersResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	return pis.saveObject(s, integrations.UsersPipe, b)
}

func (pis *ImportsPostgresStorage) SaveClientsFor(s integrations.ExternalService, res toggl.ClientsResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	return pis.saveObject(s, integrations.ClientsPipe, b)
}

func (pis *ImportsPostgresStorage) SaveProjectsFor(s integrations.ExternalService, res toggl.ProjectsResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	return pis.saveObject(s, integrations.ProjectsPipe, b)
}

func (pis *ImportsPostgresStorage) SaveTasksFor(s integrations.ExternalService, res toggl.TasksResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	return pis.saveObject(s, integrations.TasksPipe, b)
}

func (pis *ImportsPostgresStorage) SaveTodoListsFor(s integrations.ExternalService, res toggl.TasksResponse) error {
	b, err := json.Marshal(res)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}

	return pis.saveObject(s, integrations.TodoListsPipe, b)
}

func (pis *ImportsPostgresStorage) LoadAccountsFor(s integrations.ExternalService) (*toggl.AccountsResponse, error) {
	b, err := pis.loadObject(s, integrations.AccountsPipe)
	if err != nil || b == nil {
		return nil, err
	}

	var accountsResponse toggl.AccountsResponse
	err = json.Unmarshal(b, &accountsResponse)
	if err != nil {
		return nil, err
	}
	return &accountsResponse, nil
}

func (pis *ImportsPostgresStorage) LoadUsersFor(s integrations.ExternalService) (*toggl.UsersResponse, error) {
	b, err := pis.loadObject(s, integrations.UsersPipe)
	if err != nil || b == nil {
		return nil, err
	}

	var usersResponse toggl.UsersResponse
	err = json.Unmarshal(b, &usersResponse)
	if err != nil {
		return nil, err
	}
	return &usersResponse, nil
}

func (pis *ImportsPostgresStorage) LoadClientsFor(s integrations.ExternalService) (*toggl.ClientsResponse, error) {
	b, err := pis.loadObject(s, integrations.ClientsPipe)
	if err != nil || b == nil {
		return nil, err
	}

	var clientsResponse toggl.ClientsResponse
	err = json.Unmarshal(b, &clientsResponse)
	if err != nil {
		return nil, err
	}
	return &clientsResponse, nil
}

func (pis *ImportsPostgresStorage) LoadProjectsFor(s integrations.ExternalService) (*toggl.ProjectsResponse, error) {
	b, err := pis.loadObject(s, integrations.ProjectsPipe)
	if err != nil || b == nil {
		return nil, err
	}

	var projectsResponse toggl.ProjectsResponse
	err = json.Unmarshal(b, &projectsResponse)
	if err != nil {
		return nil, err
	}

	return &projectsResponse, nil
}

func (pis *ImportsPostgresStorage) LoadTodoListsFor(s integrations.ExternalService) (*toggl.TasksResponse, error) {
	b, err := pis.loadObject(s, integrations.TodoListsPipe)
	if err != nil || b == nil {
		return nil, err
	}

	var tasksResponse toggl.TasksResponse
	err = json.Unmarshal(b, &tasksResponse)
	if err != nil {
		return nil, err
	}
	return &tasksResponse, nil
}

func (pis *ImportsPostgresStorage) LoadTasksFor(s integrations.ExternalService) (*toggl.TasksResponse, error) {
	b, err := pis.loadObject(s, integrations.TasksPipe)
	if err != nil || b == nil {
		return nil, err
	}

	var tasksResponse toggl.TasksResponse
	err = json.Unmarshal(b, &tasksResponse)
	if err != nil {
		return nil, err
	}
	return &tasksResponse, nil
}

func (pis *ImportsPostgresStorage) DeleteAccountsFor(s integrations.ExternalService) error {
	_, err := pis.db.Exec(clearImportsSQL, s.GetWorkspaceID(), s.KeyFor(integrations.AccountsPipe))
	return err
}

func (pis *ImportsPostgresStorage) DeleteUsersFor(s integrations.ExternalService) error {
	_, err := pis.db.Exec(clearImportsSQL, s.GetWorkspaceID(), s.KeyFor(integrations.UsersPipe))
	return err
}

func (pis *ImportsPostgresStorage) loadObject(s integrations.ExternalService, pid integrations.PipeID) ([]byte, error) {
	var result []byte
	rows, err := pis.db.Query(loadImportsSQL, s.GetWorkspaceID(), s.KeyFor(pid))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	if err := rows.Scan(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (pis *ImportsPostgresStorage) saveObject(s integrations.ExternalService, pid integrations.PipeID, b []byte) error {
	_, err := pis.db.Exec(saveImportsSQL, s.GetWorkspaceID(), s.KeyFor(pid), b)
	if err != nil {
		bugsnag.Notify(err)
		return err
	}
	return nil
}