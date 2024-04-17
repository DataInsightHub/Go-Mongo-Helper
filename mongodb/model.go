package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// An interface for [BaseModel]
	//
	// This interface can be used as a query constraint for documents that wrap [BaseModel] or a similar struct.
	Document[T any] interface {
		InitMongoID() 
		SetUpdatedAt(updatedAt time.Time) 
		SetCreatedAt(createdAt time.Time) 
		InitDocument() 
		ResetMongoID() 
	}

	// BaseModel contains all the fields that most documents should have
	BaseModel struct {
		MongoID   primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
		CreatedAt time.Time           `bson:"createdAt" json:"createdAt,omitempty"`
		UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt,omitempty"`
	}
)

// InitMongoID creates a new MongoID if the existing one is Zero value.
func (b *BaseModel) InitMongoID() {
	if b.MongoID.IsZero() {
		b.MongoID = primitive.NewObjectID()
	}
}

// InitDocument inits a new Document so that it can be inserted into the DB.
// A new MongoDB is generated, and the createdAt and updatedAt are set to the current date.
func (b *BaseModel) InitDocument() {
	b.InitMongoID()
	b.SetCreatedAt(time.Now())
	b.SetUpdatedAt(time.Now())
}

// Sets the MongoID to the zero value.
func (b *BaseModel) ResetMongoID() {
	b.MongoID = primitive.NilObjectID
}

func (b *BaseModel) SetCreatedAt(createdAt time.Time) {
	b.CreatedAt = createdAt
}

func (b *BaseModel) SetUpdatedAt(updatedAt time.Time) {
	b.UpdatedAt = updatedAt
}

func (b *BaseModel) GetMongoID() primitive.ObjectID {
	return b.MongoID
}

func (b *BaseModel) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b *BaseModel) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}
