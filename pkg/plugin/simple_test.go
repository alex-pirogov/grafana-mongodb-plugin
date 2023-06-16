package plugin

import (
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestTypes(t *testing.T) {
	var json_string string
	// var err error
	var result []bson.D

	json_string = `{
			"$group": { "_id": "$goal", "total": { "$sum": 1 } }
		}`

	result, _ = ParseAggregations(json_string)
	fmt.Printf("%s\n", result)

	json_string = `[
		{
			"$group": { "_id": "$goal", "total": { "$sum": 1 } }
		},
		{
			"$sort": { "_id": -1 }
		}
	]
	`
	result, _ = ParseAggregations(json_string)
	fmt.Printf("%s\n", result)
}
