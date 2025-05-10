package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
	default:
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}

type Command struct {
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Params []string `json:"params"`
}

func parseCommand(arg string, cmds []Command, params []string) string {
	parsedCmd := arg
	for _, cmd := range cmds {
		if cmd.Name == parsedCmd {
			parsedCmd = cmd.Url

			if len(params) > 0 {
				for i := 0; i < len(params); i++ {
					parsedCmd = strings.ReplaceAll(parsedCmd, ":"+strconv.Itoa(i+1), url.QueryEscape(params[i]))
				}
			} else {
				parsedCmd = strings.ReplaceAll(parsedCmd, ":1", "")
			}
		}
	}

	return parsedCmd
}

func usage() {
	fmt.Printf("Usage: %s url|command [...params]\n", os.Args[0])
}

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	arg := os.Args[1]
	params := os.Args[2:]

	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(dirname + "/.config/easyopen.cmds.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s/.config/easyopen.cmds.json does not exist, please ensure this is setup", dirname)
	} else if err != nil {
		panic(err)
	}

	cmds := []Command{}
	fileBytes, _ := os.ReadFile(dirname + "/.config/easyopen.cmds.json")
	if err := json.Unmarshal(fileBytes, &cmds); err != nil {
		panic(err)
	}

	parsedCmd := parseCommand(arg, cmds, params)
	rawUri := "https://" + parsedCmd
	if _, err := url.ParseRequestURI(rawUri); err != nil {
		panic(err)
	}

	fmt.Printf("Opening Url: %s \n", rawUri)
	if err := openURL(rawUri); err != nil {
		panic(err)
	}
}
