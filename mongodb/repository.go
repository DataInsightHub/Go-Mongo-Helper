package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	FindOne[T any] interface {
		// Tries to find a Document that matches the given filter, and returns it.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.FindOne]
		FindOne(ctx context.Context, filter bson.M, opts ...*options.FindOneOptions) (T, error)
	}

	FindMany[T any] interface {
		// Finds all Documents that match the given filter, and returns them as a slice.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.Find]
		FindMany(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error)
	}

	InsertOne[T any, D Document[T]] interface {
		// Inserts a document in the db.
		// The document gets a new MongoID, if not already set, and the CreatedAt and UpdatedAt fields are set to the current time.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.InsertOne]
		InsertOne(ctx context.Context, doc D, opts ...*options.InsertOneOptions) (T, error)
	}

	InsertMany[T any, D Document[T]] interface {
		// Inserts multiple documents in the db.
		// All the documents get a new MongoID, if not already set, and the CreatedAt and UpdatedAt are set to the current time.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.InsertMany]
		InsertMany(ctx context.Context, docs []D, opts ...*options.InsertManyOptions) ([]T, error)
	}

	UpdateOne interface {
		// Updates a single document that matches the given filter. updatedAt is automatically set to the current date for the updated document.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.UpdateOne]
		UpdateOne(ctx context.Context, filter bson.M, data primitive.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	}

	UpdateMany interface {
		// Updates multiple document that matches the given filter. updatedAt is automatically set to the current date for the updated documents.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.UpdateMany]
		UpdateMany(ctx context.Context, filter bson.M, data primitive.M, opts ...*options.UpdateOptions) error
	}

	ReplaceOne[T any, D Document[T]] interface {
		// Replaces the specified document.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.ReplaceOne]
		ReplaceOne(ctx context.Context, filter bson.M, doc D, opts ...*options.ReplaceOptions) (T, error)
	}

	DeleteOne[T any] interface {
		// Deletes one document that matches the given filter
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.DeleteOne]
		DeleteOne(ctx context.Context, filter bson.M, opts ...*options.DeleteOptions) error
	}

	DeleteMany[T any] interface {
		// Deletes multiple documents, and returns the number of documents that were deleted
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.DeleteMany]
		DeleteMany(ctx context.Context, filter bson.M, opts ...*options.DeleteOptions) (int, error)
	}

	BulkWrite interface {
		// Does multiple Write and Update operations in one go.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.Bulkwrite]
		//
		// While the mongo-Method returns an error if 0 operations are passed, this method returns an empty result and no error.
		BulkWrite(ctx context.Context, Documents []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	}

	Aggregater interface {
		// Runs an aggregation pipeline.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.Aggregate]
		Aggregate(ctx context.Context, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	}

	Counter interface {
		// Returns the number of documents that match the given filter.
		//
		// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.CountDocuments]
		CountDocuments(ctx context.Context, filter bson.M, opts ...*options.CountOptions) (int, error)
	}

	// RepositoryI is an interfaces for a single mongoDB collection. All mongodb operations are permitted on this repository
	//
	// Please note that a repository always contains data for multiple company.
	// Therefore, most query filters should filter for a specific companyID, see [mongodb.NewFilter] and [mongodb.WithCompanyID]
	RepositoryI[T any, D Document[T]] interface {
		FindOne[T]
		FindMany[T]
		InsertOne[T, D]
		InsertMany[T, D]
		UpdateOne
		UpdateMany
		ReplaceOne[T, D]
		DeleteOne[T]
		DeleteMany[T]
		BulkWrite
		Aggregater
		Counter
	}

	// A Repository represents a single mongoDB collection.
	//
	// Please note that a repository always contains data for multiple company.
	// Therefore, most query filters should filter for a specific companyID, see [mongodb.NewFilter] and [mongodb.WithCompanyID]
	Repository[T any, D Document[T]] struct {
		db *mongo.Collection
	}
)

// Creates a new repository for the specified mongo collection.
func NewRepository[T any, D Document[T]](collection *mongo.Collection) RepositoryI[T, D] {
	return &Repository[T, D]{
		db: collection,
	}
}

// Tries to find a Document that matches the given filter, and returns it.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.FindOne]
func (r *Repository[T, D]) FindOne(ctx context.Context, filter bson.M, opts ...*options.FindOneOptions) (T, error) {
	var res T
	err := r.db.FindOne(ctx, filter, opts...).Decode(&res)

	return res, err
}

// Finds all Documents that match the given filter, and returns them as a slice.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.Find]
func (r *Repository[T, D]) FindMany(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]T, error) {
	var res []T
	cur, err := r.db.Find(ctx, filter, opts...)

	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Inserts a document in the db.
// The document gets a new MongoID, and the CreatedAt and UpdatedAt fields are set to the current time.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.InsertOne]
func (r *Repository[T, D]) InsertOne(ctx context.Context, doc D, opts ...*options.InsertOneOptions) (T, error) {
	doc.InitMongoID()
	doc.SetCreatedAt(time.Now())
	doc.SetUpdatedAt(time.Now())

	_, err := r.db.InsertOne(ctx, doc, opts...)
	if err != nil {
		return doc.GetDoc(), err
	}

	return doc.GetDoc(), nil
}

// Inserts multiple documents in the db.
// All the documents get a new MongoID, if not already set, and the CreatedAt and UpdatedAt are set to the current time.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.InsertMany]
func (r *Repository[T, D]) InsertMany(ctx context.Context, documents []D, opts ...*options.InsertManyOptions) ([]T, error) {
	if len(documents) <= 0 {
		// mongoDB does not allow inserting 0 documents, but that is not an error for us.
		return nil, nil
	}

	docs := make([]interface{}, len(documents))
	copyDocuments := make([]T, len(documents))

	for i, doc := range documents {
		doc.InitMongoID()
		doc.SetCreatedAt(time.Now())
		doc.SetUpdatedAt(time.Now())

		docs[i] = doc
		copyDocuments[i] = doc.GetDoc()
	}

	_, err := r.db.InsertMany(ctx, docs, opts...)
	if err != nil {
		return nil, err
	}

	return copyDocuments, nil
}

// Updates a single document that matches the given filter. updatedAt is automatically set to the current date for the updated document.
// The data parameter determines which fields are set to what value. Operations other than $set are not possible.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.UpdateOne]
func (r *Repository[T, D]) UpdateOne(ctx context.Context, filter bson.M, data primitive.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := r.db.UpdateOne(ctx, filter, bson.M{"$set": data, "$currentDate": bson.M{"updatedAt": true}}, opts...)
	if err != nil {
		return updateResult, fmt.Errorf("%v: %w", "mongodb.Repository.UpdateOne", err)
	}

	return updateResult, nil
}

// Updates multiple document that matches the given filter. updatedAt is automatically set to the current date for the updated documents.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.UpdateMany]
func (r *Repository[T, D]) UpdateMany(ctx context.Context, filter bson.M, data primitive.M, opts ...*options.UpdateOptions) error {
	_, err := r.db.UpdateMany(ctx, filter, bson.M{"$set": data, "$currentDate": bson.M{"updatedAt": true}}, opts...)
	return err
}

// Replaces the specified document.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.ReplaceOne]
func (r *Repository[T, D]) ReplaceOne(ctx context.Context, filter bson.M, doc D, opts ...*options.ReplaceOptions) (T, error) {
	doc.SetUpdatedAt(time.Now())
	_, err := r.db.ReplaceOne(ctx, filter, doc, opts...)
	return doc.GetDoc(), err
}

// Deletes one document that matches the given filter
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.DeleteOne]
func (r *Repository[T, D]) DeleteOne(ctx context.Context, filter bson.M, opts ...*options.DeleteOptions) error {
	if len(filter) == 0 {
		return fmt.Errorf("DeleteOne: Filter can not be empty. Filter: %v", filter)
	}
	_, err := r.db.DeleteOne(ctx, filter, opts...)
	return err
}

// Deletes multiple documents, and returns the number of documents that were deleted
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.DeleteMany]
func (r *Repository[T, D]) DeleteMany(ctx context.Context, filter bson.M, opts ...*options.DeleteOptions) (int, error) {
	if len(filter) == 0 {
		return 0, fmt.Errorf("DeleteMany: Filter can not be empty. Filter: %v", filter)
	}
	res, err := r.db.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return int(res.DeletedCount), err
}

// Does multiple Write and Update operations in one go.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.Bulkwrite]
//
// While the mongo-Method returns an error if 0 operations are passed, this method returns an empty result and no error.
func (r *Repository[T, D]) BulkWrite(ctx context.Context, Documents []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {

	if len(Documents) == 0 {
		return &mongo.BulkWriteResult{}, nil
	}

	return r.db.BulkWrite(ctx, Documents, opts...)
}

// Runs an aggregation pipeline.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.Aggregate]
func (r *Repository[T, D]) Aggregate(ctx context.Context, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return r.db.Aggregate(ctx, pipeline, opts...)
}

// Returns the number of documents that match the given filter.
//
// See [https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Collection.CountDocuments]
func (r *Repository[T, D]) CountDocuments(ctx context.Context, filter bson.M, opts ...*options.CountOptions) (int, error) {
	count, err := r.db.CountDocuments(ctx, filter, opts...)
	return int(count), err
}
