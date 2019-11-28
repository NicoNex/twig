package main

import (
	"os"
	"fmt"
)

const (
	NINDENT = "    "
	SINDENT = "│   "
)

// Checks if the file is hidden.
func isHidden() bool {

}

func printTree(root, prefix string, depth int) {
	f, err := os.Open(root)
	if err != nil {
		return
	}

	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return
	}

	for i, file := range fileInfo {
		var isLast = i == len(fileInfo)-1
		var fname = file.Name()

		if isLast {
			fmt.Printf("%s└── %s\n", prefix, fname)
		} else {
			fmt.Printf("%s├── %s\n", prefix, fname)
		}

		if file.IsDir() {
			path := fmt.Sprintf("%s/%s", root, fname)
			if isLast {
				printTree(path, prefix+NINDENT, depth+1)
			} else {
				printTree(path, prefix+SINDENT, depth+1)
			}
		}
	}
}

func main() {
	var arglen = len(os.Args)

	if arglen > 1 {
		fmt.Println(os.Args[arglen-1])
		printTree(os.Args[arglen-1], "", 0)
	}
}
