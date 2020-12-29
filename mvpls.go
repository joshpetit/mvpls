package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type stack []string

func (s stack) Push(v string) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, string) {
	l := len(s)
	if l == 0 {
		return s, ""
	}
	return s[:l-1], s[l-1]
}

func main() {
	s := make(stack, 0)
	flag.Parse()
	tail := flag.Args()
	for i := 0; i < len(tail)-1; i++ {
		MoveFile(tail[i], tail[len(tail)-1])
	}
	_, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	s = s.Push("testdir")
	var dir string
	for len(s) != 0 {
		s, dir = s.Pop()
		err = filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() != dir {
					fmt.Println(info.Name())
					s = s.Push(path)
				} else if info.Name() != dir {
					fmt.Println(info.Name())
				}
				return nil
			})

		s, dir = s.Pop()
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
