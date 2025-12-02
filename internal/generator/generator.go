package generator

import (
	"fmt"
	"io/fs"
	"github.com/chingiz/mobwiz/internal/config"
)

// GenerateModule generates the files for a module based on configuration
func GenerateModule(cfg config.Config, embeddedFS fs.FS) error {
	engine := NewEngine(cfg)
	
	// Load template configuration
	tmplConfig, err := LoadTemplatesConfig(embeddedFS)
	if err != nil {
		return fmt.Errorf("failed to load template configuration: %w", err)
	}

	// Get templates for the selected platform and architecture
	templates, platformKey, archKey, err := tmplConfig.GetTemplatesFor(cfg.Project.Platform, cfg.Project.Architecture)
	if err != nil {
		return err
	}

	// Prepare data for templates
	packageName := cfg.Project.PackageName
	if packageName == "" {
		packageName = "com.example"
	}

	data := struct {
		Name        string
		PackageName string
	}{
		Name:        cfg.Module.Name,
		PackageName: packageName,
	}

	fmt.Printf("Generating %s module '%s' using %s architecture...\n", cfg.Project.Platform, cfg.Module.Name, cfg.Project.Architecture)

	for _, tmplDef := range templates {
		// Load template content from file
		content, err := LoadTemplate(platformKey, archKey, tmplDef.TemplateFile, embeddedFS)
		if err != nil {
			return err
		}

		// Render output path
		path, err := engine.RenderPath(tmplDef.OutputPath, data)
		if err != nil {
			return fmt.Errorf("failed to render path %s: %w", tmplDef.OutputPath, err)
		}

		// Render content
		renderedContent, err := engine.Render(content, data)
		if err != nil {
			return fmt.Errorf("failed to render content for %s: %w", path, err)
		}

		// Write file
		if err := engine.WriteFile(path, renderedContent); err != nil {
			return err
		}
		fmt.Printf("Created %s\n", path)
	}

	return nil
}
