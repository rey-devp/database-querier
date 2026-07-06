package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ExecuteFind(ctx context.Context, coll *mongo.Collection, filter interface{}, projection interface{}, sort interface{}, limit int64) ([]bson.M, int, error) {
	if filter == nil {
		filter = bson.M{}
	}

	findOpts := options.Find()
	if projection != nil {
		findOpts.SetProjection(projection)
	}
	if sort != nil {
		findOpts.SetSort(sort)
	}
	if limit > 0 {
		findOpts.SetLimit(limit)
	}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, int(total), nil
}

func ExecuteAggregate(ctx context.Context, coll *mongo.Collection, pipeline interface{}) ([]bson.M, int, error) {
	if pipeline == nil {
		return nil, 0, fmt.Errorf("pipeline cannot be empty for aggregate operation")
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	return results, len(results), nil
}

func ExecuteCountDocuments(ctx context.Context, coll *mongo.Collection, filter interface{}) ([]bson.M, int, error) {
	if filter == nil {
		filter = bson.M{}
	}

	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// For countDocuments, we just return the count as the total, and empty docs
	return []bson.M{}, int(total), nil
}
