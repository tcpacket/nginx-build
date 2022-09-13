package modules

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/tcpacket/nginx-build/command"
	"github.com/tcpacket/nginx-build/util"

	"github.com/go-git/go-git/v5"
)

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

	switch form {
	case "git":
		_, err := git.PlainClone(m.Name, false, &git.CloneOptions{
			URL: m.URL,
		})
		return err
	case "hg":
		args := []string{form, "clone", m.URL}
		if command.VerboseEnabled {
			return command.Run(args)
		}

		f, err := os.Create(logName)
		if err != nil {
			return command.Run(args)
		}
		util.Fclose(f)

		cmd, err := command.Make(args)
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(f)
		defer writer.Flush()

		cmd.Stderr = writer

		return cmd.Run()
	case "local": // not implemented yet
		return nil
	}

	return fmt.Errorf("form=%s is not supported", form)
}
