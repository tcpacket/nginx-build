package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tcpacket/waf-builder/builder"
	"github.com/tcpacket/waf-builder/command"
	"github.com/tcpacket/waf-builder/util"
)

const DefaultDownloadTimeout = time.Duration(900) * time.Second

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func download(b *builder.Builder) error {
	log.Trace().Msgf("Downloading %s to %s", b.DownloadURL(), b.ArchivePath())
	c := &http.Client{
		Timeout: DefaultDownloadTimeout,
	}
	res, err := c.Get(b.DownloadURL())
	if err != nil {
		return err
	}
	defer util.Fclose(res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s. status code: %d", b.DownloadURL(), res.StatusCode)
	}

	tmpFileName := b.ArchivePath() + ".download"
	f, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}
	defer util.Fclose(f)
	if _, err := io.Copy(f, res.Body); err != nil && err != io.EOF {
		return err
	}
	if err := os.Rename(tmpFileName, b.ArchivePath()); err != nil {
		return err
	}
	return nil
}

func downloadAndExtract(b *builder.Builder) error {
	if !util.FileExists(b.SourcePath()) {
		if !util.FileExists(b.ArchivePath()) {
			if err := download(b); err != nil {
				return fmt.Errorf("failed to download %s. %s", b.SourcePath(), err.Error())
			}
		}
		log.Printf("Extract %s.....", b.ArchivePath())

		if err := extractArchive(b.ArchivePath()); err != nil {
			return fmt.Errorf("failed to extract %s. %s", b.ArchivePath(), err.Error())
		}
	} else {
		log.Printf("%s already exists.", b.SourcePath())
	}
	return nil
}

func downloadAndExtractParallel(b *builder.Builder) {
	if err := downloadAndExtract(b); err != nil {
		util.PrintFatalMsg(err, b.LogPath())
	}
}
