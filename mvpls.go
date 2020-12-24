package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var oldFile = flag.String("i", "", "old file")
var newFile = flag.String("o", "", "new file")

func moveFile(oldFile, newFile string) {
	if oldFile == "" || newFile == "" {
		return
	}
	oldFile, pathErr := filepath.Abs(oldFile)
	if pathErr != nil {
		fmt.Println("Error finding file %s", pathErr)
	}
	newFile, pathErr = filepath.Abs(newFile)

	err := os.Rename(oldFile, newFile)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	flag.Parse()

	fmt.Println(*oldFile)
	fmt.Println(*newFile)

	moveFile(*oldFile, *newFile)
}
