package auth

import (
	"fmt"
  "context"
	"log"
	"net/http"
	"crypto/rand"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"cloud.google.com/go/firestore"
	oidc "github.com/coreos/go-oidc"

	"../codes"
)

type AuthEnv struct {
	Issuer string
	ClientID string
	ClientSecret string
}

func Auth(ctx context.Context, auth_env AuthEnv) (gin.HandlerFunc) {
	return func(c *gin.Context) {
    provider, err := oidc.NewProvider(c, auth_env.Issuer)
    if err != nil {
        log.Fatal(err)
    }
    config := oauth2.Config{
        ClientID:     auth_env.ClientID,
        ClientSecret: auth_env.ClientSecret,
        Endpoint:     provider.Endpoint(),
        RedirectURL:  "http://localhost:8081/callback",
        Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
    }

		state, err := Rand_Str(13)
		if err != nil {
			log.Fatal("rand state", err)
		}
		nonce, err := Rand_Str(17)
		if err != nil {
			log.Fatal("rand nonce", err)
		}
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		session.Set("state", state)
		session.Set("nonce", nonce)
		session.Save()

    authURL:= config.AuthCodeURL(state, oidc.Nonce(nonce))
		c.Redirect(http.StatusFound, authURL)
		fmt.Printf("Auth finish")
	}
}

func Callback(ctx context.Context, auth_env AuthEnv, client *firestore.Client) (gin.HandlerFunc) {
	return func(c *gin.Context) {
    //ctx := r.Context()

    // この部分は /auth のコードと同じ
    provider, err := oidc.NewProvider(ctx, auth_env.Issuer)
    if err != nil {
        log.Fatal(err)
    }
    config := oauth2.Config{
        ClientID:     auth_env.ClientID,
        ClientSecret: auth_env.ClientSecret,
        Endpoint:     provider.Endpoint(),
        RedirectURL:  "http://localhost:8081/login",
        Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
    }

		// session(cookie)
    s := c.Request.URL.Query().Get("state")
		session := sessions.Default(c)
		
    // stateが返ってくるので認証画面へのリダイレクト時に渡したパラメータと矛盾がないか検証
		if s != session.Get("state") {
			log.Fatal("incorrect state")
			return
		}

    // codeをもとにトークンエンドポイントから IDトークン を取得
    code := c.Request.URL.Query().Get("code")
    oauth2Token, err := config.Exchange(ctx, code)
    if err != nil {
        //http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
				log.Fatal("code", err)
        return
    }

    // IDトークンを取り出す
    rawIDToken, ok := oauth2Token.Extra("id_token").(string)
    if !ok {
        //http.Error(w, "missing token", http.StatusInternalServerError)
				log.Fatal("missing token")
        return
    }

    oidcConfig := &oidc.Config{
        ClientID: auth_env.ClientID,
    }
		// use the nonce source to create a custom ID Token verifier
    verifier := provider.Verifier(oidcConfig)

    // IDトークンの正当性の検証
    idToken, err := verifier.Verify(ctx, rawIDToken)
    if err != nil {
        //http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
				log.Fatal("Failed to verify ID token", err)
        return
    }
		if idToken.Nonce != session.Get("nonce") {
			log.Fatal("incorrect nonce")
			return
		}

    // アプリケーションのデータ構造におとすときは以下のように書く
    idTokenClaims := map[string]interface{}{}
    if err := idToken.Claims(&idTokenClaims); err != nil {
        //http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)
        return
    }

		// session clear
		session.Clear()
		session.Save()
		session.Set("id_token", rawIDToken)
		//session.Set("access_token", oauth2Token)
		session.Set("profile", idTokenClaims)
    fmt.Printf("%#v", idTokenClaims)

    fmt.Printf("認証成功")
		Login(c, ctx, client, idTokenClaims["email"].(string))
		//c.Redirect(http.StatusOK, "http://localhost:8080")
	}
}

func Login(c *gin.Context, ctx context.Context, client *firestore.Client, email string) {

	defer client.Close()
	iter := client.Collection("users").Where(email, "==", true).Documents(ctx)
	// email is uniqu
	_, err := iter.Next()
	fmt.Printf("%s", err.Error())

	if err.Error() == codes.NotFound {
		_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
			"email": email,
		})
		if err != nil {
			log.Fatal("Log in write email", err)
		}
	} else if err != nil {
		log.Fatal("Log in", err)
	}

	c.Redirect(http.StatusOK, "http://localhost:8081")
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