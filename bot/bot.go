package bot

import (
	"caco/bot/actions"
	"caco/settings"
	"context"
	"log"

	"github.com/shomali11/slacker"
)

// InitBot ...
func InitBot() {
	if settings.Config.Debug {
		log.Println("Iniciando bot")
	}

	// make the actions available
	println("starting bot service")
	bot := slacker.NewClient(settings.Config.SlackBotToken, settings.Config.SlackAppToken, slacker.WithDebug(settings.Config.Debug))

	bot.DefaultCommand(messageHandler)

	bot.Command("prs <project>", actions.TeamPRs_CommandDefinition)

	bot.Err(func(err string) {
		log.Println(err)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func messageHandler(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	// event := botCtx.Event()
	// result, err := services.DetectIntent(event.User, event.Text)
	// if err != nil {
	// 	response.Reply(err.Error(), slacker.WithThreadReply(true))
	// }

	// action := actions.GetAction(result.GetAction())
	// if action != nil {
	// 	action.Command(result, request, response)
	// 	return
	// }

	// log.Println("não encontrou action")

	response.Reply("Oh oh... um tanto quanto perdido aqui!", slacker.WithThreadReply(true))
}
