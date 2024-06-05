package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Playlist	string	`json:"playlist"`
}

type Cache struct {
	Playing		int	`json:"playing"`
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
	err := ipcConnection.Connection.Close()
	check(err)
}

func (ipcConnection *IpcConnection) runCommand(options []any) interface{} {
	commandJson, err := json.Marshal(map[string][]any{"command": options})
	check(err)

	_, err = ipcConnection.Connection.Write([]byte(string(commandJson) + "\n"))
	check(err)

	response := make([]byte, 1024)
	n, err := ipcConnection.Connection.Read(response)
	check(err)

	response = []byte(strings.Split(string(response[:n]), "\n")[0])

	var responseJson map[string]interface{}
	err = json.Unmarshal(response, &responseJson)
	check(err)

	if responseJson["error"].(string) != "success" {
		log.Fatal(responseJson["error"].(string))
	}

	return responseJson["data"]
}

func writeToCache () {
	playing := int(ipcConnection.runCommand([]any{"get_property", "playlist-pos"}).(float64))
	if playing == -1 {
		playing = 0
	}
	cache.Playing = playing
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

	if len(os.Args) == 1 {
		panic("No arguments passed")
	}
	arg := os.Args[1]

	playingFlag := !ipcConnection.runCommand([]any{"get_property", "idle-active"}).(bool)
	currentPlaying := cache.Playing

	switch arg {
	case "toggle":
		if playingFlag {
			ipcConnection.runCommand([]any{"stop"})
		} else {
			ipcConnection.runCommand([]any{"loadfile", config.Playlist, "replace", -1, "playlist-start=" + strconv.Itoa(currentPlaying)})
			writeToCache()
		}
	case "next":
		if playingFlag {
			ipcConnection.runCommand([]any{"playlist-next"})
			writeToCache()
		}
	case "previous":
		if playingFlag {
			ipcConnection.runCommand([]any{"playlist-prev"})
			writeToCache()
		}
	case "status":
		text 	:= ""
		class	:= ""

		if playingFlag {
			text = ipcConnection.runCommand([]any{"get_property", "media-title"}).(string)
			if strings.HasPrefix(text, "playlist") {
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
