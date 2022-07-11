package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrWrongDataCopied       = errors.New("wrong data copied")
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

	if limit == 0 {
		limit = fileSize
	}

	inputBytes := ReadBytes(srcFile, offset, limit)
	reader := bytes.NewReader(inputBytes)

	destFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	writer := bufio.NewWriter(destFile)

	bytesCopied, err := io.CopyN(writer, reader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if bytesCopied > limit {
		return ErrWrongDataCopied
	}

	return nil
}

func ReadBytes(file *os.File, offset, limit int64) []byte {
	buffer := make([]byte, limit)
	bytesRead, err := file.ReadAt(buffer, offset)
	if err == io.EOF {
		if len(buffer) > 0 {
			return buffer[0:bytesRead]
		}

		return nil
	}
	return buffer[0:bytesRead]
}
