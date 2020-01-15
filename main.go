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
	"regexp"
	"runtime"

	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	NINDENT = "    "
	SINDENT = "│   "
)

var ndirs int
var nfiles int

var maxDepth int
var printAll bool
var dirsOnly bool
var pattern string

var au aurora.Aurora

func isatty() bool {
	fd := os.Stdout.Fd()
	return terminal.IsTerminal(int(fd))
}

func printErr(err error) {
	fmt.Println(au.Red(err).Bold())
}

func printFile(prefix string, info os.FileInfo, isLast bool) {
	var fname = info.Name()
	var mode = info.Mode()
	if isLast {
		prefix += "└── "
	} else {
		prefix += "├── "
	}

	if info.IsDir() {
		fmt.Printf("%s%s\n", prefix, au.Blue(fname).Bold())
	} else if !dirsOnly {
		if mode&0111 != 0 { // checks if it's executable
			fmt.Printf("%s%s\n", prefix, au.Green(fname).Bold())
		} else {
			fmt.Printf("%s%s\n", prefix, fname)
		}
	}
}

func filterDirs(files []os.FileInfo) (dirs []os.FileInfo, length int) {
	for _, f := range files {
		if f.IsDir() {
			fname := f.Name()
			if fname[0] != '.' || printAll {
				dirs = append(dirs, f)
				length++
			}
		}
	}
	return
}

func filterExpr(files []os.FileInfo) (dirs []os.FileInfo, length int) {
	for _, f := range files {
		ok, err := regexp.MatchString(pattern, f.Name())
		if ok {
			dirs = append(dirs, f)
			length++
		} else if err != nil {
			printErr(err)
		}
	}
	return
}

func filterHidden(files []os.FileInfo) (farr []os.FileInfo, length int) {
	for _, f := range files {
		fname := f.Name()
		if fname[0] != '.' {
			farr = append(farr, f)
			length++
		}
	}
	return
}

func readDir(filename string) ([]os.FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []os.FileInfo{}, err
	}
	defer file.Close()
	return file.Readdir(-1)
}

// TODO: Design the pattern matching properly.
func walkDir(root string, prefix string, depth int) {
	if depth != maxDepth {
		var arrlen int

		files, err := readDir(root)
		if err != nil {
			printErr(err)
			return
		}
		arrlen = len(files)
		if !printAll {
			files, arrlen = filterHidden(files)
		}
		if dirsOnly {
			files, arrlen = filterDirs(files)
		}

		for i, finfo := range files {
			var isLast = i == arrlen-1

			printFile(prefix, finfo, isLast)
			if finfo.IsDir() {
				newPath := fmt.Sprintf("%s/%s", root, finfo.Name())
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
}

func main() {
	var rootdir = "./"
	if narg := flag.NArg(); narg > 0 {
		rootdir = flag.Arg(narg-1)
	}

	fmt.Println(rootdir)
	walkDir(rootdir, "", 0)
	fmt.Printf("\n%d directories, %d files\n", ndirs, nfiles)
}

func init() {
	var colours = runtime.GOOS != "windows"

	flag.BoolVar(&printAll, "a", false, "Prints all files including the hidden ones.")
	flag.BoolVar(&dirsOnly, "d", false, "Prints only the directories.")
	flag.StringVar(&pattern, "e", "", "Prints only the files that match the regex. (Coming soon...)")
	flag.IntVar(&maxDepth, "l", -1, "Max display depth of the directory tree.")
	flag.BoolVar(&colours, "c", true, "Set to false to disable colours.")

	flag.Parse()
	au = aurora.NewAurora(colours && isatty())
}
