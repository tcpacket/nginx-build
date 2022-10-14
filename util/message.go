package util

import (
	"bufio"
	"os"

	"github.com/tcpacket/waf-builder/command"
	"github.com/tcpacket/waf-builder/logger"
)

var log = logger.Get()

func PrintFatalMsg(err error, path string) {
	if command.VerboseEnabled {
		log.Fatal().Msgf("%v", err)
	}

	f, err2 := os.Open(path)
	if err2 != nil {
		log.Printf("error-log: %s is not found\n", path)
		log.Fatal().Msgf("%v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		_, _ = os.Stderr.Write(scanner.Bytes())
		_, _ = os.Stderr.Write([]byte("\n"))
	}

	log.Fatal().Msgf("%v", err)
}
