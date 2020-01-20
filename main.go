package main

import (
	"github.com/sirupsen/logrus"

	"praslar.com/gotasma/internal/app/api"
	"praslar.com/gotasma/internal/pkg/http/server"
)

func main() {

	// env := flag.String("env", "", "env file")
	// flag.Parse()
	// if *env != "" {
	// 	if err := envconfig.SetEnvFromFile(*env); err != nil {
	// 		logrus.Panicf("failed to set env, err: %v", err)
	// 	}
	// }
	// logrus.Infof("initializing HTTP routing...")

	router, err := api.NewRouter()
	if err != nil {
		logrus.Panic("Cannot init Router, err: ", err)
	}
	severConf := server.LoadConfigFromEnv()
	server.ListenAndServe(severConf, router)
}
