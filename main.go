// Description:
//   `fn` is a tool for generating, and parsing, file names based on
//   current date, time, process id and gitsha.
//   use `fn` to generate a file name.
//   use `fn -r [dir]` to get the most recent file (in dir).
//   more options listed below.
// Usage:
//   fn [-m] [-t]
//   fn -g
//   fn -p
//   fn -r [-a|-A] [-f] [-i] [-n <num>] [<dir>]
//   fn -l [-a|-A] [-f] [-i] [-n <num>] [<dir>]
//   fn -L [-a|-A] [-f] [-i] [-n <num>] [<dir>]
//   fn -s [<dir>]
// Options:
//   -m          include milliseconds.
//   -g          return current git sha.
//   -p          return a prochash.
//   -t          return timestamp only.
//   -r          return all files with the most recent prochash.
//                 sorted by time, asc.
//   -l          return all files with current git sha. sorted by time, asc.
//   -L          return all files with most recent git sha. sorted by time, asc.
//   -i          reverse output.
//   -n NUM      limit to n rows.
//   -s          return most recent prochash.
//   -a          show file name only.
//   -A          show absolute paths.
//   -f          remove filename extension. resulting duplicates will be removed.
//   -h --help   show this screen.
//   --version   show version.
//
// Example: 20200927-201950--289b6e2a
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"strings"
	"time"
)

var (
	milliseconds        bool
	gitSHA              bool
	prochash            bool
	timestamp           bool
	filesRecentProcHash bool
	currentGit          bool
	recentGit           bool
	reverse             bool
	limit               int
	recentProcHash      bool
	namesonly           bool
	absolute            bool
	removeExtension     bool
	currentTime         string
)

func init() {
	flag.BoolVar(&milliseconds, "m", false, "include milliseconds.")
	flag.BoolVar(&gitSHA, "g", false, "return current git sha.")
	flag.BoolVar(&prochash, "p", false, "return a prochash.")
	flag.BoolVar(&timestamp, "t", false, "return timestamp only.")
	flag.BoolVar(&filesRecentProcHash, "r", false, "return all files with the most recent prochash.")
	flag.BoolVar(&currentGit, "l", false, "return all files with current git sha. sorted by time, asc.")
	flag.BoolVar(&recentGit, "L", false, "return all files with most recent git sha. sorted by time, asc.")
	flag.BoolVar(&reverse, "i", false, "reverse output.")
	flag.IntVar(&limit, "n", 0, "limit to n rows.")
	flag.BoolVar(&recentProcHash, "s", false, "return all files with the most recent prochash.")
	flag.BoolVar(&namesonly, "a", false, "show file name only.")
	flag.BoolVar(&absolute, "A", false, "show absolute paths.")
	flag.BoolVar(&removeExtension, "f", false, "remove filename extension. resulting duplicates will be removed.")
	usage := strings.NewReplacer("'", "`").Replace(`fn

Description:
	'fn' is a tool for generating, and parsing, file names based on
	current date, time, process id and gitsha.

	use 'fn' to generate a file name.
	use 'fn -r [dir]' to get the most recent file (in dir).
	more options listed below.


Usage:
	fn [-m] [-t]
	fn -g
	fn -p
	fn -r [-a|-A] [-f] [-i] [-n <num>] [<dir>]
	fn -l [-a|-A] [-f] [-i] [-n <num>] [<dir>]
	fn -L [-a|-A] [-f] [-i] [-n <num>] [<dir>]
	fn -s [<dir>]


Options:`)
	flag.Usage = func() {
		fmt.Println(usage)

		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	time := getTime(milliseconds)
	filename := fmt.Sprintf("%s-%s-%s", time, "", getHash(time))
	fmt.Println(filename)
	// fmt.Println("milliseconds:", milliseconds)
	// fmt.Println("gitSHA:", gitSHA)
	// fmt.Println("prochash:", prochash)
	// fmt.Println("timestamp:", timestamp)
	// fmt.Println("filesRecentProcHash:", filesRecentProcHash)
	// fmt.Println("currentGit:", currentGit)
	// fmt.Println("recentGit:", recentGit)
	// fmt.Println("reverse:", reverse)
	// fmt.Println("limit:", limit)
	// fmt.Println("recentProcHash:", recentProcHash)
	// fmt.Println("namesonly:", namesonly)
	// fmt.Println("absolute:", absolute)
	// fmt.Println("remove:", removeExtension)
}

func getTime(milliseconds bool) string {
	milli := ""
	currentTime := time.Now()
	if milliseconds {
		milli = fmt.Sprintf("_%s", currentTime.Format(".000000")[1:])
	}
	time := fmt.Sprintf("%s%s", currentTime.Format("20060102-150405"), milli)
	return time
}

func getHash(data string) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
	return string(hash[:8])
}
