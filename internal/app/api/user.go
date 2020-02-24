package api

import "github.com/gotasma/internal/app/user"

func newUserService(policy user.PolicyService, workload user.WorkLoad) (*user.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := user.NewMongoDBRespository(s)
	return user.New(repo, policy, workload), nil
}

func newUserHandler(srv *user.Service) *user.Handler {
	return user.NewHandler(srv)
}
