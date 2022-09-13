package main

import (
	"fmt"

	"github.com/tcpacket/nginx-build/builder"
)

func versionsGenNginx() []string {
	return []string{
		fmt.Sprintf("nginx-%s", builder.NginxVersion),
	}
}

func versionsGenOpenResty() []string {
	return []string{
		fmt.Sprintf("openresty-%s", builder.OpenRestyVersion),
	}
}

func versionsGenTengine() []string {
	return []string{
		fmt.Sprintf("tengine-%s", builder.TengineVersion),
	}
}

func printNginxVersions() {
	var versions []string
	versions = append(versions, versionsGenNginx()...)
	versions = append(versions, versionsGenOpenResty()...)
	versions = append(versions, versionsGenTengine()...)
	for _, v := range versions {
		fmt.Println(v)
	}
}

func versionCheck(version string) {
	if len(version) == 0 {
		log.Warn().Msgf("%v", "nginx version is not set.")
		log.Warn().Msgf("nginx-build use %s.", builder.NginxVersion)
	}
}
