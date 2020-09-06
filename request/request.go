package request

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"server_module/convert"
	"server_module/query"
)

func Search(ctx context.Context, client *firestore.Client, colName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		console := c.Query("console")
		tag := c.Query("tag")
		page := c.Query("page")
		isLatest := convert.Str2bool(c.DefaultQuery("isLatest", "true"))
		snaps := make([]*firestore.DocumentSnapshot, 20)
		var err error

		if console != "" && tag != "" {
			snaps, err = query.FireRead(ctx, client, colName, [2]string{console, tag}, isLatest, page)
		} else if console != "" {
			snaps, err = query.FireRead(ctx, client, colName, [2]string{console, ""}, isLatest, page)
		} else if tag != "" {
			snaps, err = query.FireRead(ctx, client, colName, [2]string{"", tag}, isLatest, page)
		}

		if err != nil || (console == "" && tag == "") {
			e := errors.New("cannot get query")
			c.Error(e).SetType(gin.ErrorTypePublic)
			return
		}

		data, endPage := query.FireReadContent(snaps)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
			"page": endPage,
		})
	}
}

func Create(ctx context.Context, client *firestore.Client, colName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := convert.BindJson2map(c, colName)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}

		page, err := query.FireCreateBoard(ctx, client, colName, data)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"page": page,
		})
	}
}

// switch IsPublic
// {"page": pageurl, ...}
func Update(ctx context.Context, client *firestore.Client, colName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := convert.BindJson2map(c, colName)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}

		err = query.FireUpdateBoard(ctx, client, colName, data)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}

		c.String(http.StatusNoContent, "OK")
	}
}
