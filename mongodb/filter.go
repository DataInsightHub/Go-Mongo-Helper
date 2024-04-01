package mongodb

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// FilterOption is query building block that can be combined into a full filter for a mongodb query.
	//
	// See [NewFilter]
	FilterOption interface {
		// Apply applies the given FilterOption to the final filter.
		Apply(primitive.M)
	}
)

// NewFilter builds a new MongoDB query condition, depending on the FilterOptions passed.
//
// It can be used like this to create the common companyID filter:
//
//	filter := NewFilter(WithCompanyID(companyID))
//
// There are also convenience methods to create a filter by MongoID and by CompanyID: [MongoIDFilter] and [CompanyIDFilter]
func NewFilter(opts ...FilterOption) primitive.M {
	f := primitive.M{}

	for _, opt := range opts {
		opt.Apply(f)
	}

	return f
}

type withMongoID primitive.ObjectID

func (w withMongoID) Apply(m primitive.M) {
	m["_id"] = primitive.ObjectID(w)
}

// WithMongoID creates a new [FilterOption] by the mongoID.
func WithMongoID(id primitive.ObjectID) FilterOption {
	return withMongoID(id)
}

// MongoIDFilter creates a new filter by the mongoID.
//
// CAUTION: A query should almost always contain the companyID, or the competitorID for additional safety.
func MongoIDFilter(id primitive.ObjectID) primitive.M {
	return NewFilter(WithMongoID(id))
}

// In creates an $in query-condition for the given array.
// The result is not intended to be used as the root of a query, but as a field-query.
//
//	repository.UpdateMany(ctx,
//		bson.M{"_id": mongodb.In(outboundLogIds)},
//		bson.M{"isFinished": true}
//	)
func In[T comparable](array []T) primitive.M {
	return primitive.M{"$in": array}
}
