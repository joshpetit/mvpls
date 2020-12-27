package main

import (
	"os"
	"testing"
)

func TestMoveFile(t *testing.T) {
	createTestFile("testFile", t)

	MoveFile("testFile", "testFileMoved")
	if _, err := os.Stat("testFileMoved"); err != nil && os.IsNotExist(err) {
		t.Error("File Not Moved")
	}
	os.Remove("testFileMoved")

	createTestFile("testDir/", t)

	MoveFile("testDir", "testingDir")
	if _, err := os.Stat("testingDir"); err != nil && os.IsNotExist(err) {
		t.Error("Test directory not moved")
	}

	createTestFile("beforeDirMove", t)
	MoveFile("beforeDirMove", "testingDir/")

	if _, err := os.Stat("testingDir/beforeDirMove"); err != nil && os.IsNotExist(err) {
		t.Error("file not moved to directory")
	}
	createTestFile("fullMove", t)
	MoveFile("fullMove", "testingDir/")

	if _, err := os.Stat("testingDir/fullMove"); err != nil && os.IsNotExist(err) {
		t.Error("file not moved to directory")
	}

	os.Remove("testingDir/afterMove")
	os.Remove("testingDir")
}

func createTestFile(file string, t *testing.T) {
	var err error
	if file[len(file)-1] == '/' {
		err = os.Mkdir(file[0:len(file)-1], 0777)
	} else {
		_, err = os.Create(file)
	}

	if err != nil {
		t.Error("Error creating {}", file)
	}

}
