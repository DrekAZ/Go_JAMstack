package main
import (
  "fmt"
  "log"
  "context"
	//"net/http"
  "os"
  /*"io"
  "io/ioutil"
  "time"
  "strings"
  "flag"*/
	"crypto/rand"
  "encoding/json"
  "reflect"
	"golang.org/x/crypto/bcrypt"
  "github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
  "cloud.google.com/go/firestore"
  //"cloud.google.com/go/storage"
  //"google.golang.org/api/iterator"
  "google.golang.org/api/option"

  "./query"
	"./auth"
	
)

type UserInfo struct {
  Name string `json: "Name"`
	IconPath string `json: "IconPath"`
}
type NewUser struct {
	Name string `json: "Name"`
	IconPath string `json: "IconPath"`
	Address string `json: "Address"`
	Password string `json: "Password"`
}
type Content struct {
  UserID string `json: "UserID"`
  Title string `json: "Title"`
  Markdown string `json:"Markdown"`
  Tags []string `json: "Tags"`
}
type Env struct {
	projectID string
	jsonPath string
	bucket string
	cookieSecret string
}
/*type AuthEnv struct {
	issuer string
	clientID string
	clientSecret string
}*/

func main() {
	env, auth_env := Get_Env()
  ctx := context.Background()

  /*s_client, err := storage.NewClient(ctx, option.WithCredentialsFile(path))
  if err != nil {
    log.Fatalf("Failed to create client: %v", err)
  }*/

  f_client, err := firestore.NewClient(ctx, env.projectID, option.WithCredentialsFile(env.jsonPath)) 
  if err != nil {
    log.Fatalf("Failed to create client: %v", err)
  }

	store := cookie.NewStore([]byte(env.cookieSecret))
  router := gin.Default()
	router.Use(sessions.Sessions("useron", store))

	//router.POST("/signup", signup(ctx, f_client))
	//router.GET("/login", login(ctx, f_client))
	//router.GET("/authUser", authUser(ctx, f_client))
	router.GET("/auth", auth.Auth(ctx, auth_env))
	router.GET("/callback", auth.Callback(ctx, auth_env, f_client))
	//router.GET("/login", login(ctx, f_client))

  router.GET("/get", func(c *gin.Context) {
    name := c.Query("name")
    if name == "" {
      log.Fatalf("No data")
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
    c.JSON(200, gin.H {
			"hello": "U",
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

  router.Run(":8081")
  fmt.Printf("DONE \n")
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
func authUser(ctx context.Context, client *firestore.Client) (gin.HandlerFunc) {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token").(string)

		refs := query.Fire_Read(ctx, client, "users", token)
		if refs == nil {
			log.Fatalf("cannot find address")
		}

		doc, err := refs.Next()
		fmt.Println(reflect.TypeOf(doc.Data()["Address"])) ////
		err = bcrypt.CompareHashAndPassword([]byte(token), []byte(doc.Data()["Address"].(string)))
		if err != nil {
			session.Clear()
			log.Fatal("incorrect token", err)
			/// c.Abort()
		}

		c.JSON(200, gin.H {
			"ok": true,
		})
	}
}

func Struct2Map(data interface{}) (map[string]interface{}) {
  B, err := json.Marshal(data)
  if err != nil {
    log.Fatal("marshal err", err)
    return nil
  }

  var m map[string]interface{}
  err = json.Unmarshal(B, &m)
  if err != nil {
    log.Fatal("unmarshal err", err)
    return nil
  }
  return m 
}

func Get_Env() (Env, auth.AuthEnv) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Main env cannnot load")
	}
	var gcp_env Env
	gcp_env.projectID = os.Getenv("PROJECT_ID")
	gcp_env.jsonPath = os.Getenv("JSON_PATH")
	gcp_env.bucket = os.Getenv("BUCKET")
	gcp_env.cookieSecret = os.Getenv("COOKIE_SECRET")

	var auth_env auth.AuthEnv
	auth_env.Issuer = os.Getenv("AUTH0_DOMAIN")
	auth_env.ClientID = os.Getenv("AUTH0_CLIENT_ID")
	auth_env.ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")
	return gcp_env, auth_env
}


func Rand_Str(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}