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

var (
	BLACKLISTED_FILE_EXTENSIONS = []string{".url"}
	BLACKLISTED_FOLDER_NAMES    = []string{"node_modules"}
)

func main() {
	flag.StringVar(&dir, "dir", "", "Absolute path for target directory")

	flag.Parse()
	if len(dir) < 1 {
		log.Fatal("No --dir is given")
	}

	Scan(dir)
}

func Scan(root string) {
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		size := fileInfo.Size()
		if size > 104857600 { // 1024 * 1024
			fmt.Println(d.Name(), size)
		}

		if d.IsDir() {
			for _, BLACKLISTED_FOLDER_NAME := range BLACKLISTED_FOLDER_NAMES {
				if d.Name() == BLACKLISTED_FOLDER_NAME {
					fmt.Println("> ", d.Name())
				}
			}
		}

		for _, BLACKLISTED_FILE_EXTENSION := range BLACKLISTED_FILE_EXTENSIONS {
			if filepath.Ext(d.Name()) == BLACKLISTED_FILE_EXTENSION {
				fmt.Println("> ", d.Name())
			}
		}

		return nil
	})
}
