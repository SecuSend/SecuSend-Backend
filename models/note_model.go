package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	Id                primitive.ObjectID `json:"id,omitempty"`
	Key               string             `json:"key,omitempty" validate:"required"`
	Data              string             `json:"data,omitempty" validate:"required"`
	PasswordProtected bool               `json:"passwordProtected,omitempty"`
	SelfDestruct      bool               `json:"selfDestruct,omitempty"`
	ExpireAt          time.Time          `json:"expireAt,omitempty"`
	CreatedAt         time.Time          `json:"createdAt,omitempty" validate:"required"`
}
