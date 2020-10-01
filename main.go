package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"google.golang.org/api/option"

	"server_module/request"
	"server_module/setting"
)

func main() {
	env, err := setting.GetEnv()
	if err != nil {
		log.Println(err.Error())
	}
	ctx := context.Background()

	fireClient, err := firestore.NewClient(ctx, env.ProjectID, option.WithCredentialsFile(env.JSONPath))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	router := gin.Default()
	router.Use(setting.ErrorMiddleware())

	v1 := router.Group("game-board/v1")
	{
		v1.GET("/search/OnceTeam", request.Search(ctx, fireClient, "OnceTeam"))
		v1.GET("/search/Group", request.Search(ctx, fireClient, "Group"))
		v1.POST("/create/OnceTeam", request.Create(ctx, fireClient, "OnceTeam"))
		v1.POST("/create/Group", request.Create(ctx, fireClient, "Group"))
		v1.PUT("/update/OnceTeam", request.Update(ctx, fireClient, "OnceTeam"))
		v1.PUT("/update/Group", request.Update(ctx, fireClient, "Group"))
	}

	router.Run(":8000")
}
