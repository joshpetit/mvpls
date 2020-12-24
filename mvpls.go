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
		moveFile(tail[i], tail[len(tail)-1])
	}
}

func moveFile(oldFile, newFile string) {
	if oldFile == "" || newFile == "" {
		return
	}

	oldFile, pathErr := filepath.Abs(oldFile)
	if pathErr != nil {
		log.Fatal(pathErr)
	}

	newFile, pathErr = filepath.Abs(newFile)

	info, statErr := os.Stat(newFile)
	if statErr != nil {
		log.Fatal(statErr)
	}
	fmt.Println(info)
	if info.IsDir() {
		fmt.Println("NEW FILE IS DIR")
		info, statErr := os.Stat(oldFile)
		newFile = path.Join(newFile, "/"+info.Name())
		if statErr != nil {
			log.Fatal(statErr)
		}

	}

	fmt.Println(newFile)
	fmt.Println(oldFile)
	err := os.Rename(oldFile, newFile)
	if err != nil {
		log.Fatal(err)
	}
}
