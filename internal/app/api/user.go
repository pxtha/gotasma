package api

import "github.com/gotasma/internal/app/user"

func newUserService(policy user.PolicyService, project user.ProjectService) (*user.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := user.NewMongoDBRespository(s)
	return user.New(repo, policy, project), nil
}

func newUserHandler(srv *user.Service) *user.Handler {
	return user.NewHandler(srv)
}
