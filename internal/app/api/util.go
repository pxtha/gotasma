package api

import (
	"sync"

	"github.com/globalsign/mgo"

	"github.com/gotasma/internal/pkg/db/mongodb"
)

var (
	session     *mgo.Session
	sessionOnce sync.Once
)

func dialDefaultMongoDB() (*mgo.Session, error) {
	repoConf := mongodb.LoadConfigFromEnv()
	var err error
	sessionOnce.Do(func() {
		session, err = mongodb.Dial(repoConf)
	})
	if err != nil {
		return nil, err
	}
	s := session.Clone()
	return s, nil
}
