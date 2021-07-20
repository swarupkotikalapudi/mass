package config

import (
	"context"
	"errors"
	"path/filepath"
)

var (
	ErrConfigFileNotDefined = errors.New("config file not defined")
	ErrBadTlsConfig         = errors.New("invalid TLS config")
)

type ConfigKey string
type ctxFilename struct{}

type MetricsConfig struct {
	Bind   string `yaml:"bind"`
	Port   int    `yaml:"port"`
	CACert string `yaml:"ca_cert,omitempty"`
	Cert   string `yaml:"cert,omitempty"`
	Key    string `yaml:"priv_key,omitempty"`
}

type TlsConfig struct {
	CACert      string `yaml:"ca_cert,omitempty"`
	Cert        string `yaml:"cert,omitempty"`
	Key         string `yaml:"priv_key,omitempty"`
	SkipCaCheck bool   `yaml:"insecure_ca,omitempty"`
}

type RackConfig struct {
	BasePath       string        `yaml:"-"`
	SystemID       string        `yaml:"-"`
	Secret         string        `yaml:"-"`
	MaasUrl        StringArray   `yaml:"maas_url,flow"`
	TftpRoot       string        `yaml:"tftp_root,omitempty"`
	TftpPort       int           `yaml:"tftp_port,omitempty"`
	ClusterUUID    string        `yaml:"cluster_uuid,omitempty"`
	Debug          bool          `yaml:"debug,omitempty"`
	Metrics        MetricsConfig `yaml:"metrics,omitempty"`
	Tls            TlsConfig     `yaml:"rpc,omitempty"`
	SupervisordURL string        `yaml:"supervisord,omitempty"`
	Proxy          bool          `yaml:"proxy"`
	NTPBindAddr    string        `yaml:"ntp_bind_addr"`
	NTPRefreshRate int           `yaml:"ntp_refresh_rate"`
}

const (
	systemIDFile ConfigKey = "maas_id"
	secretFile   ConfigKey = "secret"
)

// Config global rackd configuration
var (
	Config *RackConfig = new()
)

func Load(ctx context.Context, filename string) (_ context.Context, err error) {
	err = load(ctx, filename)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, ctxFilename{}, filename)
	return ctx, nil
}

func Save(ctx context.Context) (err error) {
	fname, ok := ctx.Value(ctxFilename{}).(string)
	if !ok {
		return ErrConfigFileNotDefined
	}

	if err = setConfigToFile(systemIDFile, Config.SystemID); err != nil {
		return
	}
	if err = setConfigToFile(secretFile, Config.Secret); err != nil {
		return
	}

	return saveGlobal(ctx, fname)
}

func Reload(ctx context.Context) (err error) {
	fname, ok := ctx.Value(ctxFilename{}).(string)
	if !ok {
		return ErrConfigFileNotDefined
	}

	return load(ctx, fname)
}

func load(ctx context.Context, filename string) (err error) {
	err = loadGlobal(ctx, filename)
	if err != nil {
		return err
	}

	Config.SystemID, err = getConfigFromFile(systemIDFile)
	if err != nil {
		return err
	}

	Config.Secret, err = getConfigFromFile(secretFile)
	if err != nil {
		return err
	}
	return
}

func getAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(Config.BasePath, path)
	}
	return path
}

func GetTftpPath(path string) string {
	return filepath.Join(getAbsPath(Config.TftpRoot), path)
}

func SupervisordURL() string {
	return Config.SupervisordURL
}
