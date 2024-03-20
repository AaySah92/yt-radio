package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Playlist	[]string	`json:"playlist"`
}

type Cache struct {
	Playing		int		`json:"playing"`
}

func readFromFile[T Config | Cache] (objectPath string, object *T) {
	data, err := os.ReadFile(objectPath)
	if errors.Is(err, fs.ErrNotExist) {
		writeToFile(objectPath, object)
		return
	} else {
		check(err)
	}

	err = json.Unmarshal(data, object)
	check(err)
}

func writeToFile[T Config | Cache] (objectPath string, object *T) {
	data, err := json.Marshal(*object)
	check(err)

	dir := path.Dir(objectPath)
	err = os.MkdirAll(dir, 0755)
	check(err)

	err = os.WriteFile(objectPath, []byte(data), 0644)
	check(err)
}

type IpcConnection struct {
	Connection	net.Conn
}

func (ipcConnection *IpcConnection) openConnection() {
	connection, err := net.Dial("unix", "/tmp/mpvsocket")
	check(err)
	ipcConnection.Connection = connection
}

func (ipcConnection *IpcConnection) closeConnection() {
	ipcConnection.Connection.Close()
}

func (ipcConnection *IpcConnection) sendCommand(options []string) {
	commandJson, err := json.Marshal(map[string][]string{"command": options})
	check(err)

	_, err = ipcConnection.Connection.Write([]byte(string(commandJson) + "\n"))
	check(err)
}

func (ipcConnection *IpcConnection) getCommand(options []string) interface{} {
	ipcConnection.sendCommand(options)
	response := make([]byte, 1024)
	n, err := ipcConnection.Connection.Read(response)
	check(err)

	var responseJson map[string]interface{}
	err = json.Unmarshal(response[:n], &responseJson)
	check(err)

	if responseJson["error"].(string) != "success" {
		panic(responseJson["error"].(string))
	}

	return responseJson["data"]
}

func play (currentPlaying int) {
	ipcConnection.sendCommand([]string{"loadfile", config.Playlist[currentPlaying]})
	cache.Playing = currentPlaying
	writeToFile(cachePath, &cache)
}

var configPath		string	= "/waybar/yt-radio/config.json"
var cachePath		string	= "/yt-radio/yt-radio-cache"

var config		Config
var cache		Cache
var ipcConnection	IpcConnection 

func main() {
	ipcConnection.openConnection()
	defer ipcConnection.closeConnection()

	userConfigDir, err := os.UserConfigDir()
	check(err)

	userCacheDir, err := os.UserCacheDir()
	check(err)

	configPath	= userConfigDir	+ configPath
	cachePath	= userCacheDir	+ cachePath

	readFromFile(configPath, &config)
	readFromFile(cachePath, &cache)

	if len(config.Playlist) == 0 {
		panic("No tracks in playlist")
	}

	if len(os.Args) == 0 {
		panic("No arguments passed")
	}
	arg := os.Args[1]

	playingFlag := !ipcConnection.getCommand([]string{"get_property", "idle-active"}).(bool)
	currentPlaying := cache.Playing

	switch arg {
	case "toggle":
		if playingFlag {
			ipcConnection.sendCommand([]string{"stop"})
		} else {
			play(currentPlaying)
		}
	case "next":
		if playingFlag {
			currentPlaying++
			if currentPlaying >= len(config.Playlist) {
				currentPlaying = 0
			}
			play(currentPlaying)
		}
	case "previous":
		if playingFlag {
			currentPlaying--
			if currentPlaying < 0 {
				currentPlaying = len(config.Playlist) - 1
			}
			play(currentPlaying)
		}
	case "status":
		text 	:= ""
		class	:= ""

		if playingFlag {
			text = ipcConnection.getCommand([]string{"get_property", "media-title"}).(string)
			if strings.HasPrefix(text, "watch") {
				text = "Connecting"
			} else if len(text) > 30 {
				text = text[:30]
			}

			text += "..."
			class = "playing"
		}

		responseJson, err := json.Marshal(map[string]string{"text": text, "class": class})
		check(err)
		fmt.Println(string(responseJson))
	}
}
