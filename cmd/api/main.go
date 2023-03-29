package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"flag"
	"os"
	"runtime/debug"
	"sync"

	"herrfennessey/book-recommender-api-v2/internal/env"
	"herrfennessey/book-recommender-api-v2/internal/leveledlog"
)

func main() {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)

	err := run(logger)
	if err != nil {
		trace := debug.Stack()
		logger.Fatal(err, trace)
	}
}

type config struct {
	env          string
	gcpProjectId string
	httpPort     int
}

type application struct {
	config   config
	logger   *leveledlog.Logger
	dsClient *datastore.Client
	wg       sync.WaitGroup
}

func run(logger *leveledlog.Logger) error {
	var cfg config

	cfg.env = env.GetString("ENV", "prod")
	cfg.gcpProjectId = env.GetString("GCP_PROJECT_ID", "test-project")
	cfg.httpPort = env.GetInt("PORT", 8080)

	flag.Parse()

	ctx := context.Background()
	client, err := createDatastoreClient(ctx, cfg.gcpProjectId)
	if err != nil {
		logger.Info("Failed to connect to datastore client with error %s", err)
		return err
	}
	logger.Info("Connected to Datastore client!")

	app := &application{
		config:   cfg,
		logger:   logger,
		dsClient: client,
	}

	return app.serveHTTP()
}

func createDatastoreClient(ctx context.Context, projectId string) (*datastore.Client, error) {
	dsClient, err := datastore.NewClient(ctx, projectId)
	if err != nil {
		return &datastore.Client{}, err
	}

	return dsClient, nil
}
