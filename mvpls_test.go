package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//REFACTOR: Move all test dirs to default temp directory
func TestMoveFile(t *testing.T) {
	testDir := createTestDir()
	defer os.RemoveAll(testDir)

	oldFile := filepath.Join(testDir, "file11.png")
	newFile := filepath.Join(testDir, "file11.new")
	MoveFile(oldFile, newFile)
	if _, err := os.Stat(newFile); err != nil && os.IsNotExist(err) {
		t.Error("File Not Moved")
	}

	oldFile = filepath.Join(testDir, "level2/file22.png")
	newFile = filepath.Join(testDir, "file22.png")
	MoveFile(oldFile, newFile)
	if _, err := os.Stat(newFile); err != nil && os.IsNotExist(err) {
		t.Error("file not moved to directory")
	}
	os.Remove("testingDir/beforeDirMove")

	oldFile = filepath.Join(testDir, "file11.jpg")
	newFile = filepath.Join(testDir, "level2/file11.png")

	MoveFile(oldFile, newFile)
	if _, err := os.Stat(newFile); err != nil && os.IsNotExist(err) {
		t.Error("file not moved to directory")
	}

	oldFile = filepath.Join(testDir, "level2")
	newFile = filepath.Join(testDir, "level2.moved")

	MoveFile(oldFile, newFile)
	if _, err := os.Stat(newFile); err != nil && os.IsNotExist(err) {
		t.Error("Test directory not moved")
	}
}

func createTestDir() (testDir string) {
	pdir := os.TempDir()
	testDir, _ = ioutil.TempDir(pdir, "*-mvplsTest")

	fmt.Println(testDir)
	movedDir := filepath.Join(testDir, "moved")
	os.Mkdir(movedDir, 0777)
	testSubs := filepath.Join(testDir, "level2/level3/level4")
	os.MkdirAll(testSubs, 0777)

	subs := [...]string{"file11.jpg", "file11.png", "level2/file23.jpg", "level2/file21.png", "level2/file22.png", "level2/level3/file31.jpg", "level2/level3/file32.png"}
	var sub string

	for i := 0; i < len(subs); i++ {
		sub = filepath.Join(testDir, subs[i])
		os.Create(sub)
	}

	return
}
