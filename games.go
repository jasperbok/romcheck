package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Game struct {
	Name string `xml:"name,attr"`
	Rom  Rom    `xml:"rom"`
}

type Rom struct {
	Name   string `xml:"name,attr"`
	Size   int    `xml:"size,attr"`
	Md5    string `xml:"md5,attr"`
	Status string `xml:"status,attr"`
}

type Result struct {
	Game []Game `xml:"game"`
}

func LoadGamesFromFile(filename string) (games []Game, err error) {
	games = []Game{}

	file, err := os.Open(filename)
	if err != nil {
		return games, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return games, err
	}

	result := Result{}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		return games, err
	}

	return result.Game, nil
}

func BuildMd5Map(games []Game) map[string]string {
	m := make(map[string]string)

	for _, g := range games {
		m[g.Rom.Md5] = g.Rom.Name
	}

	return m
}
