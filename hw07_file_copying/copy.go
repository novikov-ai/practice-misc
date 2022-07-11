package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

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

	progressBar := progressBarInit(offset, limit, fileSize)
	inputBytes := ReadBytesWithProgress(srcFile, offset, limit, progressBar)

	reader := bytes.NewReader(inputBytes)
	writer := bufio.NewWriter(destFile)

	_, err = io.CopyN(writer, reader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	progressBar.Finish()

	return nil
}

func ReadBytesWithProgress(file *os.File, offset, limit int64, progressBar *pb.ProgressBar) []byte {
	buffer := make([]byte, limit)
	reader := bufio.NewReader(file)

	var byteIndex, bytesRead int64

	for {
		symbol, err := reader.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil
		}

		if byteIndex < offset {
			byteIndex++
			continue
		}

		buffer[bytesRead] = symbol
		bytesRead++

		progressBar.Increment()

		if bytesRead == limit {
			break
		}
	}

	return buffer[0:bytesRead]
}

func progressBarInit(offset, limit, fileSize int64) *pb.ProgressBar {
	progressBarMaxValue := limit
	if limit > fileSize {
		progressBarMaxValue = fileSize
	} else if offset+limit > fileSize {
		progressBarMaxValue = fileSize - offset
	}
	return pb.Start64(progressBarMaxValue)
}
