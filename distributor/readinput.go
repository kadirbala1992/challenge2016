package distribution

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func GetInput() []string {
	var status string
	var inputLines []string

	fmt.Printf("Permissions for: ")

	for status != "pass" {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			status = "pass"
		} else {
			inputLines = append(inputLines, text)
		}

	}

	return inputLines

}

func GetDistType() string {
	var status string
	var text string

	for status != "pass" {
		fmt.Printf("\nEnter the distributor type (Direct/Indirect): ")
		scanner.Scan()
		text = strings.ToLower(scanner.Text())
		if text == "direct" {
			status = "pass"
		} else if text == "indirect" {
			fmt.Printf("\nEnter the Root - Distributor name: ")
			scanner.Scan()
			text = scanner.Text()
			return text
		}

	}
	return text
}
