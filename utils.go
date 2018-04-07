package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

func GetFileMd5(filename string) (md5Hash string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return md5Hash, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5Hash, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	md5Hash      = strings.ToUpper(hex.EncodeToString(hashInBytes))

	return md5Hash, nil
}
