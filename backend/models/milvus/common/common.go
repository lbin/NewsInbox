package common

import (
	"context"
	"learnerai/conf"
	appContext "learnerai/golib/app/context"
	"learnerai/golib/app/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

const (
	databaseArxiv = "arxiv"
)

var commonClient client.Client

func Setup() {
	ctx := appContext.GetGinContextWithRequestId()
	c, err := newClient(ctx)
	if err != nil {
		panic(err)
	}
	commonClient = c
}

func newClient(ctx *gin.Context) (client.Client, error) {
	c, err := client.NewClient(ctx, client.Config{
		Address:  conf.GlobalEdgeSetting.Milvus.Host,
		Username: conf.GlobalEdgeSetting.Milvus.User,
		Password: conf.GlobalEdgeSetting.Milvus.Password,
	})
	if err != nil {
		logger.Errorf(ctx, "failed to create new client: %+v", err)
		return nil, err
	}
	ver, err := c.GetVersion(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to retrieve server version: %+v", err)
		return nil, err
	}
	logger.Infof(ctx, "connected to milvus server: %v", ver)

	err = createDatabaseIfNotExist(ctx, c, databaseArxiv)
	if err != nil {
		logger.Errorf(ctx, "failed to touch database %q: %+v", databaseArxiv, err)
		return nil, err
	}

	err = c.UsingDatabase(ctx, databaseArxiv)
	if err != nil {
		logger.Errorf(ctx, "failed to use database %q: %+v", databaseArxiv, err)
		return nil, err
	}

	return c, nil
}

// NOTE: user should not close the shared client!
func GetClient() client.Client {
	return commonClient
}

func createDatabaseIfNotExist(ctx *gin.Context, c client.Client, dbName string) error {
	dbs, err := c.ListDatabases(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to list databases: %+v", err)
		return err
	}
	for i := range dbs {
		if dbs[i].Name == dbName {
			logger.Infof(ctx, "found required database %q", dbName)
			return nil
		}
	}
	err = c.CreateDatabase(ctx, dbName)
	if err != nil {
		logger.Errorf(ctx, "failed to create database %q: %+v", dbName, err)
		return err
	}
	logger.Infof(ctx, "database %q created", dbName)
	return nil
}

func GetContext(ctx *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Minute)
}
