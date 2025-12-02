package cmd

import (
	"fmt"
	"github.com/chingiz/mobwiz/internal/config"
	"github.com/chingiz/mobwiz/internal/generator"
	"github.com/chingiz/mobwiz/internal/prompt"

	"github.com/spf13/cobra"
)

var (
	moduleName   string
	platform     string
	architecture string
	packageName  string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new module",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		var err error

		// If flags are provided, use them; otherwise, run interactive mode
		if moduleName != "" {
			cfg.Module.Name = moduleName
			cfg.Project.Platform = platform
			cfg.Project.Architecture = architecture
			cfg.Options.StateManagement = "BLoC"
			cfg.Options.IncludeTests = true
			cfg.Options.IncludeDI = true

			// Set package name for Android (default to com.example if not provided)
			if platform == "Android (Kotlin)" || platform == "Android" {
				if packageName == "" {
					cfg.Project.PackageName = "com.example"
				} else {
					cfg.Project.PackageName = packageName
				}
			}
		} else {
			fmt.Println("Starting interactive mode...")
			cfg, err = prompt.RunInteractiveFlow(GetEmbeddedFS())
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		fmt.Printf("Generating module '%s' for %s...\n", cfg.Module.Name, cfg.Project.Platform)

		// Generate module using generic generator
		if err := generator.GenerateModule(cfg, GetEmbeddedFS()); err != nil {
			fmt.Printf("Generation failed: %v\n", err)
			return
		}

		fmt.Println("✓ Module generated successfully!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&moduleName, "name", "n", "", "Module name")
	createCmd.Flags().StringVarP(&platform, "platform", "p", "Flutter", "Platform (Flutter, Android, iOS)")
	createCmd.Flags().StringVarP(&architecture, "architecture", "a", "Clean Architecture + BLoC", "Architecture pattern")
	createCmd.Flags().StringVarP(&packageName, "package", "k", "", "Package name (Android only, defaults to com.example)")
}
