package project

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

func (r *MongoDBRepository) FindByName(ctx context.Context, name string, createrID string) (*types.Project, error) {
	selector := bson.M{"name": name, "creater_id": createrID}
	s := r.session.Clone()
	defer s.Close()
	var project *types.Project
	if err := r.collection(s).Find(selector).One(&project); err != nil {
		return nil, err
	}
	return project, nil
}
func (r *MongoDBRepository) FindByProjectID(ctx context.Context, projectID string) (*types.Project, error) {
	selector := bson.M{"project_id": projectID}
	s := r.session.Clone()
	defer s.Close()
	var project *types.Project
	if err := r.collection(s).Find(selector).One(&project); err != nil {
		return nil, err
	}
	return project, nil
}

func (r *MongoDBRepository) FindAllByUserID(ctx context.Context, id string, role types.Role) ([]*types.Project, error) {
	searchBy := "devs_id"
	if role == types.PM {
		searchBy = "creater_id"
	}
	selector := bson.M{searchBy: id}
	s := r.session.Clone()
	defer s.Close()
	var project []*types.Project
	if err := r.collection(s).Find(selector).All(&project); err != nil {
		return nil, err
	}
	return project, nil
}

func (r *MongoDBRepository) Create(ctx context.Context, project *types.Project) error {
	s := r.session.Clone()
	defer s.Clone()

	project.CreatedAt = time.Now()
	project.UpdateAt = project.CreatedAt

	return r.collection(s).Insert(project)

}

func (r *MongoDBRepository) Delete(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()
	if err := r.collection(s).Remove(bson.M{"project_id": id}); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) UpdateDevsID(ctx context.Context, devsID []string, projectID string, addToSet bool) error {

	s := r.session.Clone()
	defer s.Clone()
	//add data to array if not exist
	action := "$addToSet"
	data := bson.M{
		"devs_id": bson.M{
			"$each": devsID,
		},
	}
	//pull data out of array
	if !addToSet {
		action = "$pull"
		data = bson.M{
			"devs_id": devsID[0],
		}
	}

	return r.collection(s).Update(bson.M{"project_id": projectID}, bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
		action: data,
	},
	)
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("project")
}
