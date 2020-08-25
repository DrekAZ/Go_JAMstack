package auth

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	fireClient, err := firestore.NewClient(ctx, "xxx", option.WithCredentialsFile(".json"))
	err = Login(ctx, fireClient, "a@a.com")
	t.Log(err)

}
