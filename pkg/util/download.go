package util

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type Downloader struct {
	concurrency int
	resume      bool

	bar *progressbar.ProgressBar
}

func NewDownloader(concurrency int, resume bool) *Downloader {
	return &Downloader{concurrency: concurrency, resume: resume}
}

func (d *Downloader) Download(strURL, filename string) error {
	if filename == "" {
		// 解析 URL 获取文件名
		parsedURL, err := url.Parse(strURL)
		if err != nil {
			return fmt.Errorf("failed to parse URL: %v", err)
		}

		filename = path.Base(parsedURL.Path)
	}

	resp, err := http.Head(strURL)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		return d.multiDownload(strURL, filename, int(resp.ContentLength))
	}

	return d.singleDownload(strURL, filename)
}

func (d *Downloader) multiDownload(strURL, filename string, contentLen int) error {

	d.setBar(contentLen)

	partSize := contentLen / d.concurrency

	// 创建部分文件的存放目录
	partDir := d.getPartDir(filename)
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	eg := errgroup.Group{}
	eg.SetLimit(d.concurrency)

	rangeStart := 0

	for i := 0; i < d.concurrency; i++ {
		i := i
		rangeStart := rangeStart
		eg.Go(func() error {

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if i == d.concurrency-1 {
				rangeEnd = contentLen
			}

			downloaded := 0
			if d.resume {
				partFileName := d.getPartFilename(filename, i)
				content, err := os.ReadFile(partFileName)
				if err != nil {
					return errors.WithStack(err)
				}

				downloaded = len(content)

				err = d.bar.Add(downloaded)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			return d.downloadPartial(strURL, filename, rangeStart+downloaded, rangeEnd, i)
		})

		rangeStart += partSize + 1
	}

	err := eg.Wait()

	if err != nil {
		return errors.WithStack(err)
	}

	err = d.merge(filename)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (d *Downloader) downloadPartial(strURL, filename string, rangeStart, rangeEnd, i int) (err error) {
	if rangeStart >= rangeEnd {
		return
	}

	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	if d.resume {
		flags |= os.O_APPEND
	}

	partFile, err := os.OpenFile(d.getPartFilename(filename, i), flags, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(partFile, d.bar), resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return errors.WithStack(err)
	}
	return nil
}

func (d *Downloader) merge(filename string) error {
	destFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer destFile.Close()

	for i := 0; i < d.concurrency; i++ {
		partFileName := d.getPartFilename(filename, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = io.Copy(destFile, partFile)
		if err != nil {
			return errors.WithStack(err)
		}
		err = partFile.Close()
		if err != nil {
			return errors.WithStack(err)
		}
		err = os.Remove(partFileName)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// getPartDir 部分文件存放的目录
func (d *Downloader) getPartDir(filename string) string {
	return strings.SplitN(filename, ".", 2)[0] + "_tmp"
}

// getPartFilename 构造部分文件的名字
func (d *Downloader) getPartFilename(filename string, partNum int) string {
	partDir := d.getPartDir(filename)
	return fmt.Sprintf("%s/%s-%d", partDir, filename, partNum)
}

func (d *Downloader) singleDownload(strURL, filename string) error {
	resp, err := http.Get(strURL)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	d.setBar(int(resp.ContentLength))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), resp.Body, buf)
	return errors.WithStack(err)
}

func (d *Downloader) setBar(length int) {
	d.bar = progressbar.NewOptions(
		length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("downloading..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
