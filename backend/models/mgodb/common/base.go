package common

import (
	"context"
	"fmt"
	"learnerai/conf"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Mongo *mongo.Client
}

var DB *Database

// 初始化
func Setup() {
	DB = &Database{
		Mongo: SetConnect(),
	}
}

// 连接设置
func SetConnect() *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", conf.GlobalEdgeSetting.Mongo.User, conf.GlobalEdgeSetting.Mongo.Password, conf.GlobalEdgeSetting.Mongo.Host)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetMaxPoolSize(20)) // 连接池
	if err != nil {
		fmt.Println("Connet To mongo err: %+v %+v", uri, err)
		os.Exit(1)
	}
	return client
}

func GetContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return
}
