package main

import (
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/api"
	"github.com/vseriousv/go_diffie_hellman_chat_server/internal/config"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	defer log.Sync()

	c := config.DefaultConfig()
	app := api.AppStruct{}
	app.RunApp(c, log)
}
