package main

import (
	"context"
	Server "cwgo_test/server"

	"cwgo_test/biz/doc/model/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	// 设置MongoDB连接信息
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// 建立MongoDB连接
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// 选择数据库和集合
	collection := client.Database("users").Collection("user")

	// 创建一个新的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := user.NewUser()

	// 向新创建的 User 实例添加数据
	newUser.Id = 1
	newUser.Username = "john_doe"
	newUser.Age = 30
	newUser.City = "New York"
	newUser.Banned = false
	newUser.Contact = &user.UserContact{
		Email: "john@example.com",
		Phone: "123-456-7890",
	}
	newUser.Yd = []user.YDType{user.YDType_UP, user.YDType_DOWN}

	Server.Test(ctx, collection, newUser)

}
