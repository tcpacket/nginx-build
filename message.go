package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

var (
	NginxBuildVersion string
)

func nginxBuildVersion() string {
	if NginxBuildVersion != "" {
		return NginxBuildVersion
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}

func printNginxBuildVersion() {
	fmt.Printf(`nginx-build %s
Compiler: %s %s
Copyright (C) 2014-2022 Tatsuhiko Kubo <tcpacket@gmail.com>
`,
		nginxBuildVersion(),
		runtime.Compiler,
		runtime.Version())

}

func printConfigureOptions() error {
	cmd := exec.Command("objs/nginx", "-V")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func printFirstMsg() {
	fmt.Printf(`nginx-build: %s
Compiler: %s %s
`,
		nginxBuildVersion(),
		runtime.Compiler,
		runtime.Version())
}

func printLastMsg(workDir, srcDir string, openResty, configureOnly bool) {
	log.Info().Msgf("%v", "Complete building nginx!")

	if !openResty {
		if !configureOnly {
			fmt.Println()
			err := printConfigureOptions()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	fmt.Println()

	lastMsgFormat := `Enter the following command for install nginx.

   $ cd %s/%s%s
   $ sudo make install
`
	if configureOnly {
		log.Printf(lastMsgFormat, workDir, srcDir, "\n   $ make")
	} else {
		log.Printf(lastMsgFormat, workDir, srcDir, "")
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.VisitAll(func(f *flag.Flag) {
		if !isNginxBuildOption(f.Name) {
			return
		}
		s := fmt.Sprintf("  -%s", f.Name)
		s += "\n\t"
		s += f.Usage
		defValue := defaultStringValue(f.Name)
		if defValue != "" {
			s += fmt.Sprintf(" ( default: %s )", defValue)
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s\n", s)
	})
}
