package main

import (
	"flag"
	"fmt"
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

var regexFlag = flag.String("r", "", "The regex to be used to match files")

var s stack

func main() {
	s = make(stack, 0)
	flag.Parse()

	tail := flag.Args()
	tailLen := len(tail) - 1
	target := tail[tailLen]

	if *regexFlag == "" {
		for i := 0; i < tailLen; i++ {
			MoveFile(tail[i], target)
		}
		return
	}

	reg, comp_err := regexp.Compile(*regexFlag)

	if comp_err != nil {
		log.Fatal(comp_err)
		return
	}
	for i := 0; i < tailLen; i++ {
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
				if path != dir {
					if info.IsDir() {
						s = s.Push(path)
					} else if reg.MatchString(info.Name()) {
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
	mvDir := os.IsPathSeparator(newFile[len(newFile)-1])

	oldFile, pathErr := filepath.Abs(oldFile)
	if pathErr != nil {
		log.Fatal(pathErr)
		return
	}

	oldInfo, statErr := os.Stat(oldFile)

	if statErr != nil && os.IsNotExist(statErr) {
		log.Fatal(statErr)
		return
	}

	newFile, pathErr = filepath.Abs(newFile)
	newInfo, statErr := os.Stat(newFile)

	if statErr == nil && newInfo.IsDir() {

		fmt.Println("ok")
		newInfo, _ := os.Stat(oldFile)
		newFile = path.Join(newFile, newInfo.Name())
	} else if statErr != nil && mvDir && !oldInfo.IsDir() && os.IsNotExist(statErr) {
		log.Fatal(statErr)
		return
	}

	fmt.Println(oldFile, "->", newFile)

	err := os.Rename(oldFile, newFile)
	if err != nil {
		log.Fatal(err)
	}
}
