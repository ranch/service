package main

import (
	"github.com/synoday/service/task"
)

func main() {
	service := task.NewService()

	service.Serve()
}
