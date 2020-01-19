package api

import "praslar.com/gotasma/internal/app/user"

func newUserService() (*user.Service, error) {
	s, err := dialDefaultMongoDB()
	if err != nil {
		return nil, err
	}

	repo := user.NewMongoDBRespository(s)
	return user.New(repo), nil
}

func newUserHandler(srv *user.Service) *user.Handler {
	return user.NewHandler(srv)
}
