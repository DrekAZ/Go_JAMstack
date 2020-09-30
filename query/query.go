package query

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func FireRead(ctx context.Context, client *firestore.Client, colName string, value [3]string, isLatest bool, endPage string) ([]*firestore.DocumentSnapshot, error) {
	snaps := make([]*firestore.DocumentSnapshot, 20)
	var err error
	q := client.Collection(colName).Where("IsPublic", "==", true)

	if value[0] != "" {
		q = q.Where("Tool", "==", value[0])
	}
	if value[1] != "" {
		q = q.Where("Console", "==", value[1])
	}
	if value[2] != "" {
		q = q.Where("GameTag", "array-contains", value[2])
	}

	// latest
	if isLatest {
		q = q.OrderBy("UpdateTime", firestore.Desc)
	} else {
		q = q.OrderBy("UpdateTime", firestore.Asc)
	}

	// pagenation
	if endPage != "" {
		pageSnap, _ := client.Collection(colName).Doc(endPage).Get(ctx)
		snaps, err = q.StartAfter(pageSnap.Data()["UpdateTime"]).Limit(20).Documents(ctx).GetAll() //gametag -> console ???
	} else {
		snaps, err = q.Limit(20).Documents(ctx).GetAll()
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return snaps, nil
}

func FireCreateBoard(ctx context.Context, client *firestore.Client, colName string, data map[string]interface{}) (string, error) {
	data["UpdateTime"] = firestore.ServerTimestamp
	ref, _, err := client.Collection(colName).Add(ctx, data)

	return ref.ID, err
}

func FireUpdateBoard(ctx context.Context, client *firestore.Client, colName string, data map[string]interface{}, page string) error {
	data["UpdateTime"] = firestore.ServerTimestamp
	_, err := client.Collection(colName).Doc(page).Set(ctx, data, firestore.MergeAll)

	return err
}

func FireReadContent(snaps []*firestore.DocumentSnapshot) ([]map[string]interface{}, string) {
	m := make([]map[string]interface{}, 20)

	for i, doc := range snaps {
		doc.DataTo(&m[i])
	}

	endPage := ""
	if len(snaps) == 20 {
		endPage = snaps[len(snaps)-1].Ref.ID
	}

	return m, endPage
}
