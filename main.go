package main

import (
	"github.com/ebadfd/jira_sucks/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = bootstrap.RootApp.Execute()
}
