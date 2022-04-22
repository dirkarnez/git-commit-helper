package main

import (
	"flag"
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

func Scan(root string, callback func(s string, d fs.DirEntry) error) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
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
					c.Println("blacklisted folder: ", s)
				}
			}
		}

		for _, BLACKLISTED_FILE_EXTENSION := range BLACKLISTED_FILE_EXTENSIONS {
			if filepath.Ext(d.Name()) == BLACKLISTED_FILE_EXTENSION {
				c.Println("file with blacklisted extension: ", s)
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
	BLACKLISTED_FILE_EXTENSIONS = []string{".url", ".log", ".DS_Store"}
	BLACKLISTED_FOLDER_NAMES    = []string{"node_modules", ".tmp.driveupload", ".git", "vendor", "tmp", "__MACOSX"}
)

func main() {
	flag.StringVar(&dir, "dir", "", "Absolute path for target directory")

	flag.Parse()
	if len(dir) < 1 {
		log.Fatal("No --dir is given")
	}

	ic, _ := inifile.NewIniConfigFromPath(filepath.Join(dir, ".git", "config"))
	originURL, _ := ic.Value("remote \"origin\"", "url")
	if len(originURL) > 0 {
		u, _ := url.Parse(originURL)
		parts := strings.Split(u.Hostname(), ".")
		domain := parts[len(parts)-2] + "." + parts[len(parts)-1]

		switch {
		case domain == "github.com":
			c.Println("GitHub repository detected...")

			Scan(dir, func(s string, d fs.DirEntry) error {
				fileInfo, err := d.Info()
				if err != nil {
					return err
				}

				size := fileInfo.Size()
				if size > 104857600 { // 1024 * 1024
					c.Println("big file: ", s)
				}
				return nil
			})

		case domain == "bitbucket.com":
			Scan(dir, func(s string, d fs.DirEntry) error {
				return nil
			})
		}
	}

}
