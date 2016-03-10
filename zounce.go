package main

import (
	"fmt"
	"log"
)
import (
	"github.com/zamN/zounce/config"
	"github.com/zamN/zounce/server"
)

func main() {
	fmt.Println("Zounce started.")

	c, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	for _, u := range c.Users {
		fmt.Printf("Welcome %s\n", u.Nick)
	}

	server.Start()

}
