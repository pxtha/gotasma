package api

import "praslar.com/gotasma/internal/app/project"

func newProjectService(policy project.PolicyService) (*project.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := project.NewMongoDBRespository(s)

	return project.New(repo, policy), nil
}

func newProjectHandler(srv *project.Service) *project.Handler {
	return project.NewHandler(srv)
}
