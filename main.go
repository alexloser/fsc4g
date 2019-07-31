// A simple parallel file scancer, faster than other tools
package main

import (
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	if len(os.Args) != 3 && len(os.Args) != 2 {
		print("Usage:\n\t"+os.Args[0], "PATH REGEX")
		print("\n\t"+os.Args[0], os.Args[0], "PATH")
		print("Example:\n\t"+os.Args[0], "/home/user/tmp \"^.*\\.txt$\"")
		os.Exit(-1)
	}

	var dir, sre string
	dir = os.Args[1]
	stat, err := os.Stat(dir)

	if err != nil || !stat.IsDir() {
		print("Error:", dir, "not existed!!!")
		os.Exit(-1)
	}

	if len(os.Args) == 3 {
		sre = os.Args[2]
	}

	FileScan(dir, sre, func(info string) {
		print(info)
	})

}
