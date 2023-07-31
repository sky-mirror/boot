package boot

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

// App is cli wrapper that do some common operation and creates signal handler.
type App struct {
	Flags []cli.Flag
	Main  func(ctx context.Context, c *cli.Context)
}

func (a *App) before(c *cli.Context) (err error) {
	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered: ", r)
			if DebugMode {
				debug.PrintStack()
			}
			err = errors.New("init failed")
		}
	}()

	return Initialize(c)
}

func (a *App) after(c *cli.Context) error {
	return Finalize(c)
}

func (a *App) wrapMain(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("\nReceives signal: %v\n", sig)
		cancel()
	}()

	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Main recovered: ", r)
			if DebugMode {
				debug.PrintStack()
			}
		}
	}()

	a.Main(ctx, c)
	time.Sleep(3 * time.Second)
	log.Println("terminated")

	return nil
}

// GetApp gets the internal app instance.
func (a *App) GetApp() *cli.App {
	app := cli.NewApp()
	app.Flags = append(a.Flags, Flags()...)
	app.Before = a.before
	app.After = a.after
	app.Action = a.wrapMain

	return app
}

// Run setups everything and runs Main.
func (a *App) Run() {
	// maxprocs.Set(maxprocs.Logger(log.Printf), maxprocs.Min(4))
	app := a.GetApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
