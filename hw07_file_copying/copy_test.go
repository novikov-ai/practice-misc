package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestDataPath    = "testdata/"
	TempFilePrefix  = "test_*.txt"
	InvalidFilePath = "/dev/urandom"
	SourceFilePath  = "testdata/input.txt"
)

func TestCopyOffsetValue(t *testing.T) {
	tmpDirPath, tmpFile := initTempStorage(t)
	defer cleanUpTempStorage(t, tmpDirPath, tmpFile)

	fileSize, err := getFileSizeFromPath(SourceFilePath)
	require.Nil(t, err)

	testCases := []struct {
		testName    string
		offsetValue int64
		validCase   bool
	}{
		{testName: "Offset less than source file size.", offsetValue: fileSize - 1, validCase: true},
		{testName: "Offset equal to source file size.", offsetValue: fileSize, validCase: true},
		{testName: "Offset more than source file size.", offsetValue: fileSize + 1, validCase: false},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.testName, func(t *testing.T) {
			err = Copy(SourceFilePath, tmpFile.Name(), testCase.offsetValue, 0)
			if testCase.validCase {
				require.Nil(t, err)
			} else {
				require.True(t, err.Error() == ErrOffsetExceedsFileSize.Error())
			}
		})
	}
}

func TestCopyInvalidArguments(t *testing.T) {
	tmpDirPath, tmpFile := initTempStorage(t)
	defer cleanUpTempStorage(t, tmpDirPath, tmpFile)

	testCases := []struct {
		offset, limit int64
	}{
		{offset: -1, limit: 0},
		{offset: 0, limit: -1},
		{offset: -1, limit: -1},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run("Invalid args", func(t *testing.T) {
			err := Copy(SourceFilePath, tmpFile.Name(), testCase.offset, testCase.limit)
			require.True(t, err.Error() == ErrWrongArguments.Error())
		})
	}
}

func TestCopyInvalidFiles(t *testing.T) {
	tmpDirPath, tmpFile := initTempStorage(t)
	defer cleanUpTempStorage(t, tmpDirPath, tmpFile)

	t.Run("Coping invalid files.", func(t *testing.T) {
		err := Copy(InvalidFilePath, tmpFile.Name(), 0, 0)
		require.True(t, err.Error() == ErrUnsupportedFile.Error())
	})
}

func TestCopy(t *testing.T) {
	testCases := []struct {
		outputPath    string
		offset, limit int64
	}{
		{outputPath: "out_offset0_limit0.txt", offset: 0, limit: 0},
		{outputPath: "out_offset0_limit10.txt", offset: 0, limit: 10},
		{outputPath: "out_offset0_limit1000.txt", offset: 0, limit: 1000},
		{outputPath: "out_offset0_limit10000.txt", offset: 0, limit: 10000},
		{outputPath: "out_offset100_limit1000.txt", offset: 100, limit: 1000},
		{outputPath: "out_offset6000_limit1000.txt", offset: 6000, limit: 1000},
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("Reference %s\n", testCase.outputPath), func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", TempFilePrefix)
			if err != nil {
				t.Fatal(err)
			}
			err = Copy(SourceFilePath, tmpFile.Name(), testCase.offset, testCase.limit)
			require.Nil(t, err)

			dstFile, err := os.Open(tmpFile.Name())
			require.Nil(t, err)

			refFilePath := fmt.Sprintf("%s%s", TestDataPath, testCase.outputPath)
			refFile, err := os.Open(refFilePath)
			require.Nil(t, err)

			isCopiedEqualToRef, err := isEqualFiles(*dstFile, *refFile)
			assert.Nil(t, err)
			assert.True(t, isCopiedEqualToRef)

			refFile.Close()
			dstFile.Close()

			os.Remove(tmpFile.Name())
		})
	}
}

func isEqualFiles(first, second os.File) (bool, error) {
	firstFileBytes, err := getFileBytes(&first)
	if err != nil {
		return false, err
	}

	secondFileBytes, err := getFileBytes(&second)
	if err != nil {
		return false, err
	}

	return bytes.Equal(firstFileBytes, secondFileBytes), nil
}

func getFileBytes(file *os.File) ([]byte, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return readBytes(file, 0, fileInfo.Size()), nil
}

func readBytes(file *os.File, offset, limit int64) []byte {
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

func getFileSizeFromPath(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func initTempStorage(t *testing.T) (string, *os.File) {
	t.Helper()

	tmpDirPath, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}

	tmpFile, err := os.CreateTemp("", TempFilePrefix)
	if err != nil {
		t.Fatal(err)
	}

	return tmpDirPath, tmpFile
}

func cleanUpTempStorage(t *testing.T, dirPath string, tmpFile *os.File) {
	t.Helper()

	err := os.Remove(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		t.Fatal(err)
	}
}
