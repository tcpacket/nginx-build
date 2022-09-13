package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tcpacket/nginx-build/builder"
	"github.com/tcpacket/nginx-build/command"
	"github.com/tcpacket/nginx-build/util"
)

const DefaultDownloadTimeout = time.Duration(900) * time.Second

func extractArchive(path string) error {
	return command.Run([]string{"tar", "zxvf", path})
}

func download(b *builder.Builder) error {
	c := &http.Client{
		Timeout: DefaultDownloadTimeout,
	}
	res, err := c.Get(b.DownloadURL())
	if err != nil {
		return err
	}
	defer util.Fclose(res.Body)

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

			log.Printf("Download %s.....", b.SourcePath())

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
