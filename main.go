package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/graniticio/inifile"
)

var (
	dir string
	c   = color.New(color.FgCyan)
)

func Scan(root string, callback func(s string, d fs.DirEntry) error) error {
	return filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			c.Printf("prevent panic by handling failure accessing a path %q: %v\n", s, e)
			return e
		}

		// if d.IsDir() && d.Name() == ".git" {
		// 	fmt.Printf("skipping a dir without errors: %+v \n", d.Name())
		// 	return filepath.SkipDir
		// }

		if d.IsDir() {
			for _, BLACKLISTED_FOLDER_NAME := range BLACKLISTED_FOLDER_NAMES {
				if d.Name() == BLACKLISTED_FOLDER_NAME {
					return fmt.Errorf("blacklisted folder: %s", s)
				}
			}

			for _, WARNING_FOLDER_NAME := range BLACKLISTED_FOLDER_NAMES {
				if d.Name() == WARNING_FOLDER_NAME {
					c.Println("warning folder: ", s)
				}
			}
		} else {
			for _, BLACKLISTED_FILE_EXTENSION := range BLACKLISTED_FILE_EXTENSIONS {
				if filepath.Ext(d.Name()) == BLACKLISTED_FILE_EXTENSION {
					return fmt.Errorf("file with blacklisted extension: %s", s)
				}
			}

			for _, WARNING_FILE_EXTENSION := range BLACKLISTED_FILE_EXTENSIONS {
				if filepath.Ext(d.Name()) == WARNING_FILE_EXTENSION {
					c.Println("file with warning extension: ", s)
				}
			}
		}

		err := callback(s, d)
		if err != nil {
			return err
		}

		return nil
	})
}

var (
	WARNING_FILE_EXTENSIONS     = []string{".url"}                                                          // maybe useful in some case
	WARNING_FOLDER_NAMES        = []string{".git"}                                                          // maybe its .git folder
	BLACKLISTED_FILE_EXTENSIONS = []string{".log", ".DS_Store"}                                             // to be deleted
	BLACKLISTED_FOLDER_NAMES    = []string{"node_modules", ".tmp.driveupload", "vendor", "tmp", "__MACOSX"} // to be deleted
)

func main() {
	flag.StringVar(&dir, "dir", "", "Absolute path for target directory")

	flag.Parse()
	if len(dir) < 1 {
		log.Fatal("No --dir is given")
	}

	var domain = "github.com"

	ic, _ := inifile.NewIniConfigFromPath(filepath.Join(dir, ".git", "config"))
	if ic != nil {
		originURL, _ := ic.Value("remote \"origin\"", "url")
		if len(originURL) > 0 {
			u, _ := url.Parse(originURL)
			parts := strings.Split(u.Hostname(), ".")
			domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
		}
	}

	var err error
	switch {
	case domain == "github.com":
		c.Println("GitHub repository detected...")

		err = Scan(dir, func(s string, d fs.DirEntry) error {
			fileInfo, err := d.Info()
			if err != nil {
				return err
			}

			size := fileInfo.Size()
			if size > 104857600 { // 1024 * 1024
				return fmt.Errorf("big file: %s", s)
			}
			return nil
		})

	case domain == "bitbucket.com":
		err = Scan(dir, func(s string, d fs.DirEntry) error {
			return nil
		})
	}

	if err != nil {
		log.Fatal(err)
	}
}
