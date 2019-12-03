
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

	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	NINDENT = "    "
	SINDENT = "│   "
)

var ndirs int
var nfiles int

var printAll bool
var dirsOnly bool
var pattern string

var au aurora.Aurora

func isatty() bool {
	fd := os.Stdout.Fd()
	return terminal.IsTerminal(int(fd))
}

func printFile(prefix string, info os.FileInfo) {
	var fname = info.Name()
	var mode = info.Mode()

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
			fmt.Println(au.Red(err).Bold())
		}
	}
	return
}

// TODO: Design the pattern matching properly.
func walkDir(root string, prefix string, depth int) {
	var arrlen int
	f, err := os.Open(root)
	if err != nil {
		return
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return
	}

	if dirsOnly {
		files, arrlen = filterDirs(files)
	} else {
		arrlen = len(files)
	}

	// if pattern != "" {
	// 	files, arrlen = filterExpr(files)
	// } else if dirsOnly {
	// 	files, arrlen = filterDirs(files)
	// } else {
	// 	arrlen = len(files)
	// }

	for i, finfo := range files {
		var line = prefix
		var isLast = i == arrlen-1
		var fname = finfo.Name()

		if fname[0] != '.' || printAll {
			if isLast {
				line += "└── "
			} else {
				line += "├── "
			}
			printFile(line, finfo)

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
		} else {
			i--
		}
	}
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
	flag.BoolVar(&printAll, "a", false, "Prints all files including the hiddnen ones.")
	flag.BoolVar(&dirsOnly, "d", false, "Prints only the directories.")
	flag.StringVar(&pattern, "e", "", "Prints only the files that match the regex. (Coming soon...)")
	colours := flag.Bool("c", true, "Set to false to disable colours.")
	flag.Parse()
	au = aurora.NewAurora(*colours && isatty())
}
