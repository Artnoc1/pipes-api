package storage

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/toggl/pipes-api/pkg/integration"
	"github.com/toggl/pipes-api/pkg/pipe"
)

// AuthorizationsPostgresStorage SQL queries
const (
	selectAuthorizationSQL = `SELECT workspace_id, service, workspace_token, data
		FROM authorizations
		WHERE workspace_id = $1
		AND service = $2
		LIMIT 1
  `
	insertAuthorizationSQL = `WITH existing_auth AS (
		UPDATE authorizations SET data = $4, workspace_token = $3
		WHERE workspace_id = $1 AND service = $2
		RETURNING service
	),
	inserted_auth AS (
		INSERT INTO
		authorizations(workspace_id, service, workspace_token, data)
		SELECT $1, $2, $3, $4
		WHERE NOT EXISTS (SELECT 1 FROM existing_auth)
		RETURNING service
	)
	SELECT * FROM inserted_auth
	UNION
	SELECT * FROM existing_auth
  `
	deleteAuthorizationSQL   = `DELETE FROM authorizations WHERE workspace_id = $1 AND service = $2`
	truncateAuthorizationSQL = `TRUNCATE TABLE authorizations`
)

type AuthorizationsPostgresStorage struct {
	db *sql.DB
}

func NewAuthorizationsPostgresStorage(db *sql.DB) *AuthorizationsPostgresStorage {
	return &AuthorizationsPostgresStorage{db: db}
}

func (ps *AuthorizationsPostgresStorage) SaveAuthorization(a *pipe.Authorization) error {
	_, err := ps.db.Exec(insertAuthorizationSQL, a.WorkspaceID, a.ServiceID, a.WorkspaceToken, a.Data)
	if err != nil {
		return err
	}
	return nil
}

func (ps *AuthorizationsPostgresStorage) LoadAuthorization(workspaceID int, externalServiceID integration.ID, a *pipe.Authorization) error {
	rows, err := ps.db.Query(selectAuthorizationSQL, workspaceID, externalServiceID)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return rows.Err()
	}
	err = rows.Scan(&a.WorkspaceID, &a.ServiceID, &a.WorkspaceToken, &a.Data)
	if err != nil {
		return err
	}
	return nil
}

func (ps *AuthorizationsPostgresStorage) DeleteAuthorization(workspaceID int, externalServiceID integration.ID) error {
	_, err := ps.db.Exec(deleteAuthorizationSQL, workspaceID, externalServiceID)
	return err
}

// LoadWorkspaceAuthorizations loads map with authorizations status for each externalService.
// Map format: map[externalServiceID]isAuthorized
func (ps *AuthorizationsPostgresStorage) LoadWorkspaceAuthorizations(workspaceID int) (map[integration.ID]bool, error) {
	authorizations := make(map[integration.ID]bool)
	rows, err := ps.db.Query(`SELECT service FROM authorizations WHERE workspace_id = $1`, workspaceID)
	if err != nil {
		return authorizations, err
	}
	defer rows.Close()
	for rows.Next() {
		var service integration.ID
		if err := rows.Scan(&service); err != nil {
			return authorizations, err
		}
		authorizations[service] = true
	}
	return authorizations, nil
}