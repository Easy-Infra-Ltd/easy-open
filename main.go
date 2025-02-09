package main

import (
	"encoding/json"
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
	default: // "linux", "freebsd", "openbsd", "netbsd"
		// Check if running under WSL
		if isWSL() {
			// Use 'cmd.exe /c start' to open the URL in the default Windows browser
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			// Use xdg-open on native Linux environments
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux
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
					parsedCmd = strings.Replace(parsedCmd, ":"+strconv.Itoa(i+1), url.QueryEscape(params[i]), -1)
				}
			} else {
				parsedCmd = strings.Replace(parsedCmd, ":1", "", -1)
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

	cmds := []Command{}
	fileBytes, _ := os.ReadFile(dirname + "/.config/easyopen.cmds.json")
	jErr := json.Unmarshal(fileBytes, &cmds)
	if jErr != nil {
		panic(jErr)
	}

	parsedCmd := parseCommand(arg, cmds, params)
	rawUri := "https://" + parsedCmd
	_, uriErr := url.ParseRequestURI(rawUri)
	if uriErr != nil {
		panic(uriErr)
	}

	fmt.Printf("Opening Url: %s \n", rawUri)
	openURL(rawUri)
}
