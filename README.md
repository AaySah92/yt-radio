# YT-Radio - a custom [waybar](https://github.com/Alexays/Waybar/wiki) module

YT-Radio is a custom waybar module written in [Go](https://go.dev/). It's like a YouTube radio. You can have music playing in the background as you work. It's literally as easy as clicking  a button.

A list of YouTube links (your playlist) can be set in a JSON based configuration file.

## Screenshots

![ezgif-7-aae398a7c8](https://github.com/AaySah92/yt-radio/assets/26904734/40d55d63-2bcf-415c-9158-fd971e07ef48)

#### Demo on YouTube:

[![Demo Video](https://img.youtube.com/vi/XWXR4vf2Xng/0.jpg)](https://www.youtube.com/watch?v=XWXR4vf2Xng)

## Features

- Clicking the YT button starts / stops the player.
-  Scrolling up on the YT button goes to the next track.
- Scrolling down on the YT button goes to the previous track.
- The title of the current track is shown next to the YT button.

## Dependencies
- [mpv](https://mpv.io/installation/) - a free (as in freedom) media player for the command line.
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - a feature-rich command-line audio/video downloader.

## Installation
The installtion involves 3 steps:
- Download the binary
- Add a custom waybar module
- Add style

#### Download the binary
```bash
mkdir -p $XDG_CONFIG_HOME/waybar/yt-radio && curl -o $XDG_CONFIG_HOME/waybar/yt-radio/yt -L https://raw.githubusercontent.com/AaySah92/yt-radio/main/yt 
```

#### Add a custom waybar module
You can configure waybar as you like.
Here's a sample:
```json
"modules-right": [
    "group/yt-radio",
],
"group/yt-radio": {
    "modules": [
        "custom/yt-radio-title",
        "custom/yt-radio-button",
    ],
    "orientation": "inherit",
},
"custom/yt-radio-title": {
    "tooltip": false,
    "return-type": "json",
    "interval": 1,
    "exec": "$HOME/.config/waybar/yt-radio/yt status",
    "signal": 10,
},
"custom/yt-radio-button": {
    "format": "",
    "on-click": "$HOME/.config/waybar/yt-radio/yt toggle; pkill -SIGRTMIN+10 waybar",
    "on-scroll-up": "$HOME/.config/waybar/yt-radio/yt next; pkill -SIGRTMIN+10 waybar",
    "on-scroll-down": "$HOME/.config/waybar/yt-radio/yt previous; pkill -SIGRTMIN+10 waybar",
    "tooltip": false,
    "smooth-scrolling-threshold": 13,
    "return-type": "json",
    "interval": "once",
    "exec": "$HOME/.config/waybar/yt-radio/yt status",
    "signal": 10,
    "exec-on-event": false,
},

```

#### Add style
You can style the module as you like. Here's a sample:
```css
#custom-yt-radio-title {
    margin-right: 0px;
    border-top-right-radius: 0px;
    border-bottom-right-radius: 0px;
    background-color: @color7;
    border-color: @color7;
    color: #000000;
    min-width: 100px;
}

#custom-yt-radio-button {
    margin-left: 0px;
    font-size: 30px;
    padding: 2px 15px;
    background-color: #990000; 
    border-color: #990000;
}

#custom-yt-radio-button.playing {
    border-top-left-radius: 0px;
    border-bottom-left-radius: 0px;
}

```

```bash
git clone https://github.com/AaySah92/yt-radio.git $XDG_CONFIG_HOME/
```
    
## Configuration

#### mpv
Make sure mpv is running in the background with the following options:
```bash
mpv --no-terminal --no-video --idle=yes --input-ipc-server=/tmp/mpvsocket
```

#### Playlist
- File: $XDG_CONFIG_HOME/waybar/yt-radio/config.json
- Format: JSON
- Sample: 
```json
{
	"playlist": [
		"https://www.youtube.com/watch?v=4xDzrJKXOOY",
		"https://www.youtube.com/watch?v=jfKfPfyJRdk"
	]
}

```
If the config file is not created before the module is loaded for the first time, it will be created with empty values. You can add in your YouTube links later. No need to restart waybar.
## Usage

```bash
$XDG_CONFIG_HOME/yt-radio/yt <command>
```

| Command	| Description	| Return	|
| -------	| -----------	| ------	|
| Toggle	| Play / Stop	| NA		|
| Next		| Play next		| NA		|
| Previous	| Play previous	| NA		|
| Status	| Get status	| {"text": "$text", "class": "$class"}	|

## Build
You can edit the file `yt.go` and complie it using Go. This will generate a binary `yt` which can be used by waybar.
```bash
go build yt.go
```
