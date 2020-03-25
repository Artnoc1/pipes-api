package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"

	"github.com/toggl/pipes-api/pkg/integration"
	"github.com/toggl/pipes-api/pkg/toggl"
)

type Service struct {
	WorkspaceID int
	token       oauth.Token
}

func (s *Service) ID() integration.ID {
	return integration.GitHub
}

func (s *Service) GetWorkspaceID() int {
	return s.WorkspaceID
}

func (s *Service) KeyFor(objectType integration.PipeID) string {
	return fmt.Sprintf("github:%s", objectType)
}

func (s *Service) SetAuthData(b []byte) error {
	if err := json.Unmarshal(b, &s.token); err != nil {
		return err
	}
	return nil
}

func (s *Service) Accounts() ([]*toggl.Account, error) {
	var accounts []*toggl.Account
	account := toggl.Account{ID: 1, Name: "Self"}
	accounts = append(accounts, &account)
	return accounts, nil
}

// Map Github repos to projects
func (s *Service) Projects() ([]*toggl.Project, error) {
	repos, _, err := s.client().Repositories.List(context.Background(), "", nil)
	if err != nil {
		return nil, err
	}
	var projects []*toggl.Project
	for _, object := range repos {
		project := toggl.Project{
			Active:    true,
			Name:      *object.Name,
			ForeignID: strconv.FormatInt(*object.ID, 10),
		}
		projects = append(projects, &project)
	}
	return projects, nil
}

func (s *Service) SetSince(*time.Time) {}

func (s *Service) SetParams([]byte) error {
	return nil
}

func (s *Service) Users() ([]*toggl.User, error) {
	return []*toggl.User{}, nil
}

func (s *Service) Clients() ([]*toggl.Client, error) {
	return []*toggl.Client{}, nil
}

func (s *Service) Tasks() ([]*toggl.Task, error) {
	return []*toggl.Task{}, nil
}

func (s *Service) TodoLists() ([]*toggl.Task, error) {
	return []*toggl.Task{}, nil
}

func (s *Service) ExportTimeEntry(*toggl.TimeEntry) (int, error) {
	return 0, nil
}

func (s *Service) client() *github.Client {
	t := &oauth.Transport{Token: &s.token}
	return github.NewClient(t.Client())
}