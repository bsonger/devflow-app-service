package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Manifest struct {
	BaseModel     `bson:",inline"`
	ApplicationId primitive.ObjectID `json:"application_id" bson:"application_id"`
	Name          string             `json:"name" bson:"name"`
}

func (Manifest) CollectionName() string { return "manifests" }
