package task

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gotasma/internal/app/types"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoDBRespository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) FindByProjectID(ctx context.Context, projectID string) ([]*types.Task, error) {

	selector := bson.M{"project_id": projectID}
	s := r.session.Clone()
	defer s.Close()
	var tasks []*types.Task
	if err := r.collection(s).Find(selector).All(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *MongoDBRepository) Create(ctx context.Context, task *types.Task) error {

	s := r.session.Clone()
	defer s.Clone()

	task.UpdateAt = time.Now()

	return r.collection(s).Insert(task)

}

func (r *MongoDBRepository) Update(ctx context.Context, projectID string, req *types.Task) error {
	s := r.session.Clone()
	defer s.Close()

	return r.collection(s).Update(bson.M{"task_id": req.TaskID}, bson.M{
		"$set": bson.M{
			"updated_at":        time.Now(),
			"label":             req.Label,
			"type":              req.Type,
			"parent":            req.Parent,
			"effort":            req.Effort,
			"start":             req.Start,
			"duration":          req.Duration,
			"estimate_duration": req.EstimateDuration,
			"end":               req.End,
			"parents":           req.Parents,
			"children":          req.Children,
			"all_children":      req.AllChildren,
			"project_id":        projectID,
		},
	},
	)

}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("task")
}
