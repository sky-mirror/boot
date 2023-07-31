package slack

import (
	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

// Config is the config to create webhook object.
type Config struct {
	WebhookURL string
	Channel    string

	prefix string
}

// NewConfig creates the new config with prefix and register itself.
func NewConfig(prefix string) *Config {
	cfg := &Config{prefix: prefix}
	boot.Register(cfg)
	return cfg
}

func (c *Config) addPrefix(s string) string {
	if c.prefix == "" {
		return s
	}

	return c.prefix + "-" + s
}

// CliFlags implements the CliFlagers interface.
func (c *Config) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.StringFlag{
		Name:        c.addPrefix("slack-webhook-url"),
		Destination: &c.WebhookURL,
		Required:    true,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        c.addPrefix("slack-channel"),
		Destination: &c.Channel,
		Required:    true,
	})

	return flags
}
