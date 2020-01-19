package router

import (
	"net/http"

	"github.com/gorilla/mux"

	envconfig "praslar.com/gotasma/internal/pkg/env"
)

type (
	// Config hold configurations of router
	Config struct {

		// need to set manually
		Routes          []Route
		NotFoundHandler http.Handler
	}

	// Route hold configuration of routing
	Route struct {
		Desc    string
		Path    string
		Method  string
		Queries []string
		Handler http.HandlerFunc
	}
)

func New(conf *Config) (http.Handler, error) {

	r := mux.NewRouter()
	for _, rt := range conf.Routes {
		var h http.Handler
		h = http.HandlerFunc(rt.Handler)
		r.Path(rt.Path).Methods(rt.Method).Handler(h).Queries(rt.Queries...)
	}

	if conf.NotFoundHandler != nil {
		r.NotFoundHandler = conf.NotFoundHandler
	}

	return r, nil
}

func LoadConfigFromEnv() *Config {
	var conf Config
	envconfig.Load(&conf)
	return &conf
}
