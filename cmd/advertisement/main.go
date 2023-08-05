package main

import (
	"fmt"
	"github.com/romandnk/advertisement/configs"
	"log"
)

func main() {
	config, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("error initialising config: %s", err.Error())
	}
	fmt.Printf("%+v\n", config)
}
