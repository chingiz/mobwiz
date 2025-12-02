package prompt

import (
	"errors"
	"io/fs"
	"github.com/manifoldco/promptui"
	"github.com/chingiz/mobwiz/internal/config"
	"github.com/chingiz/mobwiz/internal/generator"
)

// RunInteractiveFlow gathers user input for module creation
func RunInteractiveFlow(embeddedFS fs.FS) (config.Config, error) {
	var cfg config.Config

	// 1. Select Platform
	platformPrompt := promptui.Select{
		Label: "Select platform",
		Items: []string{"Flutter", "Android (Kotlin)", "iOS (Swift)"},
	}
	_, platform, err := platformPrompt.Run()
	if err != nil {
		return cfg, err
	}
	cfg.Project.Platform = platform

	// 2. Module Name
	namePrompt := promptui.Prompt{
		Label: "Module name",
		Validate: func(input string) error {
			if len(input) < 3 {
				return errors.New("module name must be at least 3 characters")
			}
			return nil
		},
	}
	name, err := namePrompt.Run()
	if err != nil {
		return cfg, err
	}
	cfg.Module.Name = name

	// 3. Architecture Pattern (Contextual based on platform)
	tmplConfig, err := generator.LoadTemplatesConfig(embeddedFS)
	var archItems []string
	if err != nil {
		// Fallback if config fails to load (though this shouldn't happen in a valid setup)
		if platform == "Flutter" {
			archItems = []string{"Clean Architecture + BLoC"}
		} else if platform == "Android (Kotlin)" {
			archItems = []string{"MVVM + Clean Architecture"}
		} else {
			archItems = []string{"MVVM + Coordinators"}
		}
	} else {
		archItems = tmplConfig.GetArchitectures(platform)
		if len(archItems) == 0 {
			archItems = []string{"Default"}
		}
	}

	archPrompt := promptui.Select{
		Label: "Architecture pattern",
		Items: archItems,
	}
	_, arch, err := archPrompt.Run()
	if err != nil {
		return cfg, err
	}
	cfg.Project.Architecture = arch

	// 4. Package Name (for Android)
	if platform == "Android (Kotlin)" {
		packagePrompt := promptui.Prompt{
			Label:   "Package name",
			Default: "com.example",
			Validate: func(input string) error {
				if len(input) < 3 {
					return errors.New("package name must be at least 3 characters")
				}
				return nil
			},
		}
		packageName, err := packagePrompt.Run()
		if err != nil {
			return cfg, err
		}
		cfg.Project.PackageName = packageName
	}

	// 5. State Management (for Flutter)
	if platform == "Flutter" {
		smPrompt := promptui.Select{
			Label: "State management",
			Items: []string{"BLoC", "Riverpod", "Provider", "GetX"},
		}
		_, sm, err := smPrompt.Run()
		if err != nil {
			return cfg, err
		}
		cfg.Options.StateManagement = sm
	}

	// 6. Features (Simplified for now, just yes/no for common ones)
	cfg.Options.IncludeTests = true
	cfg.Options.IncludeDI = true

	return cfg, nil
}
