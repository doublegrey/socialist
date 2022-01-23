package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	Surname   string             `bson:"surname" json:"surname"`
	Status    string             `bson:"status" json:"status"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	Photos    []string           `bson:"photos" json:"photos"`
	Followed  bool               `bson:"followed" json:"followed"`
	Following []string           `bson:"following" json:"following"`
	Password  string             `bson:"password" json:"-"`
}
