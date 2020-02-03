package main

import (
	"github.com/sirupsen/logrus"

	"github.com/gotasma/internal/app/api"
	"github.com/gotasma/internal/pkg/http/server"
)

func main() {

	router, err := api.NewRouter()
	if err != nil {
		logrus.Panic("Cannot init Router, err: ", err)
	}
	severConf := server.LoadConfigFromEnv()
	server.ListenAndServe(severConf, router)
}
