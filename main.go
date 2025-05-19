package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var searchItems = []string{"Free search", "Specific search"}

func main() {
	clearScreen()
	fmt.Println("Looking for internal storage...")
	hasMoved := moveToHomeDir()

	current, err := os.Getwd()
	if err != nil {
		color.Red("Unable to switch paths")
		fmt.Println("Is this termux?")
		return
	}

	time.Sleep(time.Second)
	if hasMoved {
		fmt.Printf("Current: %v\n", current)
	} else {
		fmt.Println("Failed!")
		return
	}

	if !findDir(current + "/storage") {
		color.Red("Storage path not found!")
		fmt.Println("Termux doesen't have access to internal storage.")
		fmt.Println("Use termux-setup-storage to give it access.")
		fmt.Println("Exiting process!")
		return
	}

	time.Sleep(time.Second)
	color.Green("Successfully accessed Internal storage!")

	deepSearchCommandCenter(current + "/storage")

}

// the deep search command center
func deepSearchCommandCenter(path string) {
	clearScreen()
	var prompt = promptui.Prompt{
		Label: "Enter file name or extension to find",
	}

	res, err := prompt.Run()

	if err != nil {
		return
	}

	clearScreen()
	var choiceSelect = promptui.Select{
		Label: "Search type",
		Items: searchItems,
	}

	var _, choice, err2 = choiceSelect.Run()

	if err2 != nil {
		return
	}

	clearScreen()
	fmt.Printf("Looking for %v in storage [%v]\n", res, path)
	var exists, fpath = searchPath(path, res, choice)
	if exists {
		color.Green("File found at : %v\n", fpath)
	} else {
		color.Red("Unable to find file!")
	}
}

func moveToHomeDir() bool {
	home, err := os.UserHomeDir()

	if err != nil {
		color.Red("Error: Could not find the home directory!")
		fmt.Println("Exiting process...")

		return false
	}

	var err2 = os.Chdir(home)
	if err2 != nil {
		color.Red("Failed to change directory!")
		fmt.Println("Ripple is unable to switch to the home directory!")
		fmt.Println("Exiting process...")
		return false
	}

	return true
}

// searches a specified path
func searchPath(path string, filename string, choice string) (bool, string) {
	fmt.Printf("Search path: %v\n", path)
	files, err := os.ReadDir(path)
	if err != nil {
		color.Red("Failed to open directory!")
		return false, ""
	}

	for _, entry := range files {
		fullpath := filepath.Join(path, entry.Name())
		info, err := os.Stat(fullpath)

		if err != nil {
			color.Red("Failed to grab stat!")
			return false, ""
		}
		if info.IsDir() {
			fmt.Println("Opening new dir")
			found, fullpath := searchPath(filepath.Join(path, entry.Name()), filename, choice)
			if found {
				return true, fullpath
			}
		} else {
			if choice == searchItems[1] {
				if entry.Name() == filename {
					return true, filepath.Join(path, filename)
				}
			} else {
				if strings.Contains(entry.Name(), filename) {
					return true, filepath.Join(path, filename)
				}
			}
		}
	}
	return false, ""

}

// clears the terminal screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func findDir(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}
