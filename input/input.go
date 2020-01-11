package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	mt "github.com/cosandr/go-beat-playlist/types"
)

// GetInputNumber returns first valid number from user input
func GetInputNumber() int {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil || num < 0 {
			fmt.Printf("%s is not a valid number, try again: ", scanner.Text())
			continue
		}
		return num
	}
	return 0
}

// GetInputPlaylist returns complete path
func GetInputPlaylist(dirPath string) (path string, exists bool) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter playlist file: ")
	for scanner.Scan() {
		path = dirPath + "/" + scanner.Text()
		exists = mt.FileExists(path)
		return
	}
	return
}

// GetConfirm reads y/n answer and returns boolean, defaults to true (empty returns true)
func GetConfirm(question string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(question)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			return true
		} else if strings.ToLower(scanner.Text()) == "y" {
			return true
		} else if strings.ToLower(scanner.Text()) == "n" {
			return false
		} else {
			fmt.Printf("%s is not a valid response, try again: ", scanner.Text())
			continue
		}
	}
	return false
}
