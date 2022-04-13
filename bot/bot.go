package bot

import (
	"caco/bot/actions"
	"caco/services"
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
	bot := slacker.NewClient(settings.Config.SlackToken, slacker.WithDebug(settings.Config.Debug))
	bot.DefaultCommand(messageHandler)

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

func messageHandler(request slacker.Request, response slacker.ResponseWriter) {
	result, err := services.DetectIntent(request.Event().User, request.Event().Text)
	if err != nil {
		response.Reply(err.Error(), slacker.WithThreadReply(true))
	}

	action := actions.GetAction(result.GetAction())
	if action != nil {
		action.Command(result, request, response)
		return
	}

	log.Println("não encontrou action")

	response.Reply("Ohoh... um tanto quanto perdido aqui!", slacker.WithThreadReply(true))
}
