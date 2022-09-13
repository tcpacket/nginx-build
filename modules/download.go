package modules

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tcpacket/nginx-build/command"
	"github.com/tcpacket/nginx-build/logger"
	"github.com/tcpacket/nginx-build/util"

	"github.com/go-git/go-git/v5"
)

var log = logger.Get()

func DownloadAndExtractParallel(m Module) {
	if util.FileExists(m.Name) {
		log.Printf("%s already exists.", m.Name)
		return
	}

	if m.Form != "local" {
		if len(m.Rev) > 0 {
			log.Printf("Download %s-%s.....", m.Name, m.Rev)
		} else {
			log.Printf("Download %s.....", m.Name)
		}

		logName := fmt.Sprintf("%s.log", m.Name)

		err := download(m, logName)
		if err != nil {
			util.PrintFatalMsg(err, logName)
		}
	} else if !util.FileExists(m.URL) {
		log.Fatal().Msgf("no such directory: %s", m.URL)
	}
}

func download(m Module, logName string) error {
	form := m.Form

	f, err := os.Create(logName)
	if err != nil {
		log.Fatal().Msgf("creating log file %s failed: %s", logName, err.Error())
	}
	defer util.Fclose(f)

	switch form {
	case "git":
		_, err = git.PlainClone(m.Name, false, &git.CloneOptions{
			URL:               m.URL,
			SingleBranch:      true,
			RecurseSubmodules: 10,
			Progress:          f,
			InsecureSkipTLS:   false,
		})
		return err
	case "hg":
		args := []string{form, "clone", m.URL}
		if command.VerboseEnabled {
			err = command.Run(args)
			break
		}
		cmd, err := command.Make(args)
		if err != nil {
			break
		}
		writer := bufio.NewWriter(f)
		cmd.Stderr = writer
		err = cmd.Run()
		util.Flush(writer)
	case "local": // not implemented yet
		return nil
	default:
		err = fmt.Errorf("form=%s is not supported", form)
	}
	return err
}
