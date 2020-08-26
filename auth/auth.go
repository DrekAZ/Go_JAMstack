package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/firestore"
	oidc "github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"server_module/convert"
	"server_module/setting"
	"server_module/status_code"
)

func Auth(ctx context.Context, authEnv *setting.AuthEnv) gin.HandlerFunc {
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

		state, err := convert.Rand2base64(32)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		nonce, err := convert.Rand2base64(32)
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
	}
}

func Callback(ctx context.Context, authEnv *setting.AuthEnv, client *firestore.Client) gin.HandlerFunc {
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
			RedirectURL:  "http://localhost:8090",
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
			err = errors.New("missing token")
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
			err = errors.New("incorrect nonce")
			c.Error(err).SetType(gin.ErrorTypePublic)
			return
		}

		// アプリケーションのデータ構造におとすときは以下のように書く
		var idTokenClaims map[string]interface{}
		err = idToken.Claims(&idTokenClaims)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}

		// session clear
		session.Clear()
		session.Save()
		session.Set("id_token", rawIDToken)
		//session.Set("access_token", oauth2Token)
		//session.Set("profile", idTokenClaims)
		session.Save()
		fmt.Printf("%#v\n", rawIDToken)
		fmt.Printf("%#v\n", idTokenClaims)

		log.Println("認証成功")
		err = Login(ctx, client, idTokenClaims["email"].(string))
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		c.Redirect(http.StatusFound, "http://localhost:8090")
	}
}

func Login(ctx context.Context, client *firestore.Client, email string) error {
	defer client.Close()
	iter := client.Collection("users").Where("email", "==", email).Documents(ctx)
	// email is uniqu
	_, err := iter.Next()

	// not found email
	// sign up -> add email to firestore(GCP)
	if err != nil && err.Error() == status_code.NotFound {
		_, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
			"email": email,
			"time":  time.Now().Format(time.RFC3339Nano),
		})
		if err != nil {
			return err
		}
	}

	return err
}

func Logout(ctx context.Context, authEnv *setting.AuthEnv) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()

		logoutURL, err := url.Parse(authEnv.Issuer)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}

		logoutURL.Path += "v2/logout"
		param := url.Values{}

		var scheme string
		if c.Request.TLS == nil {
			scheme = "http"
		} else {
			scheme = "https"
		}
		returnTo, err := url.Parse(scheme + "://" + c.Request.Host)
		if err != nil {
			c.Error(err).SetType(gin.ErrorTypePrivate)
			return
		}
		param.Add("returnTo", returnTo.String())
		param.Add("client_id", authEnv.ClientID)
		logoutURL.RawQuery = param.Encode()

		c.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
		log.Println("Logout")
	}
}
