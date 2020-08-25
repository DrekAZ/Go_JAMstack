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

	/*t.Run("exist_email", func(t *testing.T) {
		t.Log("exist_email")
		_ = Login(ctx, fireClient, "kassuruv2@gmail.com")
	})

	t.Run("not_exist_email", func(t *testing.T) {
		t.Log("not_exist_email")
		_ = Login(ctx, fireClient, "k2@gmail.com")
	})*/

}
