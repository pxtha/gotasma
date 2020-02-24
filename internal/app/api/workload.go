package api

import (
	"github.com/gotasma/internal/app/workload"
)

func newWorkloadService(policy workload.PolicyService) (*workload.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := workload.NewMongoDBRespository(s)
	return workload.New(repo, policy), nil
}
