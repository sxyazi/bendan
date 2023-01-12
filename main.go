package main

import (
	"github.com/sxyazi/bendan/boot"
	"github.com/sxyazi/bendan/db"
)

func main() {
	go db.Indexes()
	boot.ServePool()
}
