package actions

import (
	"caco/services"
	"context"
	"time"

	"github.com/shomali11/slacker"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// BombeiroAction ...
var BombeiroAction = &services.DialogFlowAction{
	Name:    "input.bombeiro",
	Command: evaluateBombeiro,
}

func evaluateBombeiro(result *dialogflowpb.QueryResult, request slacker.Request, response slacker.ResponseWriter) {

	var date time.Time
	var err error
	if len(result.Parameters.Fields) > 0 {
		// verifica se é struct
		var dateTxt string = ""
		strct := result.Parameters.Fields["date-time"].GetStructValue()
		if strct != nil {
			dateTxt = strct.Fields["startDate"].GetStringValue()
			date, err = time.Parse(time.RFC3339, dateTxt)
			if err != nil {
				response.ReportError(err, slacker.WithThreadError(true))
				return
			}
			date = date.Add(24 * time.Hour)
		} else {
			dateTxt = result.Parameters.Fields["date-time"].GetStringValue()
			if dateTxt != "" {
				date, err = time.Parse(time.RFC3339, dateTxt)
				if err != nil {
					response.ReportError(err, slacker.WithThreadError(true))
					return
				}
			} else {
				date = time.Now().Truncate(24 * time.Hour)
			}
		}

	} else {
		date = time.Now().Truncate(24 * time.Hour)
	}

	ctx := context.Background()
	svc := services.NewSpredsheetService(ctx)

	rowResult, err := svc.GetRowByDate(date)
	if err != nil {
		response.ReportError(err, slacker.WithThreadError(true))
		return
	}

	if rowResult == nil {
		response.Reply("Nenhum bombeiro encontrado :sadsad:", slacker.WithThreadReply(true))
		return
	}

	response.Reply("O bombeiro que você está procurando é <@"+rowResult.SlackID+">", slacker.WithThreadReply(true))
}
