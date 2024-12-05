package main

import (
	"bufio"
	"fmt"
	"lierslie/client"
	"lierslie/controllers"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Welcome to Liars Lie!")

	for {
		fmt.Println("Enter command:")

		reader := bufio.NewReader(os.Stdin)
		inputCommand, _ := reader.ReadString('\n')

		inputCommand = strings.TrimSpace(inputCommand)
		args := strings.Fields(inputCommand)

		// Check if input is empty
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "start":
			if len(args) < 9 {
				fmt.Println("Usage: start --value <value> --max-value <max> --num-agents <numAgents> --liar-ratio <liarRatio>")
				continue
			}
			// Reset game state without terminating program
			controllers.StopGame(false)
			var value, max, numAgents int
			var liarRatio float64
			for i := 1; i < len(args); i++ {
				switch args[i] {
				case "--value":
					value, _ = strconv.Atoi(args[i+1])
				case "--max-value":
					max, _ = strconv.Atoi(args[i+1])
				case "--num-agents":
					numAgents, _ = strconv.Atoi(args[i+1])
				case "--liar-ratio":
					liarRatio, _ = strconv.ParseFloat(args[i+1], 64)
				}
			}
			controllers.StartGame(value, max, numAgents, liarRatio)

		case "kill":
			// Check if the game has been started
			if !controllers.CheckConfigExists() {
				fmt.Println("Game not started. Please start the game first.")
				continue
			}

			if len(args) != 3 {
				fmt.Println("Usage: kill --id <id>")
				continue
			}
			id, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("Invalid ID:", err)
				continue
			}
			controllers.KillAgentByID(id, true) //treu for Prints confirmation message on kill

		case "stop":
			if !controllers.CheckConfigExists() {
				fmt.Println("Game not started. Please start the game first.")
				continue
			}
			controllers.StopGame(true)

		case "play":
			if !controllers.CheckConfigExists() {
				fmt.Println("Game not started. Please start the game first.")
				continue
			}
			client := &client.Client{}
			client.ReadConfig()
			fmt.Println("Network value:", client.Play())

		case "extend":
			if !controllers.CheckConfigExists() {
				fmt.Println("Game not started. Please start the game first.")
				continue
			}
			if len(args) < 9 {
				fmt.Println("Usage: extend --value <value> --max-value <max> --num-agents <numAgents> --liar-ratio <liarRatio>")
				continue
			}
			var value, max, numAgents int
			var liarRatio float64
			for i := 1; i < len(args); i++ {
				switch args[i] {
				case "--value":
					value, _ = strconv.Atoi(args[i+1])
				case "--max-value":
					max, _ = strconv.Atoi(args[i+1])
				case "--num-agents":
					numAgents, _ = strconv.Atoi(args[i+1])
				case "--liar-ratio":
					liarRatio, _ = strconv.ParseFloat(args[i+1], 64)
				}
			}
			controllers.ExtendNetwork(value, max, numAgents, liarRatio)

		case "playexpert":
			if !controllers.CheckConfigExists() {
				fmt.Println("Game not started. Please start the game first.")
				continue
			}
			client := &client.Client{}
			client.ReadConfig()

			if len(args) < 5 {
				fmt.Println("Usage: playexpert --num-agents <number> --liar-ratio <ratio>")
				continue
			}

			var numAgents int
			var liarRatio float64
			for i := 1; i < len(args); i++ {
				switch args[i] {
				case "--num-agents":
					numAgents, _ = strconv.Atoi(args[i+1])
				case "--liar-ratio":
					liarRatio, _ = strconv.ParseFloat(args[i+1], 64)
				}
			}
			fmt.Println("Network value (expert mode):", client.PlayExpert(numAgents, liarRatio))

		default:
			fmt.Println("Unknown command:", args[0])
		}
	}
}
