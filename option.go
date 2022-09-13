package main

import (
	"runtime"
	"strconv"

	"github.com/tcpacket/nginx-build/builder"
)

type Options struct {
	Values  map[string]OptionValue
	Bools   map[string]OptionBool
	Numbers map[string]OptionNumber
}

type OptionValue struct {
	Desc    string
	Value   *string
	Default string
}

type OptionBool struct {
	Desc    string
	Enabled *bool
	// all options are false by default
	// Default bool
}

type OptionNumber struct {
	Desc    string
	Value   *int
	Default int
}

func makeNginxBuildOptions() Options {
	var nginxBuildOptions Options
	argsNumber := make(map[string]OptionNumber)
	argsBool := make(map[string]OptionBool)
	argsString := make(map[string]OptionValue)

	argsNumber["j"] = OptionNumber{
		Desc:    "jobs to build nginx",
		Default: runtime.NumCPU(),
	}

	argsBool["verbose"] = OptionBool{
		Desc: "verbose mode",
	}
	argsBool["pcre"] = OptionBool{
		Desc: "embedded PCRE staticlibrary",
	}
	argsBool["openssl"] = OptionBool{
		Desc: "embedded OpenSSL staticlibrary",
	}
	argsBool["libressl"] = OptionBool{
		Desc: "embedded LibreSSL staticlibrary",
	}
	argsBool["zlib"] = OptionBool{
		Desc: "embedded zlib staticlibrary",
	}
	argsBool["clear"] = OptionBool{
		Desc: "remove entries in working directory",
	}
	argsBool["version"] = OptionBool{
		Desc: "print nginx-build version",
	}
	argsBool["versions"] = OptionBool{
		Desc: "print nginx versions",
	}
	argsBool["openresty"] = OptionBool{
		Desc: "download openresty instead of nginx",
	}
	argsBool["tengine"] = OptionBool{
		Desc: "download tengine instead of nginx",
	}
	argsBool["configureonly"] = OptionBool{
		Desc: "configure nginx only not building",
	}
	argsBool["idempotent"] = OptionBool{
		Desc: "build nginx if already installed nginx version is different",
	}
	argsBool["help-all"] = OptionBool{
		Desc: "print all flags",
	}
	argsBool["with-compat"] = OptionBool{
		Desc: "enable extra nginx module compatibility",
	}
	argsBool["with-http_v2_module"] = OptionBool{
		Desc: "enable HTTP/2 module",
	}

	argsString["v"] = OptionValue{
		Desc:    "nginx version",
		Default: builder.NginxVersion,
	}
	argsString["c"] = OptionValue{
		Desc:    "configuration file for building nginx",
		Default: "",
	}
	argsString["m"] = OptionValue{
		Desc:    "configuration file for 3rd party modules",
		Default: "",
	}
	argsString["d"] = OptionValue{
		Desc:    "working directory",
		Default: "",
	}
	argsString["pcreversion"] = OptionValue{
		Desc:    "PCRE version",
		Default: builder.PcreVersion,
	}
	argsString["opensslversion"] = OptionValue{
		Desc:    "OpenSSL version",
		Default: builder.OpenSSLVersion,
	}
	argsString["libresslversion"] = OptionValue{
		Desc:    "LibreSSL version",
		Default: builder.LibreSSLVersion,
	}
	argsString["zlibversion"] = OptionValue{
		Desc:    "zlib version",
		Default: builder.ZlibVersion,
	}
	argsString["openrestyversion"] = OptionValue{
		Desc:    "openresty version",
		Default: builder.OpenRestyVersion,
	}
	argsString["tengineversion"] = OptionValue{
		Desc:    "tengine version",
		Default: builder.TengineVersion,
	}
	argsString["patch"] = OptionValue{
		Desc:    "patch path for applying to nginx",
		Default: "",
	}
	argsString["patch-opt"] = OptionValue{
		Desc:    "option for patch",
		Default: "",
	}
	argsString["conf-path"] = OptionValue{
		Desc:    "nginx configuration file path",
		Default: "/etc/nginx/nginx.conf",
	}
	argsString["http-log-path"] = OptionValue{
		Desc:    "nginx http log file path",
		Default: "/var/log/nginx/access.log",
	}
	argsString["error-log-path"] = OptionValue{
		Desc:    "nginx error log file path",
		Default: "/var/log/nginx/error.log",
	}
	argsString["prefix"] = OptionValue{
		Desc:    "nginx installation prefix",
		Default: "/usr/local",
	}
	/*	argsString["pid-path"] = OptionValue{
		Desc:    "nginx pid file path",
		Default: "/var/run/nginx.pid",
	}*/

	nginxBuildOptions.Bools = argsBool
	nginxBuildOptions.Values = argsString
	nginxBuildOptions.Numbers = argsNumber

	return nginxBuildOptions
}

func isNginxBuildOption(k string) bool {
	if _, ok := nginxBuildOptions.Bools[k]; ok {
		return true
	}
	if _, ok := nginxBuildOptions.Values[k]; ok {
		return true
	}
	if _, ok := nginxBuildOptions.Numbers[k]; ok {
		return true
	}
	return false
}

func defaultStringValue(k string) string {
	if _, ok := nginxBuildOptions.Values[k]; ok {
		return nginxBuildOptions.Values[k].Default
	}
	if _, ok := nginxBuildOptions.Numbers[k]; ok {
		return strconv.Itoa(nginxBuildOptions.Numbers[k].Default)
	}
	return ""
}
