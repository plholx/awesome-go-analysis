package main_test

import (
	"testing"

	aga "awesome-go-analysis"
)

func TestDownloadREADMEFile(t *testing.T) {
	filePath, err := aga.DownloadREADMEFile()
	if err != nil {
		t.Error(err)
	}
	t.Log(filePath)
}

func TestGenerateMd(t *testing.T) {
	aga.InitDB()
	aga.GenerateMd()
}
