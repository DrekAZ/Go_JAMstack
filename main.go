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

	//router.GET("/v1/items", request.Get(ctx, client, "all"))
	v1 := router.Group("/v1")
	{
		v1.GET("/search/once-team", request.Search(ctx, fireClient, "OnceTeam"))
		v1.GET("/search/group", request.Search(ctx, fireClient, "Group"))
		v1.POST("/create/once-team", request.Create(ctx, fireClient, "OnceTeam"))
		v1.POST("/create/group", request.Create(ctx, fireClient, "Group"))
		v1.PUT("/update/once-team", request.Update(ctx, fireClient, "OnceTeam"))
		v1.PUT("/update/group", request.Update(ctx, fireClient, "Group"))
	}

	router.Run(":8090")
}
