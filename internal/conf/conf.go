package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type versionEndpointConf struct {
	GoroutineCount bool `toml:"goroutine_count"`
	GoVersion      bool `toml:"go_version"`
	OS             bool `toml:"os"`
	Arch           bool `toml:"arch"`
	ServerTime     bool `toml:"server_time"`
	Uptime         bool `toml:"uptime"` // TODO: show the uptime. For now, it does nothing.
}

// Conf is a struct that holds the fields for configuring a Miresa server.
type Conf struct {
	Title           string              `toml:"title"`
	Port            int                 `toml:"port"`
	IDLength        int                 `toml:"id_length"`
	EasterEggs      bool                `toml:"enable_eastereggs"`
	LogType         string              `toml:"log_type"`
	DatabaseURL     string              `toml:"database_url"`
	VersionEndpoint versionEndpointConf `toml:"version_endpoint"`
}

// Config is the default Conf. Once LoadConf has been called, its values are
// filled with the server's config.
var Config = Conf{
	Title:       "miresa",
	Port:        8000,
	IDLength:    64,
	EasterEggs:  true,
	LogType:     "slog",
	DatabaseURL: os.Getenv("DATABASE_URL"),
	VersionEndpoint: versionEndpointConf{
		GoroutineCount: true,
		GoVersion:      true,
		OS:             true,
		Arch:           true,
		ServerTime:     true,
		Uptime:         true,
	},
}

// LoadConf loads configuration from the user's config directory.
func LoadConf() (Conf, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return Config, fmt.Errorf("failed to get user config dir: %v", err)
	}

	// For now, config can only be done in TOML.
	// TODO: check for miresa-srv.{json,yaml,yml,toml} and use those if
	// they exist.
	// TODO: use the config file from the CLI flags provided.
	confFile := filepath.Join(confDir, "miresa-srv.toml")
	data, err := os.ReadFile(confFile)
	if err != nil {
		return Config, fmt.Errorf("failed to read %s: %v", confFile, err)
	}

	err = toml.Unmarshal(data, &Config)
	if err != nil {
		return Config, fmt.Errorf("failed to unmarshal %s: %v", confFile, err)
	}
	return Config, nil
}
