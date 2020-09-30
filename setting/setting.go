package setting

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Env is Struct for Google Cloud Platform and cookie
type Env struct {
	ProjectID    string
	JSONPath     string
	Bucket       string
	CookieSecret string
}

// AuthEnv is Struct for Auth0
type AuthEnv struct {
	Issuer       string
	ClientID     string
	ClientSecret string
}

// Party is Struct for recruit once team members
type OnceTeam struct {
	Tool       string   `json:"Tool" binding:"required"`
	AppLink    string   `json:"AppLink" binding:"required"`
	Name       string   `json:"Name" binding:"required"`
	Console    string   `json:"Console" binding:"required"`
	Msg        string   `json:"Msg" binding:"required"`
	GameTag    []string `json:"GameTag" binding:"required"`
	RecruitCnt []int    `json:"RecruitCnt" binding:"required"`
	IsPublic   bool     `json:"IsPublic"`
}

// Group is Struct for recruit discord members
type Group struct {
	Tool     string   `json:"Tool" binding:"required"`
	AppLink  string   `json:"AppLink" binding:"required"`
	Name     string   `json:"Name" binding:"required"`
	Console  string   `json:"Console" binding:"required"`
	Msg      string   `json:"Msg" binding:"required"`
	GameTag  []string `json:"GameTag" binding:"required"`
	IsPublic bool     `json:"IsPublic"`
}

type UpdateOnceTeam struct {
	Data OnceTeam `json:"Data"`
	Page string   `json:"Page" binding:"required"`
}

type UpdateGroup struct {
	Data Group  `json:"Data"`
	Page string `json:"Page" binding:"required"`
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			log.Println(err.Err)
			c.AbortWithStatus(400)
		}

		err = c.Errors.ByType(gin.ErrorTypePrivate).Last()
		if err != nil {
			log.Println(err.Err)
			c.AbortWithStatus(500)
		}
	}
}

func GetEnv() (*Env, error) {
	var gcpEnv Env

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	gcpEnv.ProjectID = os.Getenv("PROJECT_ID")
	gcpEnv.JSONPath = os.Getenv("JSON_PATH")
	gcpEnv.Bucket = os.Getenv("BUCKET")
	gcpEnv.CookieSecret = os.Getenv("COOKIE_SECRET")

	return &gcpEnv, nil
}
