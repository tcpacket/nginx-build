package configure

import (
	"fmt"
	"strings"

	"github.com/tcpacket/nginx-build/builder"
	"github.com/tcpacket/nginx-build/modules"
)

const bashBreak = " \\\n"

type Conf struct {
	Mods []modules.Module
	Deps []builder.StaticLibrary
	Opts Options
	Dir  string

	PathPrefix    string
	NginxConfPath string
	HTTPLogPath   string
	ErrorLogPath  string

	OpenRestyEnabled bool
	HTTP2Enabled     bool
	WithCompat       bool
}

func br(s *strings.Builder) {
	s.WriteString(bashBreak)
}

func wr(s *strings.Builder, in string) {
	s.WriteString(in)
}

func ln(s *strings.Builder, in string) {
	wr(s, in)
	br(s)
}

func (c *Conf) Generate(initial string, Concurrency int) string {
	cf := &strings.Builder{}
	openSSLStatic := false
	if len(initial) == 0 {
		ln(cf, `#!/bin/sh

./configure \
`)
	} else {
		wr(cf, initial)
	}
	if c.OpenRestyEnabled {
		ln(cf, fmt.Sprintf("-j%d", Concurrency))
	}
	for _, d := range c.Deps {
		wr(cf, d.Option)
		wr(cf, "=../")
		wr(cf, d.Name)
		wr(cf, "-")
		ln(cf, d.Version)
		if d.Name == "openssl" || d.Name == "libressl" {
			openSSLStatic = true
		}
	}
	if openSSLStatic && !strings.Contains(initial, "--with-http_ssl_module") {
		ln(cf, "--with-http_ssl_module")
	}
	if c.WithCompat {
		ln(cf, "--with-compat")
	}
	if c.PathPrefix != "" {
		wr(cf, "--prefix=")
		ln(cf, c.PathPrefix)
	}
	if c.HTTPLogPath != "" {
		wr(cf, "--http-log-path=")
		ln(cf, c.HTTPLogPath)
	}
	if c.ErrorLogPath != "" {
		wr(cf, "--error-log-path=")
		ln(cf, c.ErrorLogPath)
	}
	if c.HTTP2Enabled {
		ln(cf, "--with-http_v2_module")
	}
	wr(cf, handleModuleFlags(c.Mods))
	for _, option := range c.Opts.Values {
		if *option.Value == "" {
			continue
		}
		switch {
		case option.Name == "--add-module":
			wr(cf, normalizeAddModulePaths(*option.Value, c.Dir, false))
		case option.Name == "--add-dynamic-module":
			wr(cf, normalizeAddModulePaths(*option.Value, c.Dir, true))
		case strings.Contains(*option.Value, " "):
			wr(cf, option.Name)
			wr(cf, "='")
			wr(cf, *option.Value)
			ln(cf, "'")
		default:
			wr(cf, option.Name)
			wr(cf, "=")
			ln(cf, *option.Value)
		}
	}
	for _, option := range c.Opts.Bools {
		if *option.Enabled {
			ln(cf, option.Name)
		}
	}
	return cf.String()
}

func handleModuleFlags(mods []modules.Module) string {
	var res = &strings.Builder{}
	for _, m := range mods {
		switch m.Dynamic {
		case true:
			res.WriteString("--add-dynamic-module=")
		default:
			res.WriteString("--add-module=")
		}
		switch m.Form {
		case "local":
			res.WriteString(m.URL)
		default:
			res.WriteString("../")
			res.WriteString(m.Name)
			res.WriteString(bashBreak)
		}
	}
	return res.String()
}
