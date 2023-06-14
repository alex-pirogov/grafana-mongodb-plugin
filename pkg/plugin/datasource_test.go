package plugin

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestQueryData(t *testing.T) {
	var client mongo.Client
	var err error
	if client, err = NewMongoClient(); err != nil {
		panic(err)
	}

	ds := &Datasource{&client}

	payload_raw := queryModel{
		"juniors",
		"user_actions",
		`{ "$group": { "_id": "$time", "total": { "$sum": 1 } } }`,
	}

	payload, _ := json.Marshal(payload_raw)

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A", JSON: payload},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}
