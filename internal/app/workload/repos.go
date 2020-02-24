package workload

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

func (r *MongoDBRepository) FindByID(ctx context.Context, projectID string, userID string) (*types.WorkLoad, error) {

	selector := bson.M{"user_id": userID, "project_id": projectID}
	s := r.session.Clone()
	defer s.Close()
	var workload *types.WorkLoad
	if err := r.collection(s).Find(selector).One(&workload); err != nil {
		return nil, err
	}
	return workload, nil
}

func (r *MongoDBRepository) Create(ctx context.Context, workload *types.WorkLoad) error {

	s := r.session.Clone()
	defer s.Clone()

	workload.CreatedAt = time.Now()
	workload.UpdateAt = workload.CreatedAt

	return r.collection(s).Insert(workload)
}

func (r *MongoDBRepository) Update(ctx context.Context, projectID string, userID string, overload map[int]int) error {
	s := r.session.Clone()
	defer s.Close()

	return r.collection(s).Update(bson.M{"project_id": projectID, "user_id": userID}, bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
			"overload":   overload,
		},
	},
	)
}

func (r *MongoDBRepository) Delete(ctx context.Context, projectID string, userID string) error {
	s := r.session.Clone()
	defer s.Close()
	if projectID == "_all_projects_" {
		_, err := r.collection(s).RemoveAll(bson.M{"user_id": userID})
		if err != nil {
			return err
		}
		return nil
	}
	if userID == "_all_devs_" {
		_, err := r.collection(s).RemoveAll(bson.M{"project_id": projectID})
		if err != nil {
			return err
		}
		return nil
	}
	return r.collection(s).Remove(bson.M{"project_id": projectID, "user_id": userID})
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("workload")
}
