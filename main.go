package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/brinwiththevlin/pokedexcli/internal/pokecache"
	"github.com/chzyer/readline"
)

func main() {
	commands := getCommands()
	cache := pokecache.NewCache(time.Second * 5)
	cfg := &config{cache: cache, pokedex: make(map[string]pokemon)}

	rl, err := readline.New("Pokedex > ")
	if err != nil {
		log.Fatalf("failed to create readline: %v", err)
	}
	defer rl.Close()

	for {
		text, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if text == "" {
				continue
			}
		} else if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error reading input: %s\n", err)
			continue
		}
		cleaned := cleanInput(text)
		cmd := cleaned[0]
		if command, ok := commands[cmd]; ok {
			err := command.callback(cfg, cleaned[1:]...)
			if err != nil {
				fmt.Printf("Error executing %s command: %s\n", cmd, err)
			}
		} else {
			fmt.Printf("%s is not a command.\n", cmd)
		}

	}

}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}
