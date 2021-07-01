package dbconn

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConn(ctx context.Context, addr, port string) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", addr, port)))
}
