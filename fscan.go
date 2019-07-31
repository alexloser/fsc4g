// A simple parallel file scancer, faster than other tools
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

func print(obj ...interface{}) {
	for _, each := range obj {
		fmt.Printf("%v ", each)
	}
	fmt.Print("\n")
}

type PathWalkerCallback func(string)

type PathWalker struct {
	fileinfo chan string
	wg       sync.WaitGroup
	callback PathWalkerCallback
}

func Make(callback PathWalkerCallback) *PathWalker {
	return &PathWalker{
		fileinfo: make(chan string),
		callback: callback,
	}
}

func (self PathWalker) dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(dir)
		return nil
	}

	return entries
}

func (self PathWalker) format(parent string, f os.FileInfo) string {
	size := float64(f.Size()) / 1024.0
	unit := "KB"

	if size > 1024*1000 {
		size /= 1024.0
		unit = "MB"
	}

	if f.IsDir() {
		return fmt.Sprintf("D %10.2f%s  %s", size, unit, filepath.Join(parent, f.Name()))
	} else {
		return fmt.Sprintf("F %10.2f%s  %s", size, unit, filepath.Join(parent, f.Name()))
	}
}

func (self *PathWalker) walk(dir string) {
	defer self.wg.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: reading \"%v\" failed!\n", dir)
			return
		}
	}()
	for _, entry := range self.dirents(dir) {
		self.fileinfo <- self.format(dir, entry)

		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			self.wg.Add(1)
			go self.walk(subdir)
		}
	}
}

func (self *PathWalker) Go(dir string) {
	self.wg.Add(1)

	go self.walk(dir)

	go func() {
		self.wg.Wait()
		close(self.fileinfo)
	}()

loop:
	for {
		select {
		case info, ok := <-self.fileinfo:
			if !ok {
				break loop
			}
			self.callback(info)
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func FileScan(dir, sre string, callabck func(info string)) {
	if len(sre) > 0 {
		regex := regexp.MustCompile(sre)
		Make(func(info string) {
			if regex.MatchString(info) {
				callabck(info)
			}
		}).Go(dir)
	} else {
		Make(func(info string) {
			callabck(info)
		}).Go(dir)
	}
}
