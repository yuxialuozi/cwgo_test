package Server

import (
	"context"
	user2 "cwgo_test/biz/doc/dao/user"
	user1 "cwgo_test/biz/doc/model/user"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func Test(ctx context.Context, collection *mongo.Collection, newUser *user1.User) {
	// 使用 NewUserRepository 创建一个新的 UserRepository 实例
	userMongo := user2.NewUserRepository(collection)

	// 调用 InsertUser 函数插入用户文档到 MongoDB
	_, err := userMongo.InsertOne(ctx, newUser)
	// 检查插入操作是否成功
	if err != nil {
		log.Println("插入用户失败:", err)
	} else {
		// 插入操作成功，打印插入的用户ID
		log.Println("成功插入用户")
	}
}
