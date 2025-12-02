package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "mobwiz",
	Short: "Mobile Module Creator",
	Long:  `A CLI tool to generate mobile app modules with consistent architecture patterns.`,
}

var embeddedFS fs.FS

func Execute(fsys fs.FS) {
	embeddedFS = fsys
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetEmbeddedFS() fs.FS {
	return embeddedFS
}
