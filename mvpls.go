package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	flag.Parse()
	tail := flag.Args()
	for i := 0; i < len(tail)-1; i++ {
		MoveFile(tail[i], tail[len(tail)-1])
	}
}

func MoveFile(oldFile, newFile string) {
	if oldFile == "" || newFile == "" {
		return
	}

	oldFile, pathErr := filepath.Abs(oldFile)
	if pathErr != nil {
		log.Fatal(pathErr)
	}

	newFile, pathErr = filepath.Abs(newFile)

	info, statErr := os.Stat(newFile)
	fmt.Println(statErr)
	if (statErr == nil || os.IsExist(statErr)) && info.IsDir() {
		info, _ := os.Stat(oldFile)
		newFile = path.Join(newFile, "/"+info.Name())
	}

	fmt.Println(oldFile)
	fmt.Println(newFile)
	err := os.Rename(oldFile, newFile)
	if err != nil {
		log.Fatal(err)
	}
}
