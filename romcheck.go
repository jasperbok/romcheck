package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

const USAGE = `romcheck

Usage:
  romcheck <datfile> <romfolder> [--rename] [--delete] [--collection]

Options:
  -h --help     Show this screen.
  --rename      Automatically rename files with matching hash
  --collection  Shows how much of the listed games you have`

const (
	WHITE = "\033[0m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	RED = "\033[31m"
)

var (
	OPT_RENAME bool     = false
	OPT_COLLECTION bool = false
)

func CheckFile(fpath string, hashmap map[string]string) error {
	finfo, err := os.Stat(fpath)
	if err != nil {
		return err
	}

	md5, err := GetFileMd5(fpath)
	if err != nil {
		return err
	}

	realName, known := hashmap[md5]

	if known && OPT_RENAME && finfo.Name() != realName {
		d, _ := filepath.Split(fpath)
		newPath := filepath.Join(d, realName)

		err = os.Rename(fpath, newPath)
		if err != nil {
			return err
		}
	}

	if known {
		if finfo.Name() == realName {
			fmt.Fprintf(os.Stdout, "%s%s%s\n", GREEN, finfo.Name(), WHITE)
		} else {
			fmt.Fprintf(os.Stdout, "%s%s%s (%s)\n", YELLOW, finfo.Name(), WHITE, realName)
		}
	} else {
		fmt.Fprintf(os.Stdout, "%s%s%s\n", RED, finfo.Name(), WHITE)
	}

	return nil
}

func CheckDir(dirpath string, hashmap map[string]string) error {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			fmt.Fprintf(os.Stdout, "%s%s%s\n", YELLOW, f.Name(), WHITE)
		} else {
			CheckFile(filepath.Join(dirpath, f.Name()), hashmap)
		}
	}

	return nil
}

func CheckCollection(dirpath string, hashmap map[string]string) error {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	games := make(map[string]bool)

	for _, f := range files {
		if !f.IsDir() {
			fpath := filepath.Join(dirpath, f.Name())

			md5, err := GetFileMd5(fpath)
			if err != nil {
				return err
			}

			romName, known := hashmap[md5]
			if known {
				games[romName] = true
				delete(hashmap, md5)
			}
		}
	}

	keys := []string{}

	for _, g := range hashmap {
		games[g] = false
	}

	for gameName, _ := range games {
		keys = append(keys, gameName)
	}

	sort.Strings(keys)

	for _, gameName := range keys {
	    owned := games[gameName]
		if owned {
			fmt.Fprintf(os.Stdout, "%s%s%s\n", GREEN, gameName, WHITE)
		} else {
			fmt.Fprintf(os.Stdout, "%s%s%s\n", RED, gameName, WHITE)
		}
	}

	return nil
}

func main() {
	opts, err := docopt.ParseDoc(USAGE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	datfile, err := opts.String("<datfile>")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	romdir, err  := opts.String("<romfolder>")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	rename, err := opts.Bool("--rename")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	OPT_RENAME = rename

	collection, err := opts.Bool("--collection")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	OPT_COLLECTION = collection

	games, err := LoadGamesFromFile(datfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "The following error occurred:\n%s\n", err)
		os.Exit(1)
	}

	hashmap := BuildMd5Map(games)

	if OPT_COLLECTION {
		CheckCollection(romdir, hashmap)
	} else {
		CheckDir(romdir, hashmap)
	}
}
