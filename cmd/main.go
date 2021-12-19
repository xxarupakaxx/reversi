package main

import (
	"github.com/xxarupakaxx/reversi/client"
	"os"
)

func main() {
	os.Exit(client.NewReversi().Run())
}
