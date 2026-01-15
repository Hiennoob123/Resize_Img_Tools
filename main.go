package main 

import (
	"fmt"
	"path/filepath"
	"github.com/davidbyttow/govips/v2/vips"
	"strings"
	"os"
	"errors"
)

var Res RES = p720
var dstdir, srcdir string

func dst_path(path string) string {
	filename := filepath.Base(path)
	return filepath.Join(dstdir, filename)
}

func check_filename(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".jpg" || ext == "jpeg" {
		return true
	}
	return false
}

func getPaths(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return make([]string, 0), err
	}
	entries := make([]string, len(files))
	for _, file := range files {
		if file.IsDir() || !check_filename(file.Name()) {
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
func parsing(args []string) error {
	if len(args) < 3 {
		return errors.New("Not enough arguments")
	}
	if !checkdir(args[1]) || !checkdir(args[2]) {
		return errors.New("Failed to parse argument")
	}
	srcdir = args[1]
	dstdir = args[2]
	if len(args) == 4 {
		if args[3] == "p1080" {
			Res = p1080
		} else if args[3] == "p720" {
			Res = p720
		} else {
			return errors.New("Failed to parse argument")
		}
	}
	if len(args) > 4 {
		return errors.New("Too many arguments")
	}
	return nil
}

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()
	fmt.Println("Hello Word")
	
	args := os.Args
	fmt.Println(len(args))
	err := parsing(args)
	if err != nil {
		fmt.Println("Error in: ", err)
		return
	}
	files, err := getPaths(srcdir)
	if err != nil {
		fmt.Println("Error in: ", err)
	}
	jobs := NewJob_pool(files)
	jobs.wait_empty()
	_, finished, failed := jobs.stats() 
	fmt.Println("FINISHED: %d tasks", len(finished))
	fmt.Println("FAILED: %d tasks", len(failed))
}
