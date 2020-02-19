package api

import (
	"github.com/gotasma/internal/app/project"
	"github.com/sirupsen/logrus"
)

func newProjectService(policy project.PolicyService, holiday project.HolidayService) (*project.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		logrus.Errorf("Fail to dial defautl mongodb API util")
		return nil, err
	}

	repo := project.NewMongoDBRespository(s)

	es, err := newElasticSearchClient()
	if err != nil {
		logrus.Errorf("Fail to new defautl elastic search API util")
		return nil, err
	}

	elastic := project.NewElasticSearchRepository(es)

	return project.New(repo, policy, elastic, holiday), nil
}

func newProjectHandler(srv *project.Service) *project.Handler {
	return project.NewHandler(srv)
}
