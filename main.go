package fn

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var gitSize = 8
var pidSize = 8
var sep = "-"

/// GetTime(milliseconds bool) is function for getting current time with ot without milliseconds.
/// Example:
/// 20200927-201950
func GetTime(milliseconds bool) string {
	milli := ""
	currentTime := time.Now()
	if milliseconds {
		milli = fmt.Sprintf("_%s", currentTime.Format(".000000")[1:])
	}
	t := fmt.Sprintf("%s%s", currentTime.Format("20060102-150405"), milli)
	return t
}

/// NewFileName is function for getting new filename.
/// Example:
/// 20200927-201950--289b6e2a
func NewFileName(t string) string {
	return fmt.Sprintf("%s-%s-%s", t, GetGitHash(), _getHash(t))
}

/// NewFileName is function for getting new filename.
/// Example:
/// 20200927-201950--289b6e2a
func NewFileNameCustom(t string, gitSize int, pidSize int) string {
	return fmt.Sprintf("%s-%s-%s", t, GetGitHashCustom(gitSize), _getHashCustom(t, pidSize))
}

/// ListFilesWithGitHashCustom is function for filtering files with git hash
func ListFilesWithGitHash(hash string) []string {
	files := ListFilesWithGitHashCustom(hash, sep, gitSize, pidSize)
	return files
}

/// ListFilesWithGitHashCustom is function for filtering files with git hash
func ListFilesWithGitHashCustom(hash string, sep string, gitSize int, pidSize int) []string {
	files := _listFilesFromDir(sep, gitSize, pidSize)
	var result []string
	for _, file := range files {
		if file["_gitsha"] == hash {
			result = append(result, file["_raw"])
		}
	}
	return result
}

/// GetGitHash() is function for getting current git hash with default sha length
func GetGitHash() string {
	return GetGitHashCustom(gitSize)
}

/// GetGitHashCustom(len int) is function for getting current git hash
func GetGitHashCustom(len int) string {
	out, err := exec.Command("git", "rev-parse", "--short="+strconv.Itoa(len), "HEAD").CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func _getHash(data string) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
	return hash[:pidSize]
}

func _getHashCustom(data string, len int) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
	return hash[:len]
}

func _listFilesFromDir(sep string, gitSize int, pidSize int) []map[string]string {
	files := _getFileNameTokenizer(sep, gitSize, pidSize)
	return files
}

func _getDirFromArgs() string {
	directory := "."
	if flag.Arg(0) != "" {
		directory = flag.Arg(0)
	}
	return directory
}

func _getFilesFromDir() []string {
	files, err := ioutil.ReadDir(_getDirFromArgs())
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list
}

func _getFileNameTokenizer(sep string, gitSize int, pidSize int) []map[string]string {
	_groups := strings.Join([]string{"^(?P<_date>[0-9]{8})",
		"(?P<_time>([0-9]{6}(_[0-9]{6})?))",
		fmt.Sprintf("(?P<_gitsha>[0-9a-z]{%d})?", gitSize),
		fmt.Sprintf("(?P<_prochash>[0-9a-z]{%d})", pidSize)}, sep)
	_ends := strings.Join([]string{"(-(?P<_num>[0-9]+))*",
		"(.(?P<_ext>[a-zA-Z0-9_-]*$))*"}, "")
	tokenString := _groups + _ends
	token := regexp.MustCompile(tokenString)
	var list []map[string]string
	for _, file := range _getFilesFromDir() {
		if token.MatchString(file) {
			match := token.FindStringSubmatch(file)
			result := make(map[string]string)
			for i, name := range token.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			result["_raw"] = file
			list = append(list, result)
		}
	}
	return list
}
