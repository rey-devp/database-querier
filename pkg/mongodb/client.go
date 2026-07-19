package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Client struct {
	client *mongo.Client
	dbName string
}

func NewClient(ctx context.Context, uri string, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
		dbName: dbName,
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) GetDatabase() *mongo.Database {
	return c.client.Database(c.dbName)
}

func (c *Client) ListCollections(ctx context.Context) ([]string, error) {
	names, err := c.GetDatabase().ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return names, nil
}
