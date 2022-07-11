package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestDataPath     = "testdata/"
	TestDataTempPath = "testdata/temp/"
	SourceFilePath   = "testdata/input.txt"
)

func TestCopyOffsetValue(t *testing.T) {
	err := os.Mkdir(TestDataTempPath, 0o740)
	if err != nil && !os.IsExist(err) {
		require.Nil(t, err)
	}
	dstFilePath := fmt.Sprintf("%s%s", TestDataTempPath, "test_offset_value.txt")

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
			err = Copy(SourceFilePath, dstFilePath, testCase.offsetValue, 0)
			if testCase.validCase {
				require.Nil(t, err)
			} else {
				require.True(t, err.Error() == ErrOffsetExceedsFileSize.Error())
			}
		})
	}

	err = os.Remove(dstFilePath)
	assert.Nil(t, err)

	err = os.Remove(TestDataTempPath)
	assert.Nil(t, err)
}

func TestCopyInvalidArguments(t *testing.T) {
	dstFilePath := fmt.Sprintf("%s%s", TestDataTempPath, "test_invalid_args.txt")

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
			err := Copy(SourceFilePath, dstFilePath, testCase.offset, testCase.limit)
			require.True(t, err.Error() == ErrWrongArguments.Error())
		})
	}
}

func TestCopyInvalidFiles(t *testing.T) {
	src := "/dev/urandom"
	dstFilePath := fmt.Sprintf("%s%s", TestDataTempPath, "test_invalid_args.txt")

	t.Run("ads", func(t *testing.T) {
		err := Copy(src, dstFilePath, 0, 0)
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

	err := os.Mkdir(TestDataTempPath, 0o740)
	if err != nil && !os.IsExist(err) {
		require.Nil(t, err)
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(fmt.Sprintf("Reference %s\n", testCase.outputPath), func(t *testing.T) {
			dstFilePath := fmt.Sprintf("%s%s", TestDataTempPath, testCase.outputPath)
			err := Copy(SourceFilePath, dstFilePath, testCase.offset, testCase.limit)
			require.Nil(t, err)

			dstFile, err := os.Open(dstFilePath)
			require.Nil(t, err)

			refFilePath := fmt.Sprintf("%s%s", TestDataPath, testCase.outputPath)
			refFile, err := os.Open(refFilePath)
			require.Nil(t, err)

			isCopiedEqualToRef, err := isEqualFiles(*dstFile, *refFile)
			assert.Nil(t, err)
			assert.True(t, isCopiedEqualToRef)

			refFile.Close()
			dstFile.Close()

			err = os.Remove(dstFilePath)
			assert.Nil(t, err)
		})
	}

	err = os.Remove(TestDataTempPath)
	assert.Nil(t, err)
}

func isEqualFiles(first, second os.File) (bool, error) {
	firstFileInfo, err := first.Stat()
	if err != nil {
		return false, err
	}
	firstFileBytes := ReadBytes(&first, 0, firstFileInfo.Size())

	secondFileInfo, err := second.Stat()
	if err != nil {
		return false, err
	}
	secondFileBytes := ReadBytes(&second, 0, secondFileInfo.Size())

	return bytes.Equal(firstFileBytes, secondFileBytes), nil
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
