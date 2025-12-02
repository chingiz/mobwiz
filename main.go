package main

import (
	"embed"
	"github.com/chingiz/mobwiz/cmd"
)

//go:embed templates/*
var embeddedTemplates embed.FS

func main() {
	// Pass the embedded filesystem to the command execution
	
	cmd.Execute(embeddedTemplates)
}