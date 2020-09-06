package query

import (
	"context"

	"cloud.google.com/go/firestore"
)

func FireRead(ctx context.Context, client *firestore.Client, colName string, value [2]string, isLatest bool, endPage string) ([]*firestore.DocumentSnapshot, error) {
	snaps := make([]*firestore.DocumentSnapshot, 20)
	var err error
	var q firestore.Query

	// value
	if value[0] != "" && value[1] != "" {
		q = client.Collection(colName).Where("Console", "==", value[0]).Where("GameTag", "array-contains", value[1])
	} else if value[0] != "" {
		q = client.Collection(colName).Where("Console", "==", value[0])
	} else if value[1] != "" {
		q = client.Collection(colName).Where("GameTag", "array-contains", value[1])
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
		return nil, err
	}
	return snaps, nil
}

func FireCreateBoard(ctx context.Context, client *firestore.Client, colName string, data map[string]interface{}) (string, error) {
	data["UpdateTime"] = firestore.ServerTimestamp
	ref, _, err := client.Collection(colName).Add(ctx, data)

	return ref.ID, err
}

func FireUpdateBoard(ctx context.Context, client *firestore.Client, colName string, data map[string]interface{}) error {
	page := data["page"].(string)
	delete(data, "page")
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
	if len(m) == 20 {
		endPage = snaps[len(snaps)-1].Ref.ID
	}

	return m, endPage
}
