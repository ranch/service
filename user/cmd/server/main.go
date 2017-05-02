package main

import (
	"github.com/synoday/service/user"
)

func main() {
	service := user.NewService()

	service.Serve()
}
