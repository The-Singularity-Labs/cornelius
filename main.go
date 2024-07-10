package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/the-singularity-labs/cornelius/log"
	"github.com/the-singularity-labs/cornelius/sync"

	"github.com/hoenirvili/skapt"
	"github.com/hoenirvili/skapt/argument"
	"github.com/hoenirvili/skapt/flag"
)

const DefaultWalletPath = "/etc/cornelius/arweave_wallet.json"

type Synchronizer interface {
	Start(context.Context) error
}

func main() {
	app := skapt.Application{
		Name:        "Cornelius",
		Description: "Sync Object Storage Objects to Arweave",
		Version:     "0.0.1",
		Handler: func(scaptCtx *skapt.Context) error {
			// ctx, _ := context.WithCancel(context.Background())
			ctx := context.Background()

			logLevel := slog.LevelInfo
			if scaptCtx.Bool("debug") {
				logLevel = slog.LevelDebug
			}

			var logger log.Logger
			if logType := scaptCtx.String("logtype"); logType == "json" {
				logger = log.NewJsonLogger(logLevel)
			} else if logType == "text" {
				logger = log.NewTextLogger(logLevel)
			} else {
				return fmt.Errorf("%q is not a valid log type", logType)
			}

			logger.Info("loading configuration")

			ardrivecliPath := scaptCtx.String("ardrivecli")
			configPath := scaptCtx.String("config")
			config, err := sync.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("unable to load config: %w", err)
			}

			logger.Info("initializing synchronizer")
			var synchronizer Synchronizer = sync.New(logger, ardrivecliPath, config)
			if err != nil {
				return fmt.Errorf("unable to initialize synchronizer: %w", err)
			}

			err = synchronizer.Start(ctx)
			if err != nil {
				return fmt.Errorf("unable to synchronize: %w", err)
			}

			return nil
		},
		Flags: flag.Flags{
			flag.Flag{
				Short: "c", Long: "config",
				Description: "Filepath of YAML config file",
				Type:        argument.String,
				Required:    true,
			},
			flag.Flag{
				Short: "x", Long: "ardrivecli",
				Description: "Filepath of YAML config file",
				Type:        argument.String,
				Required:    false,
			},
			flag.Flag{
				Short: "d", Long: "debug",
				Description: "Filepath of YAML config file",
				Type:        argument.Bool,
				Required:    false,
			},
			flag.Flag{
				Short: "l", Long: "logtype",
				Description: "Type of logger to use. Can be text or json",
				Type:        argument.String,
				Required:    false,
			},
		},
	}
	app.Exec(os.Args)
}
