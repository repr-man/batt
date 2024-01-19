package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
    listBatteries()
}

func listBatteries() {
    // Get power-related devices.
    upowerPath, err := exec.LookPath("upower")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Unable to find upower executable.")
        os.Exit(-1)
    }
    listCmd := exec.Command(upowerPath, "-e")
    rawList, err := listCmd.Output()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to enumerate devices.")
        os.Exit(-1)
    }
    list := strings.Split(string(rawList), "\n")

    for _, device := range list {
        // Find devices that are batteries.
        if strings.Contains(device, "/org/freedesktop/UPower/devices/battery") {
            batteryCmd := exec.Command(upowerPath, "-i", device)
            batteryCmdOut, err := batteryCmd.Output()
            if err != nil {
                fmt.Fprintln(os.Stderr, "Unable to find battery '", device, "'.")
                os.Exit(-1)
            }
            batteryCmdStr := string(batteryCmdOut)

            // Obtain information about devices from output of the command.
            state, batteryCmdStr := findValue(batteryCmdStr, "state:               ")
            var timeLeft string
            if state == "discharging" {
                timeLeft, batteryCmdStr = findValue(batteryCmdStr, "time to empty:       ")
            }
            percentage, _ := findValue(batteryCmdStr, "percentage:          ")

            if timeLeft == "" {
                fmt.Printf("Battery: %s, %s\n", percentage, state)
            } else {
                fmt.Printf("Battery: %s, %s, %s left\n", percentage, state, timeLeft)
            }
            
        }
    }
}

func findValue(cmdTxt string, pattern string) (string, string) {
    start := strings.Index(string(cmdTxt), pattern) + len(pattern)
    cmdTxt = cmdTxt[start:]
    nlIdx := strings.IndexByte(string(cmdTxt), '\n')
    return cmdTxt[:nlIdx], cmdTxt[nlIdx:]
}
