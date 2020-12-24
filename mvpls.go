package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	tail := flag.Args()
	for i := 0; i < len(tail)-1; i++ {
		moveFile(tail[i], tail[len(tail)-1])
	}
}

func moveFile(oldFile, newFile string) {
	if oldFile == "" || newFile == "" {
		return
	}

	info, statErr := os.Stat(oldFile)
	if statErr != nil {
		log.Fatal(statErr)
	}

	fmt.Println("ya yeet")
	oldFile, pathErr := filepath.Abs(oldFile)
	if pathErr != nil {
		log.Fatal(pathErr)
	}
	newFile, pathErr = filepath.Abs(newFile)

	err := os.Rename(oldFile, newFile)
	if err != nil {
		log.Fatal(err)
	}
}
