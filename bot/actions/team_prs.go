package actions

import (
	"caco/services"
	"fmt"
	"strings"

	"github.com/google/go-github/v43/github"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

var TeamPRs_CommandDefinition = &slacker.CommandDefinition{
	Description: "Find PRs from the team",
	Example:     "prs",
	Handler:     handleTeamPrAction,
}

func handleTeamPrAction(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	proj := request.Param("project")

	if strings.TrimSpace(proj) == "" {
		response.Reply("Uh oh, nenhum projeto informado", slacker.WithThreadReply(true))
		return
	}

	svc := services.NewGithubService()

	// get all prs for projects
	//for _, p := range settings.Config.GithubProjects {
	prs, err := svc.GetOpenPullRequests(proj)

	if err != nil {
		// abort
		errResponse := err.(*github.ErrorResponse)
		if errResponse.Response.StatusCode == 404 {
			response.Reply("Na não. Esse projeto não existe. Você digitou corretamente?", slacker.WithThreadReply(true))
			return
		}
		response.ReportError(err, slacker.WithThreadError(true))
		return
	}

	for _, pr := range prs {
		go replyPR(response, pr)
	}
	// }
}

func replyPR(response slacker.ResponseWriter, pr *github.PullRequest) {
	// response.Reply("", slacker.WithAttachments(getPullRequestAttachment(pr)), slacker.WithThreadReply(true))
	response.Reply("", slacker.WithBlocks(getPullRequestBlocks(pr)), slacker.WithThreadReply(true))
}

func getPullRequestAttachment(pr *github.PullRequest) []slack.Attachment {
	attachments := make([]slack.Attachment, 0)

	// basic structure
	attachment := slack.Attachment{
		Title: *pr.Title,
		// Color:      getColor(pr),
		AuthorName: *pr.User.Login,
		AuthorIcon: *pr.User.AvatarURL,
		FooterIcon: *pr.User.AvatarURL,
		Fields: []slack.AttachmentField{
			{
				Title: "Data",
				Value: pr.UpdatedAt.Format("2 Jan 2006"),
				Short: true,
			},
			{
				Title: "Projeto",
				Value: *pr.Base.Repo.Name,
				Short: true,
			},
		},
		Actions: []slack.AttachmentAction{
			{
				Type:  "button",
				Name:  "repo",
				Text:  "Acessar",
				URL:   *pr.HTMLURL,
				Style: "primary",
			},
		},
	}

	// if mr.Upvotes > 0 {
	// 	attachment.Fields = append(attachment.Fields, slack.AttachmentField{
	// 		Title: "Aprovações",
	// 		Value: ":thumbsup: " + strconv.Itoa(mr.Upvotes),
	// 		Short: true,
	// 	})
	// }

	// if mr.Downvotes > 0 {
	// 	attachment.Fields = append(attachment.Fields, slack.AttachmentField{
	// 		Title: "Rejeições",
	// 		Value: ":thumbsdown: " + strconv.Itoa(mr.Upvotes),
	// 		Short: true,
	// 	})
	// }

	attachments = append(attachments, attachment)

	attachments = append(attachments, slack.Attachment{
		Title: "Branches",
		Color: "#3AA3E3",
		Actions: []slack.AttachmentAction{
			{
				Type:  "button",
				Name:  "repo",
				Text:  *pr.Head.Ref,
				Style: "primary",
			},
			{
				Type:  "button",
				Name:  "repo",
				Text:  *pr.Base.Ref,
				Style: "danger",
			},
		},
	})

	return attachments
}

func getPullRequestBlocks(pr *github.PullRequest) []slack.Block {
	blocks := []slack.Block{}

	title := slack.NewTextBlockObject("plain_text", *pr.Title, true, false)
	if pr.GetDraft() {
		title.Text = "DRAFT: " + title.Text
	}
	headerBl := slack.NewHeaderBlock(title)
	blocks = append(blocks, headerBl)

	sectionText := fmt.Sprintf(":github: *projeto:* <%s|%s> \n\n :pr-merged: *%s* :right-arrow: *%s*", pr.Base.Repo.GetHTMLURL(), *pr.Base.Repo.Name, *pr.Head.Ref, *pr.Base.Ref)
	section := slack.SectionBlock{
		Type: slack.MBTSection,
		Text: &slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: sectionText,
		},
	}
	blocks = append(blocks, section)

	if !pr.GetDraft() {
		btnTxt := slack.NewTextBlockObject("plain_text", "Revisar esse PR", false, false)
		nextBtn := slack.NewButtonBlockElement("", "click_me_123", btnTxt)
		nextBtn.Style = "primary"
		nextBtn.URL = *pr.HTMLURL
		actionBlock := slack.NewActionBlock("", nextBtn)
		blocks = append(blocks, actionBlock)
	}

	imageEl := slack.NewImageBlockElement(*pr.User.AvatarURL, *pr.User.Login)
	authorEl := &slack.TextBlockObject{
		Type: slack.MarkdownType,
		Text: fmt.Sprintf("*<%s|%s>* criou esse pull request em *%s*", *pr.User.HTMLURL, *pr.User.Login, pr.CreatedAt.Format("2 Jan 2006")),
	}
	// metaEl := &slack.TextBlockObject{
	// 	Type: slack.MarkdownType,
	// 	Text: fmt.Sprintf("*%d* commits: *%d* adições, *%d* remoções, *%d* arquivos alterados", pr.GetCommits(), pr.GetAdditions(), pr.GetDeletions(), pr.GetChangedFiles()),
	// }
	// meta_reviewEl := &slack.TextBlockObject{
	// 	Type: slack.MarkdownType,
	// 	Text: fmt.Sprintf("*%d* comentários de revisões", pr.GetReviewComments()),
	// }
	// contextBlock := slack.NewContextBlock("contextBl", imageEl, authorEl, metaEl, meta_reviewEl)
	contextBlock := slack.NewContextBlock("contextBl", imageEl, authorEl)
	blocks = append(blocks, contextBlock)

	blocks = append(blocks, slack.NewDividerBlock())
	return blocks
}
