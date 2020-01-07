package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// GetInputNumber returns first valid number from user input
func GetInputNumber() int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter number: ")
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
