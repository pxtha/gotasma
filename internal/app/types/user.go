package types

import "time"

const (
	PM Role = iota
	DEV
	TEMP
	ADMIN
)

type (
	UserStatus string
	Role       int
	User       struct {
		ID        string     `json:"id,omitempty" bson:"_id,omitempty"`
		Email     string     `json:"email,omitempty" bson:"email,omitempty"`
		Password  string     `json:"-" bson:"password,omitempty"`
		FirstName string     `json:"first_name,omitempty" bson:"first_name,omitempty"`
		LastName  string     `json:"last_name,omitempty" bson:"last_name,omitempty"`
		UserID    string     `json:"user_id,omitempty" bson:"user_id,omitempty"`
		Role      Role       `json:"role,omitempty" bson:"role,omitempty"`
		CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdateAt  *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	RegisterRequest struct {
		Email     string `json:"email,omitempty" validate:"required,email"`
		Password  string `json:"password" validate:"required,gt=3"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Role      Role   `json:"role,omitempty" validate:"required,lt=3" bson:"role,omitempty"`
	}
)

func (user *User) Strip() *User {
	stripedUser := User(*user)
	stripedUser.Password = ""
	return &stripedUser
}
