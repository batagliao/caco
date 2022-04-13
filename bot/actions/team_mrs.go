package actions

import (
	"caco/services"
	"caco/settings"
	"fmt"
	"strconv"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// TeamMRsAction ...
var TeamMRsAction = &services.DialogFlowAction{
	Name:    "input.team-mrs",
	Command: evaluateTeamMRs,
}

const (
	APPROVAL_UPVOTES = 2
)

func evaluateTeamMRs(result *dialogflowpb.QueryResult, request slacker.Request, response slacker.ResponseWriter) {
	team := settings.Config.GitlabTeam

	if team == "" {
		response.Reply("Hi-ho, nenhum time aqui", slacker.WithThreadReply(true))
		return
	}

	svc := services.NewGitlabService()

	group, err := svc.GetGroupByName(team)
	if err != nil {
		response.ReportError(err, slacker.WithThreadError(true))
		return
	}

	if group == nil {
		s := fmt.Sprintf("Grupo `%s` não encontrado no GitLab", team)
		response.Reply(s, slacker.WithThreadReply(true))
	}

	projects, err := svc.GetGroupProjects(group.ID)
	if err != nil {
		response.ReportError(err, slacker.WithThreadError(true))
		return
	}

	for _, proj := range projects {
		go drillDownProjects(response, svc, proj)
	}
}

func drillDownProjects(response slacker.ResponseWriter, svc *services.GitlabService, proj *gitlab.Project) {
	if settings.Config.Debug {
		fmt.Println("Finding MRs for project " + proj.Name)
	}

	mrs, err := svc.GetProjectOpenedMergeRequests(proj.ID)

	if err != nil {
		response.ReportError(err, slacker.WithThreadError(true))
		return
	}

	for _, mr := range mrs {
		go replyMR(response, proj, mr)
	}

}

func replyMR(response slacker.ResponseWriter, proj *gitlab.Project, mr *gitlab.MergeRequest) {
	response.Reply("", slacker.WithAttachments(getMergeRequestAttachment(proj, mr)), slacker.WithThreadReply(true))
}

func getMergeRequestAttachment(proj *gitlab.Project, mr *gitlab.MergeRequest) []slack.Attachment {
	attachments := make([]slack.Attachment, 0)

	// basic structure
	attachment := slack.Attachment{
		Title:      mr.Title,
		Color:      getColor(mr),
		AuthorName: mr.Author.Name,
		AuthorIcon: mr.Author.AvatarURL,
		FooterIcon: mr.Author.AvatarURL,
		Fields: []slack.AttachmentField{
			{
				Title: "Data",
				Value: mr.UpdatedAt.Format("2 Jan 2006"),
				Short: true,
			},
			{
				Title: "Projeto",
				Value: proj.Name,
				Short: true,
			},
		},
		Actions: []slack.AttachmentAction{
			{
				Type:  "button",
				Name:  "repo",
				Text:  getButtonText(mr),
				URL:   mr.WebURL,
				Style: "primary",
			},
		},
	}

	if mr.Upvotes > 0 {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Aprovações",
			Value: ":thumbsup: " + strconv.Itoa(mr.Upvotes),
			Short: true,
		})
	}

	if mr.Downvotes > 0 {
		attachment.Fields = append(attachment.Fields, slack.AttachmentField{
			Title: "Rejeições",
			Value: ":thumbsdown: " + strconv.Itoa(mr.Upvotes),
			Short: true,
		})
	}

	attachments = append(attachments, attachment)

	attachments = append(attachments, slack.Attachment{
		Title: "Branches",
		Color: "#3AA3E3",
		Actions: []slack.AttachmentAction{
			{
				Type:  "button",
				Name:  "repo",
				Text:  mr.SourceBranch,
				URL:   mr.WebURL,
				Style: "primary",
			},
			{
				Type:  "button",
				Name:  "repo",
				Text:  mr.TargetBranch,
				URL:   mr.WebURL,
				Style: "danger",
			},
		},
	})

	return attachments
}

func getColor(mr *gitlab.MergeRequest) string {
	if mr.Upvotes >= APPROVAL_UPVOTES {
		return "#00FF00"
	} else {
		return "#FFFFFF"
	}
}

func getButtonText(mr *gitlab.MergeRequest) string {
	if mr.Upvotes >= APPROVAL_UPVOTES {
		return "Fazer o merge"
	} else {
		return "Revisar"
	}
}
