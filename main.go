/*
 * Twig
 * Copyright (C) 2019  Nicolò Santamaria
 *
 * Twig is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Twig is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"os"
	"fmt"
	"flag"

	"github.com/logrusorgru/aurora"
	// "golang.org/x/crypto/ssh/terminal"
)

const (
	NINDENT = "    "
	SINDENT = "│   "
)

var ndirs int
var nfiles int

var printAll bool

// TODO: fix colours
func printFile(prefix string, info os.FileInfo) {
	var fname = info.Name()
	var mode = info.Mode()

	if info.IsDir() {
		fmt.Printf("%s%s\n", prefix, aurora.BrightBlue(fname))
	} else if mode&(1<<0) != 0 { // checks if it's executable
		fmt.Printf("%s%s\n", prefix, aurora.BrightGreen(fname))
	} else {
		fmt.Printf("%s%s\n", prefix, fname)
	}
}

func walkDir(root string, prefix string, depth int) {
	var nfiles uint

	f, err := os.Open(root)
	if err != nil {
		return
	}
	defer f.Close()

	fileInfo, err := f.Readdir(-1)
	if err != nil {
		return
	}
	nfiles = len(fileInfo)

	for i, finfo := range fileInfo {
		// var line string
		var isLast = i == nfiles-1
		var fname = finfo.Name()

		if !printAll && fname[0] == '.' {
			continue
		}

		if isLast {
			fmt.Printf("%s└── %s\n", prefix, fname)
			// line = prefix + "└── "
		} else {
			fmt.Printf("%s├── %s\n", prefix, fname)
			// line = prefix + "├── "
		}
		// printFile(line, finfo)

		if finfo.IsDir() {
			newPath := fmt.Sprintf("%s/%s", root, fname)
			if isLast {
				walkDir(newPath, prefix+NINDENT, depth+1)
			} else {
				walkDir(newPath, prefix+SINDENT, depth+1)
			}
			ndirs++
		} else {
			nfiles++
		}
	}
}

func printHelp() {
	// Put the help message here.
}

func main() {
	var rootdir = "./"
	var nargs = flag.NArg()
	if nargs > 0 {
		rootdir = os.Args[len(os.Args)-1]
	}

	fmt.Println(rootdir)
	walkDir(rootdir, "", 0)
	fmt.Printf("\n%d directories, %d files\n", ndirs, nfiles)
}

func init() {
	flag.BoolVar(&printAll, "a", false, "Print all files including the hiddnen ones.")
	flag.Parse()
}
