package main

import (
	"fmt"
	"log"
)
import "github.com/zamN/zounce/config"

func main() {
	fmt.Println("Zounce started.")

	c, err := config.LoadConfig("config/config.toml")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	_, err = config.LoadConfig("config/bad.toml")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	fmt.Printf("Welcome %#v\n", c)
}
