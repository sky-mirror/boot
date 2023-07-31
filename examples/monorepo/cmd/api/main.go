package main

import (
	"context"

	"github.com/sky-mirror/boot"
	"github.com/sky-mirror/boot/examples/monorepo/logger"
	"github.com/sky-mirror/boot/examples/monorepo/slack"
	"github.com/urfave/cli/v2"
)

var alertChan = slack.NewConfig("alert")
var infoChan = slack.NewConfig("info")

// Main is the real app entrypoint.
func Main(ctx context.Context, c *cli.Context) {
	logger := logger.Default()
	logger.Println("starting app")

	alertSlacker := slack.NewWebhook(alertChan)
	infoSlacker := slack.NewWebhook(infoChan)

	infoSlacker.PostMessage("app started")

	<-ctx.Done()

	alertSlacker.PostMessage("app terminating")
}

func main() {
	app := boot.App{
		Main: Main,
	}

	app.Run()
}
