package services

import (
	"caco/settings"
	"context"
	"fmt"
	"log"

	"github.com/shomali11/slacker"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// DialogFlowAction ...
type DialogFlowAction struct {
	Name    string
	Command func(*dialogflowpb.QueryResult, slacker.Request, slacker.ResponseWriter)
}

// DetectIntent ...
func DetectIntent(username string, text string) (*dialogflowpb.QueryResult, error) {
	client, err := dialogflow.NewSessionsClient(context.Background())
	defer client.Close()

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	sessionID := username
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", settings.Config.GoogleProjectID, sessionID)

	response, err := client.DetectIntent(
		context.Background(),
		&dialogflowpb.DetectIntentRequest{
			Session: sessionPath,
			QueryInput: &dialogflowpb.QueryInput{
				Input: &dialogflowpb.QueryInput_Text{
					Text: &dialogflowpb.TextInput{
						Text:         text,
						LanguageCode: "pt-br",
					},
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(response)

	return response.GetQueryResult(), nil
}
