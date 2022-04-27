package settings

import (
	"log"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"
)

type config struct {
	Environment                string   `env:"ENVIRONMENT" default:"development"`
	Version                    string   `env:"VERSION" required:"true"`
	Debug                      bool     `env:"DEBUG" default:"false"`
	SlackBotToken              string   `env:"SLACK_BOT_TOKEN" required:"true"`
	SlackAppToken              string   `env:"SLACK_APP_TOKEN" required:"true"`
	GithubPersonalToken        string   `env:"GITHUB_PERSONAL_TOKEN"`
	GithubProjects             []string `env:"GITHUB_PROJECTS"`
	GithubOrg                  string   `env:"GITHUB_ORG"`
	GitlabTeam                 string   `env:"GITLAB_TEAM" default:"squad-estoque"`
	GitlabToken                string   `env:"GITLAB_TOKEN" required:"true"`
	GitlabProjectExclusions    []int    `env:"GITLAB_PROJECT_EXCLUISIONS"`
	GitlabProjectExclusionsMap map[int]bool
	GitlabURL                  string `env:"GITLAB_URL"`
	GoogleCredentialsPath      string `env:"GOOGLE_APPLICATION_CREDENTIALS"`
	GoogleProjectID            string `env:"GOOGLE_PROJEC_ID"`
}

// Config ...
var Config config

// InitConfigs inicializa as configurações de ambiente
func InitConfigs() {
	// carrega o .env (se existir)
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	// bind env vars
	if err := env.Set(&Config); err != nil {
		log.Fatal(err)
	}

	// transform exclusion array to map
	Config.GitlabProjectExclusionsMap = make(map[int]bool)
	for _, item := range Config.GitlabProjectExclusions {
		Config.GitlabProjectExclusionsMap[item] = true
	}
}
