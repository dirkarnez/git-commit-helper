package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

var (
	dir string
)

func Scan(root string) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", s, e)
			return e
		}

		// if d.IsDir() && d.Name() == ".git" {
		// 	fmt.Printf("skipping a dir without errors: %+v \n", d.Name())
		// 	return filepath.SkipDir
		// }

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		size := fileInfo.Size()
		if size > 104857600 { // 1024 * 1024
			fmt.Println("big file: ", s)
		}

		if d.IsDir() {
			for _, BLACKLISTED_FOLDER_NAME := range BLACKLISTED_FOLDER_NAMES {
				if d.Name() == BLACKLISTED_FOLDER_NAME {
					fmt.Println("blacklisted folder: ", s)
				}
			}
		}

		for _, BLACKLISTED_FILE_EXTENSION := range BLACKLISTED_FILE_EXTENSIONS {
			if filepath.Ext(d.Name()) == BLACKLISTED_FILE_EXTENSION {
				fmt.Println("file with blacklisted extension: ", s)
			}
		}

		return nil
	})
}

var (
	BLACKLISTED_FILE_EXTENSIONS = []string{".url", ".log"}
	BLACKLISTED_FOLDER_NAMES    = []string{"node_modules", ".tmp.driveupload", ".git", "vendor", "tmp"}
)

func main() {
	flag.StringVar(&dir, "dir", "", "Absolute path for target directory")

	flag.Parse()
	if len(dir) < 1 {
		log.Fatal("No --dir is given")
	}

	Scan(dir)
}
