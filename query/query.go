package query

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type JsonData struct {
	Name     string `json:"name"`
	Markdown string `json:"markdown"`
}

/*func main() {
  projectID := "you-know-275301"
  path := "./YOU-KNOW-be4a1d88e2c3.json"
  bucket := "you-know-275301.appspot.com"
  ctx := context.Background()
  client, err := storage.NewClient(ctx, projectID, option.WithCredentialsFile(path))
  if err != nil {
    log.Fatalf("Failed to create client: %v", err)
  }

  ctx, cancel := context.WithTimeout(ctx, time.Second*1000)
  defer cancel()

	data, err := read(client, bucket, name)
	if err != nil {
		log.Fatalf("Cannot", err)
	}

	str := Byte2str(data)

	c.JSON(200, gin.H {
		"content": str,
	})

  router.POST("/update", func(c *gin.Context) {
    var json_data JsonData
    if err := c.BindJSON(&json_data); err != nil {
      log.Fatal(err)
    }
    if json_data.Name == "" {
      log.Fatalf("No query")
    }

    err := write(client, bucket, json_data.Name, json_data.Markdown)
    if err != nil {
      log.Fatalf("Cannot", err)
    }

    c.JSON(200, gin.H {
      "OK": true,
    })
  })
}*/

/////////////////////////////
func FirestoreRead(ctx context.Context, client *firestore.Client, read_type string, key string) *firestore.DocumentIterator {
	refs := client.Collection(read_type).Where(key, "==", true).Documents(ctx)
	return refs
}
func FirestoreUpdate(ctx context.Context, client *firestore.Client, bucket, write_type string, key string, value string) error {
	defer client.Close()

	refs := FirestoreRead(ctx, client, write_type, key)

	for {
		ref, err := refs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		_, err = ref.Ref.Update(ctx, []firestore.Update{
			{Path: key, Value: value},
		})
	}
	return nil
}
func FireContentsWrite(ctx context.Context, client *firestore.Client, data map[string]interface{}) error {
	batch := client.Batch()
	ref := client.Collection("contents").Doc(data["Title"].(string))
	// data["UserID"] = ref.Path

	// write_type == users, contents, tags
	batch.Set(ref, data, firestore.MergeAll)
	// tags write
	tags, _ := data["Tags"].([]interface{})
	//fmt.Println(tags[0].(string))
	for _, tag := range tags { // math science ...
		t_ref := client.Collection("tags").Doc(tag.(string))
		batch.Set(t_ref, map[string]string{"content_id": ref.ID}, firestore.MergeAll)
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func StorageWrite(ctx context.Context, client *storage.Client, bucket, title string, markdown string) error {
	f, err := os.Open(title)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5000)
	defer cancel()

	wc := client.Bucket(bucket).Object(title).NewWriter(ctx)
	_, err = fmt.Fprintf(wc, markdown)
	if err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	return nil
}
func StorageRead(ctx context.Context, client *storage.Client, bucket, name string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5000)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(name).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return data, nil
}
