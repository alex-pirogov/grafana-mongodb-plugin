package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

const uri = "mongodb://juniors:123456@alex-pirogov.ru:27017/"

func NewMongoClient() (mongo.Client, error) {
	fmt.Println("Connecting to MongoDB...")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB connected")

	return *client, nil
}

func NewDatasource(_ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	fmt.Println("Datasource Init")

	var client mongo.Client
	var err error
	if client, err = NewMongoClient(); err != nil {
		panic(err)
	}

	return &Datasource{&client}, nil
}

type Datasource struct {
	MongoClient *mongo.Client
}

func (d *Datasource) Dispose() {
	if err := d.MongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	Db          string
	Collection  string
	Aggregation string
}

type Goal struct {
	User int64
	Goal string
	Time primitive.DateTime
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse
	var err error
	var qm queryModel

	fmt.Printf("JSON: %s\n", query.JSON)
	log.Printf("JSON: %s\n", query.JSON)

	err = json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	fmt.Printf("DB: %s, Coll: %s\n", qm.Db, qm.Collection)

	coll := d.MongoClient.Database(qm.Db).Collection(qm.Collection)

	var decode_result bson.D
	fmt.Printf("Raw aggregation: %s\n", qm.Aggregation)
	bson.UnmarshalExtJSON([]byte(qm.Aggregation), false, &decode_result)
	fmt.Printf("Decoded aggregation: %s\n\n", decode_result)

	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{decode_result})
	if err != nil {
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	fmt.Printf("Fetch results: %s\n\n", results)

	frame := data.NewFrame("response")

	Dataframe := make(map[string][]interface{})

	for k := range results[0] {
		Dataframe[k] = []interface{}{}
	}

	for _, result := range results {
		for k, v := range result {
			Dataframe[k] = append(Dataframe[k], v)
		}
	}

	for k := range results[0] {
		fmt.Printf("Key[%s]\n", k)
		series := Dataframe[k]

		var elementType reflect.Type

		if _, ok := series[0].(primitive.DateTime); ok {
			elementType = reflect.TypeOf(primitive.DateTime(1).Time())
		} else {
			elementType = reflect.TypeOf(series[0])
		}

		fmt.Printf("Element type[%s]\n", elementType)

		sliceType := reflect.SliceOf(elementType)
		slice := reflect.MakeSlice(sliceType, 0, 0)

		for _, v := range series {
			if t, ok := v.(primitive.DateTime); ok {
				v = t.Time()
			}

			slice = reflect.Append(slice, reflect.ValueOf(v))
		}

		fmt.Printf("Series[%s]\n", slice)

		frame.Fields = append(frame.Fields,
			data.NewField(k, nil, slice.Interface()),
		)
	}

	response.Frames = append(response.Frames, frame)
	return response
}

func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
