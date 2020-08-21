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

func GetEnv() (*Env, *AuthEnv, error) {
	var gcpEnv Env
	var authEnv AuthEnv

	err := godotenv.Load(".env")
	if err != nil {
		return &gcpEnv, &authEnv, err
	}

	gcpEnv.ProjectID = os.Getenv("PROJECT_ID")
	gcpEnv.JSONPath = os.Getenv("JSON_PATH")
	gcpEnv.Bucket = os.Getenv("BUCKET")
	gcpEnv.CookieSecret = os.Getenv("COOKIE_SECRET")

	authEnv.Issuer = os.Getenv("AUTH0_DOMAIN")
	authEnv.ClientID = os.Getenv("AUTH0_CLIENT_ID")
	authEnv.ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
	return &gcpEnv, &authEnv, nil
}
