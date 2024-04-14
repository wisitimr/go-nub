package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID   `bson:"_id" json:"id"`
	Username  string               `bson:"username" json:"username"`
	Password  string               `bson:"password" json:"password"`
	FullName  string               `bson:"fullName" json:"fullName"`
	FirstName string               `bson:"firstName" json:"firstName"`
	LastName  string               `bson:"lastName" json:"lastName"`
	Email     string               `bson:"email" json:"email"`
	Companies []primitive.ObjectID `bson:"companies" json:"companies"`
	Role      string               `bson:"role" json:"role"`
	CreatedBy primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID   `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type UserCompany struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Email     string             `bson:"email" json:"email"`
	Companies []Company          `bson:"companies" json:"companies"`
	Role      string             `bson:"role" json:"role"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type UserProfile struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Username    string             `bson:"username" json:"username"`
	FullName    string             `bson:"fullName" json:"fullName"`
	FirstName   string             `bson:"firstName" json:"firstName"`
	LastName    string             `bson:"lastName" json:"lastName"`
	Email       string             `bson:"email" json:"email"`
	Role        string             `bson:"role" json:"role"`
	AccessToken string             `bson:"accessToken" json:"accessToken"`
	Companies   []Company          `bson:"companies" json:"companies"`
}

type UpdatedUserProfile struct {
	Id        primitive.ObjectID   `bson:"_id" json:"id"`
	Username  string               `bson:"username" json:"username"`
	FullName  string               `bson:"fullName" json:"fullName"`
	FirstName string               `bson:"firstName" json:"firstName"`
	LastName  string               `bson:"lastName" json:"lastName"`
	Email     string               `bson:"email" json:"email"`
	Companies []primitive.ObjectID `bson:"companies" json:"companies"`
	Role      string               `bson:"role" json:"role"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"Password"`
}

type Company struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Address     string             `bson:"address" json:"address"`
	Phone       string             `bson:"phone" json:"phone"`
	Contact     string             `bson:"contact" json:"contact"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
