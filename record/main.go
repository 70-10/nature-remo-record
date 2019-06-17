package main

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	firebase "firebase.google.com/go"
	"github.com/70-10/nature-remo-go"
	"google.golang.org/api/option"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler() (Response, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	firestore, err := app.Firestore(ctx)

	if err != nil {
		return Response{}, err
	}
	defer firestore.Close()

	client := natureremo.NewClient("<ACCESS TOKEN>")
	devices, err := client.GetDevices()
	if err != nil {
		return Response{StatusCode: 400}, err
	}

	now := time.Now()

	for _, d := range devices {
		dRef := firestore.Collection(d.ID)
		_, err = dRef.Doc("data").Set(ctx, map[string]interface{}{
			"name":               d.Name,
			"created_at":         d.CreatedAt,
			"updated_at":         d.UpdatedAt,
			"firmware_version":   d.FirmwareVersion,
			"temperature_offset": d.TemperatureOffset,
			"humidity_offset":    d.HumidityOffset,
		})
		if err != nil {
			return Response{StatusCode: 500}, err
		}

		_, err = firestore.Collection("temperature").Doc(now.Format("2006-01-02 15:04:05")).Set(ctx, d.NewestEvents.Temperature)
		if err != nil {
			return Response{StatusCode: 500}, err
		}

		_, err = firestore.Collection("humidity").Doc(now.Format("2006-01-02 15:04:05")).Set(ctx, d.NewestEvents.Humidity)
		if err != nil {
			return Response{StatusCode: 500}, err
		}
		_, err = firestore.Collection("illumination").Doc(now.Format("2006-01-02 15:04:05")).Set(ctx, d.NewestEvents.Illumination)
		if err != nil {
			return Response{StatusCode: 500}, err
		}

	}

	var buf bytes.Buffer

	body, err := json.Marshal(devices)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
