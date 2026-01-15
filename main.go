package main 

import (
	"time"
	"fmt"
	"path/filepath"
	"github.com/davidbyttow/govips/v2/vips"
	"strings"
	"os"
	"errors"
	"flag"
	"io/fs"
)

var Res RES = p720
var dstdir, srcdir string
var rec bool

func dst_path(path string) string {
	rel, err := filepath.Rel(srcdir, path)
	if err != nil {
		fmt.Println("Failed to resolve path")
	}
	newpath := filepath.Join(dstdir, rel)
	return newpath 
}

func check_file(file fs.DirEntry) bool {
	if file.IsDir() {
		return false
	}
	ext := strings.ToLower(filepath.Ext(file.Name()))
	if ext == ".jpg" || ext == "jpeg" {
		return true
	}
	return false
}

func getPaths_recursive(dir string) ([] string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("WalkDir prevented at %s: %s\n", path, err)
			return err
		}
		if check_file(d) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}	else {
		return files, nil
	}
}

func getPaths(dir string) ([]string, error) {
	if rec {
		return getPaths_recursive(dir)
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return make([]string, 0), err
	}
	entries := make([]string, 0)
	for _, file := range files {
		if check_file(file) {
			continue
		}
		entries = append(entries, filepath.Join(dir, file.Name()))
	}
	return entries, nil
}

// Check if directory exist, if not then create it
func checkdir(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Path %s not exist", path)
		return false
	}
	fmt.Println("Failed to check path %s: ", path, err)
	return false
}
func flag_parse() error {
	resolution := flag.String("res", "p720", "resolution of the resized image")
	recursive  := flag.Bool("r", false, "recursively process image")
	flag.Parse()
	if *resolution == "p1080" {
		Res = p1080
	} else if *resolution == "p720" {
		Res = p720
	} else {
		return errors.New("wrong format for --res") 
	}

	if *recursive == true {
		rec = true
	} else {
		rec = false
	}
	return nil
}
func parsing(args []string) error {
	if len(args) < 2 {
		return errors.New("2 arguments srcdir and dstdir are required")
	}
	if !checkdir(args[0]) || !checkdir(args[1]) {
		return errors.New("Failed to parse argument")
	}
	srcdir = args[0]
	dstdir = args[1]
	return nil
}
func myLoggerHandler(messageDomain string, verbosity vips.LogLevel, message string) {
	var messageLevelDescription string
	switch verbosity {
	case vips.LogLevelError:
		messageLevelDescription = "error"
	case vips.LogLevelCritical:
		messageLevelDescription = "critical"
	case vips.LogLevelWarning:
		messageLevelDescription = "warning"
	case vips.LogLevelMessage:
		messageLevelDescription = "message"
	case vips.LogLevelInfo:
		messageLevelDescription = "info"
	case vips.LogLevelDebug:
		messageLevelDescription = "debug"
	}
	fmt.Printf("[%v.%v] %v", messageDomain, messageLevelDescription, message)
}
func main() {
	exec_time := time.Now()
	vips.LoggingSettings(myLoggerHandler, vips.LogLevelError)
	vips.Startup(nil)
	defer vips.Shutdown()
	fmt.Println("Start Resizing Images")
	err := flag_parse()
	if err != nil {
		fmt.Println("Error in flag parse", err)
		return
	}
	args := flag.Args()
	err = parsing(args)
	if err != nil {
		fmt.Println("Error in: ", err)
		return
	}
	files, err := getPaths(srcdir)
	if err != nil {
		fmt.Println("Error in: ", err)
	}
	jobs := NewJob_pool(files)
	_, finished, failed := jobs.wait_empty() 
	fmt.Println("FINISHED: %d tasks", len(finished))
	fmt.Println("FAILED: %d tasks", len(failed))
	fmt.Println("Execution time is", time.Since(exec_time))
	fmt.Println("FAILED Jobs:")
	for _, job := range failed {
		fmt.Printf("Job %d %s failed: %s\n", job.id, job.path, job.err)
	}
}
