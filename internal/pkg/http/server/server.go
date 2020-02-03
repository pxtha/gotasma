package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func ListenAndServe(conf Config, router http.Handler) {
	port := fmt.Sprint(conf.Port)
	if conf.Port == 0 {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}
	address := fmt.Sprintf("%s:%s", conf.Address, port)
	srv := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadTimeout:       conf.ReadTimeout,
		WriteTimeout:      conf.WriteTimeout,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
	}
	logrus.Infof("HTTP Server is listening on: %s", address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Panicf("listen: %s\n", err)
	}

}
