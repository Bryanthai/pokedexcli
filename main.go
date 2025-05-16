package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\033[34mPokedex > \033[0m")
		scanner.Scan()
		inputString := scanner.Text()
		cleanedInput := cleanInput(inputString)
		if len(cleanedInput) == 0{
			fmt.Printf("You have not entered anything...\n")
			continue
		}
		callFunc(cleanedInput)
	}
}