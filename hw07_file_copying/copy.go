package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWrongArguments        = errors.New("wrong arguments: offset or limit")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return ErrWrongArguments
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if limit == 0 {
		limit = fileSize
	}

	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	bar := progressBarInit(offset, limit, fileSize)

	reader := bufio.NewReader(srcFile)
	proxyReader := bar.NewProxyReader(reader)

	writer := bufio.NewWriter(destFile)

	err = WriteFromReader(proxyReader, writer, int(limit))
	if err != nil {
		return err
	}

	progressBarFinish(bar)

	return nil
}

func WriteFromReader(proxyReader *pb.Reader, writer *bufio.Writer, limit int) error {
	buffer := make([]byte, 100)
	for {
		n, err := proxyReader.Read(buffer)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if n >= limit {
			n = limit
		}
		_, err = writer.Write(buffer[:n])
		if err != nil {
			return err
		}

		limit -= n
		if limit <= 0 {
			break
		}
	}

	return writer.Flush()
}

func progressBarInit(offset, limit, fileSize int64) *pb.ProgressBar {
	pb := pb.New64(getSizeOfFileCopy(offset, limit, fileSize)).SetRefreshRate(time.Millisecond * 10)
	return pb.Start()
}

func progressBarFinish(bar *pb.ProgressBar) {
	bar.SetCurrent(bar.Total())
	bar.Finish()
}

func getSizeOfFileCopy(offset, limit, fileSize int64) int64 {
	size := limit
	if limit > fileSize {
		size = fileSize
	} else if offset+limit > fileSize {
		size = fileSize - offset
	}
	return size
}
