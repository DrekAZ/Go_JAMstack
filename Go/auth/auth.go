package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	oidc "github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"../codes"
)

type AuthEnv struct {
	Issuer       string
	ClientID     string
	ClientSecret string
}

func Auth(ctx context.Context, authEnv *AuthEnv) gin.HandlerFunc {
	return func(c *gin.Context) {
		provider, err := oidc.NewProvider(c, authEnv.Issuer)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		config := oauth2.Config{
			ClientID:     authEnv.ClientID,
			ClientSecret: authEnv.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  "http://localhost:8090/v1/auth/callback",
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		}

		state, err := randStr(13)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		nonce, err := randStr(17)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		session.Set("state", state)
		session.Set("nonce", nonce)
		session.Save()

		authURL := config.AuthCodeURL(state, oidc.Nonce(nonce))
		c.Redirect(http.StatusFound, authURL)
		log.Print("Auth finish -> go callback")
	}
}

func Callback(ctx context.Context, authEnv *AuthEnv, client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// この部分は /auth のコードと同じ
		provider, err := oidc.NewProvider(ctx, authEnv.Issuer)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		config := oauth2.Config{
			ClientID:     authEnv.ClientID,
			ClientSecret: authEnv.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  "http://localhost:8090/login",
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		}

		// session(cookie)
		s := c.Request.URL.Query().Get("state")
		session := sessions.Default(c)

		// stateが返ってくるので認証画面へのリダイレクト時に渡したパラメータと矛盾がないか検証
		if s != session.Get("state") {
			err := errors.New("incorrect state")
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		// codeをもとにトークンエンドポイントから IDトークン を取得
		code := c.Request.URL.Query().Get("code")
		oauth2Token, err := config.Exchange(ctx, code)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		// IDトークンを取り出す
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			//http.Error(w, "missing token", http.StatusInternalServerError)
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		oidcConfig := &oidc.Config{
			ClientID: authEnv.ClientID,
		}
		// use the nonce source to create a custom ID Token verifier
		verifier := provider.Verifier(oidcConfig)

		// IDトークンの正当性の検証
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			//http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			//log.Fatal("Failed to verify ID token", err)
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}
		if idToken.Nonce != session.Get("nonce") {
			//log.Fatal("incorrect nonce")
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		// アプリケーションのデータ構造におとすときは以下のように書く
		idTokenClaims := map[string]interface{}{}
		if err := idToken.Claims(&idTokenClaims); err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			c.Error(err).SetType(gin.ErrorTypePrivate)
			//log.Fatal(err)
			return
		}

		// session clear
		session.Clear()
		session.Save()
		session.Set("id_token", rawIDToken)
		//session.Set("access_token", oauth2Token)
		session.Set("profile", idTokenClaims)
		//fmt.Printf("%#v", idTokenClaims)

		log.Print("認証成功")
		err = login(ctx, client, idTokenClaims["email"].(string))
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		c.Redirect(http.StatusOK, "http://localhost:8080")
	}
}

func login(ctx context.Context, client *firestore.Client, email string) error {
	defer client.Close()
	iter := client.Collection("users").Where(email, "==", true).Documents(ctx)
	// email is uniqu
	_, err := iter.Next()
	//fmt.Printf("%s", err.Error())

	// not found email
	// sign up -> add email to firestore(GCP)
	if err.Error() == codes.NotFound {
		_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
			"email": email,
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	//c.Redirect(http.StatusOK, "http://localhost:8090")
	return nil
}

func randStr(digit uint32) (string, error) {
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
