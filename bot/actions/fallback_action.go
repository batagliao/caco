package actions

import (
	"caco/services"

	"github.com/shomali11/slacker"

	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// FallbackAction ...
var FallbackAction = &services.DialogFlowAction{
	Name:    "input.welcome",
	Command: evaluateTeamMRs,
}

func defaultResponse(result *dialogflowpb.QueryResult, request slacker.Request, response slacker.ResponseWriter) {
	response.Reply(result.GetFulfillmentText(), slacker.WithThreadReply(true))
}
