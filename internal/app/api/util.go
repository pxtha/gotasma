package api

import (
	"sync"

	"github.com/globalsign/mgo"
	"gopkg.in/olivere/elastic.v5"

	"github.com/gotasma/internal/pkg/db/elasticsearch"
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

func newElasticSearchClient() (*elastic.Client, error) {
	esConf := elasticsearch.LoadConfigFromEnv()
	es, err := elasticsearch.NewClient(esConf)
	if err != nil {
		return nil, err
	}
	return es, nil
}
