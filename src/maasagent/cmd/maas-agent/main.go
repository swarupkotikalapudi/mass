package main

/*
	Copyright 2023 Canonical Ltd.  This software is licensed under the
	GNU Affero General Public License version 3 (see the file LICENSE).
*/

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"gopkg.in/yaml.v3"
	"launchpad.net/maas/maas/src/maasagent/internal/workflow"
	wflog "launchpad.net/maas/maas/src/maasagent/internal/workflow/log"
	"launchpad.net/maas/maas/src/maasagent/pkg/workflow/codec"
)

const (
	TemporalPort = 5271
)

// config represents a neccessary set of configuration options for MAAS Agent
type config struct {
	SystemID    string   `yaml:"system_id"`
	Secret      string   `yaml:"secret"`
	Controllers []string `yaml:"controllers,flow"`
}

func Run() int {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if logLevel, err := zerolog.ParseLevel(envLogLevel); err != nil {
			log.Warn().Str("LOG_LEVEL", envLogLevel).Msg("Unknown log level, defaulting to INFO")
		} else {
			zerolog.SetGlobalLevel(logLevel)
		}
	}

	cfg, err := getConfig()
	if err != nil {
		log.Error().Err(err).Send()
		return 1
	}

	// Encryption Codec required for Temporal Workflow's payload encoding
	codec, err := codec.NewEncryptionCodec([]byte(cfg.Secret))
	if err != nil {
		log.Error().Err(err).Msg("Encryption codec setup failed")
		return 1
	}

	clientBackoff := backoff.NewExponentialBackOff()
	clientBackoff.MaxElapsedTime = 60 * time.Second

	client, err := backoff.RetryWithData(
		func() (client.Client, error) {
			return client.Dial(client.Options{
				// TODO: fallback retry if Controllers[0] is unavailable
				HostPort: fmt.Sprintf("%s:%d", cfg.Controllers[0], TemporalPort),
				Logger:   wflog.New(log.Logger),
				DataConverter: converter.NewCodecDataConverter(
					converter.GetDefaultDataConverter(),
					codec,
				),
			})
		}, clientBackoff,
	)

	if err != nil {
		log.Error().Err(err).Msg("Temporal client error")
		return 1
	}

	workerPoolBackoff := backoff.NewExponentialBackOff()
	workerPoolBackoff.MaxElapsedTime = 60 * time.Second

	_, err = backoff.RetryWithData(
		func() (*workflow.WorkerPool, error) {
			return workflow.NewWorkerPool(cfg.SystemID, client)
		}, workerPoolBackoff,
	)

	if err != nil {
		log.Error().Err(err).Msg("Temporal worker pool creation failure")
		return 1
	}

	log.Info().Msg("Service MAAS Agent started")

	sigC := make(chan os.Signal, 2)

	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)

	<-sigC

	return 0
}

// getConfig reads MAAS Agent YAML configuration file
// TODO: agent.yaml config is generated by rackd, however this behaviour
// should be changed when MAAS Agent will be a standalone service, not managed
// by the Rack Controller.
func getConfig() (*config, error) {
	fname := os.Getenv("MAAS_AGENT_CONFIG")
	if fname == "" {
		fname = "/etc/maas/agent.yaml"
	}

	data, err := os.ReadFile(filepath.Clean(fname))
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	cfg := &config{}

	err = yaml.Unmarshal([]byte(data), cfg)
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	return cfg, nil
}

func main() {
	os.Exit(Run())
}
