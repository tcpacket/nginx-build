package builder

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog/log"
)

type SourceArchive struct {
	Name      string `toml:"name"`
	Version   string `toml:"version"`
	URLPrefix string `toml:"url_prefix"`
}

type ComponentMap map[string]SourceArchive

var Components = ComponentMap{
	"nginx": {
		Name:      "nginx",
		Version:   "1.22.0",
		URLPrefix: "https://nginx.org/download",
	},
	"pcre": {
		Name:      "pcre",
		Version:   "8.32",
		URLPrefix: "https://ftp.exim.org/pub",
	},
	"openssl": {
		Name:      "openssl",
		Version:   "1.1.1q",
		URLPrefix: "https://www.openssl.org/source",
	},
	"libressl": {
		Name:      "libressl",
		Version:   "3.3.5",
		URLPrefix: "https://ftp.openbsd.org/pub/OpenBSD/LibreSSL",
	},
	"zlib": {
		Name:      "zlib",
		Version:   "1.2.12",
		URLPrefix: "https://zlib.net",
	},
	"openresty": {
		Name:      "openresty",
		Version:   "1.21.4.1",
		URLPrefix: "https://openresty.org/download",
	},
}

func (c ComponentMap) MarshalTOML() ([]byte, error) {
	return toml.Marshal(c)
}

func init() {
	var configPath string
	if home, err := os.UserHomeDir(); err != nil {
		return
	} else {
		configPath = home + "/.config/waf-builder/components.toml"
	}
	if _, err := os.Stat(configPath); err != nil {
		f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error().Err(err).Msg("failed to component config file file")
			return
		}
		defer f.Close()
		if err := toml.NewEncoder(f).Encode(Components); err != nil {
			log.Fatal().Err(err).Msg("failed to write component config file")
		}
		log.Debug().Str("caller", configPath).Msg("component config file created")
		return
	}
	f, err := os.Open(configPath)
	if err != nil {
		log.Warn().Err(err).Msg("failed to open component config file")
		return
	}
	defer f.Close()
	if err = toml.NewDecoder(f).Decode(&Components); err != nil {
		log.Warn().Err(err).Msg("failed to decode component config file")
		return
	}
	log.Info().Str("caller", configPath).Msgf("loaded %d component details from toml config",
		len(Components))
}

// component enumerations
const (
	ComponentNginx = iota
	ComponentOpenResty
	ComponentPcre
	ComponentOpenSSL
	ComponentLibreSSL
	ComponentZlib
	ComponentMax
)
