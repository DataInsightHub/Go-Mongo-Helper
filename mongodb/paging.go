package mongodb

import "go.mongodb.org/mongo-driver/mongo/options"

func GetPaginatedOpts(limit, page int64) *options.FindOptions {
	l := limit
	skip := page * limit - limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}