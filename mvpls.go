package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

var regexFlag = flag.String("r", "regex", "The regex to be used to match files")

var s stack

func main() {
	s = make(stack, 0)
	flag.Parse()
	reg, comp_err := regexp.Compile(*regexFlag)

	if comp_err != nil {
		log.Fatal(comp_err)
		return
	}

	tail := flag.Args()
	target := tail[len(tail)-1]

	if *regexFlag == "" {
		for i := 0; i < len(tail)-1; i++ {
			MoveFile(tail[i], target)
		}
		return
	}

	_, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(tail)-1; i++ {
		ProbeDirectory(tail[i], target, reg)
	}
}

func ProbeDirectory(dir, target string, reg *regexp.Regexp) {
	if dir == "" || target == "" {
		return
	}
	dir, _ = filepath.Abs(dir)
	target, _ = filepath.Abs(target)
	s = s.Push(dir)
	for len(s) != 0 {
		s, dir = s.Pop()
		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				fmt.Println(info.Name())

				if path != dir {
					if info.IsDir() {
						fmt.Println(info.Name(), "added to stack")
						s = s.Push(path)
					} else if reg.MatchString(info.Name()) {
						fmt.Println(info.Name(), "attempting move")
						MoveFile(path, target)
					}
				}

				return nil
			})
		if err != nil {
			log.Fatal(err)
		}
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
