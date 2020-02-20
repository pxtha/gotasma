package api

import "github.com/gotasma/internal/app/task"

func newTaskService(policy task.PolicyService) (*task.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := task.NewMongoDBRespository(s)
	return task.New(repo, policy), nil
}
