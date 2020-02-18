package project

import (
	"context"
	"errors"

	"github.com/gotasma/internal/app/types"

	"github.com/sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

type (
	ElasticSearchRepository struct {
		client *elastic.Client
	}
)

func NewElasticSearchRepository(client *elastic.Client) *ElasticSearchRepository {
	return &ElasticSearchRepository{
		client: client,
	}
}

func (es *ElasticSearchRepository) IndexNewHistory(ctx context.Context, project *types.ProjectHistory) error {

	exists, err := es.client.IndexExists("history").Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, err := es.client.CreateIndex("history").BodyJson(types.Mapping).Do(ctx)
		if err != nil {
			logrus.Errorf("Fail to create index: err %v", err)
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("Cannot create index")
		}
	}

	put1, err := es.client.Index().
		Index("history").
		Type("_doc").
		BodyJson(project).
		Do(ctx)

	if err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Infof("Indexed history %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	return nil
}
