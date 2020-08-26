package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/api/option"

	"server_module/auth"
	"server_module/query"
	"server_module/setting"
)

// user info data (comments, profile ...) struct
type UserInfo struct {
	Name     string `json:"Name"`
	IconPath string `json:"IconPath"`
}

// new user datas struct
type NewUser struct {
	Name     string `json:"Name"`
	IconPath string `json:"IconPath"`
	Address  string `json:"Address"`
	Password string `json:"Password"`
}

// Content Struct
type Content struct {
	UserID   string   `json:"UserID"`
	Title    string   `json:"Title"`
	Markdown string   `json:"Markdown"`
	Tags     []string `json:"Tags"`
}

func main() {
	env, authEnv, err := setting.GetEnv()
	if err != nil {
		log.Println(err.Error())
	}
	ctx := context.Background()

	/*s_client, err := storage.NewClient(ctx, option.WithCredentialsFile(path))
	  if err != nil {
	    log.Fatalf("Failed to create client: %v", err)
	  }*/

	fireClient, err := firestore.NewClient(ctx, env.ProjectID, option.WithCredentialsFile(env.JSONPath))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	store := cookie.NewStore([]byte(env.CookieSecret))
	router := gin.Default()
	router.Use(sessions.Sessions("useron", store))
	router.Use(setting.ErrorMiddleware())

	//router.POST("/signup", signup(ctx, f_client))
	//router.GET("/login", login(ctx, f_client))
	//router.GET("/authUser", authUser(ctx, f_client))
	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		s := session.Get("id_token")
		c.JSON(200, gin.H{
			"id_token": s,
		})
	})
	router.GET("/v1/auth/login", auth.Auth(ctx, authEnv))
	router.GET("/v1/auth/callback", auth.Callback(ctx, authEnv, fireClient))
	router.GET("/v1/auth/logout", auth.Logout(ctx, authEnv))
	//router.GET("/login", login(ctx, f_client))

	router.GET("/get", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			err := errors.New("no name")
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		/*err := fire_Read(ctx, client)
		  if err != nil {
		    log.Fatalf("Cannot", err)
		  }*/

		//str := Byte2str(data)
		/*session := sessions.Default(c)
		  if session.Get("hello") != "world" {
		    session.Set("hello", "world")
		    session.Save()
		  }*/
		c.JSON(200, gin.H{
			"hello": name,
		})
	})

	/*router.POST("/update", func(c *gin.Context) {
	  var data Content
	  if err := c.BindJSON(&data); err != nil {
	    log.Fatal(err)
	  }
	  if data.Markdown == "" {
	    log.Fatalf("No query")
	  }*/

	/*err = query.Storage_Write(ctx, s_client, bucket, data.Title, data.Markdown)
	  if err != nil {
	    log.Fatalf("storage", err)
	  }
	  data.Markdown = "https://storage.cloud.google.com/"+ bucket + data.Title*/

	// contents
	//var ref *firestore.DocumentRef

	//data_map := Struct2Map(data)
	//refs := Fire_Read(ctx, f_client, "contents", data_map["Title"].(string))
	/*if refs != nil { // title already exists -> content update but later
	  c.JSON(200, gin.H {
	    "OK": false,
	  })
	}*/

	/*err := Contents_Write(ctx, f_client, data_map)
	  if err != nil {
	    log.Fatal("Contents_Write", err)
	  }*/

	// tags
	//contents_ref := map[string]string{ "contents_ref":ref.Path }
	//err = Tags_Write(ctx, f_client, data_map[""], contents_ref)

	/*c.JSON(200, gin.H {
	    "OK": true,
	  })
	})*/

	router.Run(":8090")
	log.Print("SET UP")
}

/*func login(ctx context.Context, client *firestore.Client, data map[string]interface{}) (gin.HandlerFunc) {
  return func(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    var content Content
    if err := c.BindJSON(&content); err != nil {
      log.Fatal(err)
    }
    data := Struct2Map(content)

    refs := query.Fire_Read(ctx, client, "users", data["Address"])
    if refs == nil {
      log.Fatalf("cannot find address")
    }

    // mail address is unique だから見つかるのは一つだけ
    doc, err := refs.Next()
    fmt.Println(reflect.TypeOf(doc.Data()["Password"])) ////
    err = bcrypt.CompareHashAndPassword([]byte(doc.Data()["Password"].(string)), []byte(data["Password"]))
    if err != nil {
      log.Fatal("incorrect password", err)
      /// c.Abort()
    }

    token, err := bcrypt.GenerateFromPassword([]byte(data["Password"]), bcrypt.DefaultCost) // hash
    if err != nil {
      log.Fatalf("cannnot crypt address")
    }

    /// session store in Firestore
    defer client.Close()
    _, err = doc.Ref.Set(ctx, map[string]interface{}{
      "token": token,
    }, firestore.MergeAll)
    if err != nil {
      log.Fatal("error login add session store", err)
    }

    session.Set("token", token)
    session.Options(sessions.Options{
      MaxAge: 604800,
      Secure: true,
      HttpOnly: true,
      SameSite: http.SameSiteLaxMode,
    })
    session.Save()
    c.JSON(200, gin.H {
      "user": true,
    })
  }
}*/

func AuthUser(ctx context.Context, client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token").(string)

		refs := query.FirestoreRead(ctx, client, "users", token)
		if refs == nil {
			log.Println("cannot find address")
		}

		doc, err := refs.Next()
		fmt.Println(reflect.TypeOf(doc.Data()["Address"])) ////
		err = bcrypt.CompareHashAndPassword([]byte(token), []byte(doc.Data()["Address"].(string)))
		if err != nil {
			session.Clear()
			log.Println("incorrect token", err)
			/// c.Abort()
		}

		c.JSON(200, gin.H{
			"ok": true,
		})
	}
}
