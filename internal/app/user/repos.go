package user

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

func (r *MongoDBRepository) Create(ctx context.Context, user *types.User) (string, error) {
	s := r.session.Clone()
	defer s.Close()
	user.CreatedAt = time.Now()
	user.UpdateAt = user.CreatedAt

	if err := r.collection(s).Insert(user); err != nil {
		return "", err
	}
	return user.UserID, nil
}

func (r *MongoDBRepository) Delete(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Remove(bson.M{"user_id": id})
}

func (r *MongoDBRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	selector := bson.M{"email": email}
	s := r.session.Clone()
	defer s.Close()
	var user *types.User
	if err := r.collection(s).Find(selector).One(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MongoDBRepository) FindAllDev(ctx context.Context, createrID string) ([]*types.User, error) {
	selector := bson.M{"creater_id": createrID}
	s := r.session.Clone()
	defer s.Close()
	var users []*types.User
	if err := r.collection(s).Find(selector).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}
func (r *MongoDBRepository) FindDevsByID(ctx context.Context, userIDs []string) ([]*types.User, error) {
	selector := bson.M{"user_id": bson.M{
		"$in": userIDs,
	}}
	s := r.session.Clone()
	defer s.Close()
	var users []*types.User
	if err := r.collection(s).Find(selector).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoDBRepository) FindByID(ctx context.Context, UserID string) (*types.User, error) {
	selector := bson.M{"user_id": UserID}
	s := r.session.Clone()
	defer s.Close()
	var users *types.User
	if err := r.collection(s).Find(selector).One(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoDBRepository) UpdateProjectsID(ctx context.Context, userID string, projectID string, addToSet bool) error {

	s := r.session.Clone()
	defer s.Clone()
	//add data to array if not exist
	action := "$addToSet"
	//pull data out of array
	if !addToSet {
		action = "$pull"
	}

	return r.collection(s).Update(bson.M{"user_id": userID}, bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
		action: bson.M{
			"projects_id": projectID,
		},
	},
	)
}
func (r *MongoDBRepository) UpdateTasksID(ctx context.Context, userID string, taskID string, addToSet bool) error {

	s := r.session.Clone()
	defer s.Clone()
	//add data to array if not exist
	action := "$addToSet"
	//pull data out of array
	if !addToSet {
		action = "$pull"
	}
	return r.collection(s).Update(bson.M{"user_id": userID}, bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
		action: bson.M{
			"tasks_id": taskID,
		},
	},
	)
}

func (r *MongoDBRepository) FindByProjectID(ctx context.Context, projectID string) ([]*types.User, error) {

	selector := bson.M{"projects_id": projectID}
	s := r.session.Clone()
	defer s.Close()
	var users []*types.User
	if err := r.collection(s).Find(selector).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("user")
}
