package elasticsearch

import (
	envconfig "github.com/gotasma/internal/pkg/env"
	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

type (
	Config struct {
		URL []string `envconfig:"ELASTIC_URL" default:"http://elasticsearch:9200"`
	}
)

func LoadConfigFromEnv() *Config {
	var conf Config
	envconfig.Load(&conf)
	return &conf
}

func NewClient(conf *Config) (*elastic.Client, error) {
	logrus.Info(conf.URL)
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(conf.URL...),
	)
	if err != nil {
		logrus.Errorf("Fail to Create new elastic Client from env: err: %v ", err)
		return nil, err
	}
	esversion, err := client.ElasticsearchVersion(conf.URL[0])
	if err != nil {
		// Handle error
		logrus.Warning("Fail to connect elastic: err: %v ", err)
	}
	logrus.Infof("successfully connect to ES at %v and Elasticsearch version %s\n", conf.URL, esversion)
	return client, nil

}
