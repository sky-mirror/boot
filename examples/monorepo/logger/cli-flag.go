package logger

import (
	"errors"
	"os"
	"path"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

var defaultCfg config

func init() {
	boot.Register(&defaultCfg)
	boot.IsBeforer(&defaultCfg)
	boot.IsAfterer(&defaultCfg)
}

type fileConfig struct {
	enabled bool
	dir     string
	name    string
}

type config struct {
	file fileConfig
}

func (cfg *config) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.BoolFlag{
		Name:        "log-enable-file",
		EnvVars:     []string{"LOG_ENABLE_FILE"},
		Destination: &cfg.file.enabled,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "log-file-name",
		EnvVars:     []string{"LOG_FILE_NAME"},
		Usage:       "filename prefix of log file",
		Value:       path.Base(os.Args[0]),
		Destination: &cfg.file.name,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "log-file-dir",
		EnvVars:     []string{"LOG_FILE_DIR"},
		Usage:       "path of log file",
		Value:       os.TempDir(),
		Destination: &cfg.file.dir,
	})

	return flags
}

func (cfg *config) Before(c *cli.Context) error {
	if cfg.file.enabled {
		if len(cfg.file.dir) == 0 {
			return errors.New("log-file-dir must be set")
		}
		if len(cfg.file.name) == 0 {
			return errors.New("log-file-name must be set")
		}
	}

	return Initialize()
}

func (cfg *config) After() {
	Finalize()
}
