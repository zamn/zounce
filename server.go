package main

import "fmt"
import "github.com/zamN/zounce/config"

func main() {
	fmt.Println("Zounce started.")
	config.LoadConfig("config.toml")
}
