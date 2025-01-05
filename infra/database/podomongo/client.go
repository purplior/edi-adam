package podomongo

import (
	"context"
	"log"
	"time"

	"github.com/purplior/podoroot/application/config"
	"github.com/purplior/podoroot/domain/shared/constant"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	ConstructorOption struct {
		Phase         constant.Phase
		URI           string
		DefaultDbName string
	}

	Client struct {
		*mongo.Client
		databaseMap map[string]*MyMongoDatabase
		opt         ConstructorOption
	}
)

func (c *Client) Connect() error {
	ctx := context.Background()

	if client, err := mongo.Connect(ctx, c.clientOptions()); err != nil {
		return err
	} else {
		c.Client = client
	}

	defaultDBName := c.opt.DefaultDbName
	if len(defaultDBName) > 0 {
		log.Println("[mymongo] use default db name.")
		c.databaseMap[defaultDBName] = NewDatabase(c.Client.Database(defaultDBName))
	}

	if err := c.Client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	log.Println("[mymongo] is connected.")

	if err := initIndexes(ctx, c); err != nil {
		return err
	}

	log.Println("[mymongo] init indexes.")

	return nil
}

func (c *Client) Dispose(
	ctx context.Context,
) error {
	return c.Client.Disconnect(ctx)
}

func (c *Client) MyDatabase(
	name string,
	opts ...*options.DatabaseOptions,
) *MyMongoDatabase {
	db, is := c.databaseMap[name]
	if !is || db == nil {
		c.databaseMap[name] = NewDatabase(c.Database(name, opts...))
		db = c.databaseMap[name]
	}

	return db
}

func (c *Client) MyCollection(
	name string,
	opts ...*options.CollectionOptions,
) *MyMongoCollection {
	db := c.databaseMap[c.opt.DefaultDbName]

	return db.MyCollection(name, opts...)
}

func (c *Client) clientOptions() *options.ClientOptions {
	serverAPIOptions := options.
		ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(c.opt.URI).
		SetServerAPIOptions(serverAPIOptions)

	switch c.opt.Phase {
	case constant.Phase_Production:
		clientOptions.SetHeartbeatInterval(10 * time.Second)
		clientOptions.SetMaxPoolSize(100)
		clientOptions.SetMinPoolSize(1)
		clientOptions.SetMaxConnIdleTime(0)
	default:
		clientOptions.SetHeartbeatInterval(15 * time.Second)
		clientOptions.SetMaxPoolSize(5)
		clientOptions.SetMinPoolSize(1)
		clientOptions.SetMaxConnIdleTime(10 * time.Second)
	}

	return clientOptions
}

func NewClient() *Client {
	opt := ConstructorOption{
		Phase:         config.Phase(),
		URI:           config.MongoDbURI(),
		DefaultDbName: config.MongoDbName(),
	}

	client := &Client{
		opt:         opt,
		databaseMap: make(map[string]*MyMongoDatabase),
	}

	return client
}
