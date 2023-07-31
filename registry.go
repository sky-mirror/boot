package boot

import (
	"log"

	"github.com/urfave/cli/v2"
)

// DebugMode indicates whether to show more message.
var DebugMode = true

// CliFlager is interface to describe a struct
// which is a set of options to singleton in package.
// This struct has method CliFlags to returns the options this package needed.
type CliFlager interface {
	CliFlags() []cli.Flag
}

// Beforer is interface for some package may needs an before hook to init
// or validates before calls main function.
type Beforer interface {
	Before(*cli.Context) error
}

// Afterer is interface for some package may needs an after hook to destroy
// or cleanup after main function exited.
type Afterer interface {
	After()
}

var cliFlagers []CliFlager

// Register registers CliFlager, so we won't use a package without
// init it.
func Register(f CliFlager) {
	cliFlagers = append(cliFlagers, f)
}

// IsBeforer checks interface conversion.
func IsBeforer(Beforer) {}

// IsAfterer checks interface conversion.
func IsAfterer(Afterer) {}

// Flags returns flags from all registered packages.
func Flags() []cli.Flag {
	var res []cli.Flag
	for _, f := range cliFlagers {
		res = append(res, f.CliFlags()...)
	}

	return res
}

// Initialize inits all registered packages.
func Initialize(c *cli.Context) error {
	for _, f := range cliFlagers {
		b, ok := f.(Beforer)
		if ok {
			if DebugMode {
				log.Printf("running %T", b)
			}

			err := b.Before(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Finalize finalizes registered packages, its execution order is reversed.
func Finalize(c *cli.Context) error {
	for i := len(cliFlagers) - 1; i >= 0; i-- {
		f := cliFlagers[i]
		a, ok := f.(Afterer)
		if ok {
			if DebugMode {
				log.Printf("running %T", a)
			}

			a.After()
		}
	}

	return nil
}
