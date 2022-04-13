package main

import (
	"caco/bot"
	"caco/settings"
)

func main() {
	// inicializa a configuração
	settings.InitConfigs()

	bot.InitBot()
}
