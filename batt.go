package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
    "slices"
)

func main() {
    listBatteries()
}

func listBatteries() {
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

    

    statePattern := "state:               "
    timeLeftPattern := "time to empty:       "
    percentagePattern := "percentage:          "


    for _, device := range list {
        if strings.Contains(device, "/org/freedesktop/UPower/devices/battery") {
            batteryCmd := exec.Command(upowerPath, "-i", device)
            batteryCmdOut, err := batteryCmd.Output()
            if err != nil {
                fmt.Fprintln(os.Stderr, "Unable to find battery '", device, "'.")
                os.Exit(-1)
            }

            stateStart := strings.Index(string(batteryCmdOut), statePattern) + len(statePattern)
            batteryCmdOut = batteryCmdOut[stateStart:]
            state := batteryCmdOut[:strings.IndexByte(string(batteryCmdOut), '\n')]

            var timeLeft []byte
            if slices.Equal(state, []byte("discharging")) {
                timeLeftStart := strings.Index(string(batteryCmdOut), timeLeftPattern) + len(timeLeftPattern)
                batteryCmdOut = batteryCmdOut[timeLeftStart:]
                timeLeft = batteryCmdOut[:strings.IndexByte(string(batteryCmdOut), '\n')]
            }

            percentageStart := strings.Index(string(batteryCmdOut), percentagePattern) + len(percentagePattern)
            batteryCmdOut = batteryCmdOut[percentageStart:]
            percentage := batteryCmdOut[:strings.IndexByte(string(batteryCmdOut), '\n')]

            if timeLeft == nil {
                fmt.Printf("Battery: %s, %s\n", percentage, state)
            } else {
                fmt.Printf("Battery: %s, %s, %s left\n", percentage, state, timeLeft)
            }
            
        }
    }
}
