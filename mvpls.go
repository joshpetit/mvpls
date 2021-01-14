package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type Operation int

const (
	Move Operation = iota
	Copy
	Remove
)

type FileOperation func(string, string)

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
var copyFlag = flag.Bool("c", false, "Copy the requested files, please")
var removeFlag = flag.Bool("remove", false, "Copy the requested files, please")

var s stack
var version = "0.0.0"

func main() {
	flag.Parse()
	tail := flag.Args()
	tailLen := len(tail) - 1
	if tailLen == -1 {
		fmt.Printf("Mvpls!\nv%s\nRun mvpls --help to ask for help (please)\n", version)
		return
	}
	target := tail[tailLen]
	var fileOperation FileOperation
	var op Operation

	switch {
	case *copyFlag:
		op = Copy
		fileOperation = CopyFile
	case *removeFlag:
		op = Remove
		fileOperation = MoveFile
	default:
		op = Move
		fileOperation = MoveFile
	}
	if *regexFlag == "" {
		for i := 0; i < tailLen; i++ {
			switch op {
			case Copy:
				CopyFile(tail[i], target)
			default:
				MoveFile(tail[i], target)
			}
		}
		return
	}

	reg, comp_err := regexp.Compile(*regexFlag)

	if comp_err != nil {
		log.Fatal(comp_err)
		return
	}

	s = make(stack, 0)
	for i := 0; i < tailLen; i++ {
		ProbeDirectory(tail[i], target, reg, fileOperation)
	}
}

func probeable(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Error: %s does not exist", dir))
		return false
	}
	if !stat.IsDir() {
		log.Fatal(fmt.Sprintf("Error: %s is not a directory", dir))
		return false
	}
	return true
}

func ProbeDirectory(dir, target string, reg *regexp.Regexp, operate FileOperation) {
	if dir == "" || target == "" {
		return
	}
	dir, _ = filepath.Abs(dir)
	target, _ = filepath.Abs(target)
	if !probeable(dir) || !probeable(target) {
		return
	}
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
						operate(path, target)
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

func CopyFile(oldFilePath, newFilePath string) {
	oldFilePath, oldPathErr := filepath.Abs(oldFilePath)
	if oldPathErr != nil {
		log.Fatal(oldPathErr)
	}

	newFilePath, newPathErr := filepath.Abs(newFilePath)
	if newPathErr != nil {
		log.Fatal(newPathErr)
	}

	oldFile, newInfoErr := os.Open(oldFilePath)
	if newInfoErr != nil {
		log.Fatal(newInfoErr)
	}

	defer oldFile.Close()
	oldFileInfo, oldInfoErr := os.Stat(oldFilePath)
	if oldFileInfo.IsDir() {
		log.Fatal("Copying entire directories is not currently supported :(")
		return
	}

	if oldInfoErr != nil {
		log.Fatal(oldInfoErr)
		return
	}

	newFileInfo, newInfoErr := os.Stat(newFilePath)

	if newInfoErr == nil && newFileInfo.IsDir() {
		newInfo, _ := os.Stat(oldFilePath)
		newFilePath = path.Join(newFilePath, newInfo.Name())
	}

	newFile, newFileErr := os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE, oldFileInfo.Mode().Perm())
	if newFileErr != nil {
		log.Fatal(newInfoErr)
	}

	defer newFile.Close()
	_, newInfoErr = io.Copy(newFile, oldFile)
	fmt.Println(oldFilePath, "-- copied -->", newFilePath)
	if newInfoErr != nil {
		log.Fatal(newInfoErr)
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
