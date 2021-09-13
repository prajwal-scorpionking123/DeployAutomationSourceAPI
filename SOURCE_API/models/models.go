package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Source struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DeveloperName string             `json:"devname,omitempty" bson:"devname,omitempty"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty"`
	SourceLink    string             `json:"sourcelink,omitempty" bson:"sourcelink,omitempty"`
	Ticket        string             `json:"ticket,omitempty" bson:"ticket,omitempty"`
	Timestamp     time.Time          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	IsApproved    bool               `json:"isapproved,omitempty"  bson:"isapproved,omitempty"`
	IsRequested   bool               `json:"isrequested" bson:"isrequested"`
	IsVarified    bool               `json:"isvarified,omitempty"  bson:"isvarified,omitempty"`
}
type ValueAttr struct {
	Value bool `json:"value"`
}
type Deployment struct {
	SourceLink      string `json:"sourcelink,omitempty" bson:"sourcelink,omitempty"`
	DestinationLink string `json:"destinationlink,omitempty" bson:"destinationlink,omitempty"`
}
type Login struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type UserObject struct {
	FirstName    string
	LastName     string
	MobileNumber string
	Email        string
	Username     string
}
type MetaData struct {
	Source   string   `json:"source,omitempty"`
	Filepath []string `json:"filepath,omitempty"`
}
