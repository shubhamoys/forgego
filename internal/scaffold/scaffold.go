package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shubhamoys/forgego/internal/pkg/constants"
	"github.com/shubhamoys/forgego/internal/prompts"
	"github.com/shubhamoys/forgego/internal/utils"
)

func ScaffoldProject(config *prompts.ProjectConfig) error {
	// Create project directory
	projectDir := config.ProjectName
	if err := utils.CreateDir(projectDir); err != nil {
		return err
	}

	// Generate folder structure based on project type
	switch config.ProjectType {
	case constants.ProjectTypeAPI:
		if err := scaffoldAPI(projectDir, config); err != nil {
			return err
		}
	case constants.ProjectTypeCLI:
		if err := scaffoldCLI(projectDir, config); err != nil {
			return err
		}
	case constants.ProjectTypeLibrary:
		if err := scaffoldLibrary(projectDir, config); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported project type: %s", config.ProjectType)
	}

	// Write common files
	if err := renderTemplate("common/go.mod.tmpl", filepath.Join(projectDir, "go.mod"), config); err != nil {
		return err
	}
	if err := generateGoSum(projectDir); err != nil {
		return fmt.Errorf("failed to generate go.sum: %v", err)
	}
	if err := renderTemplate("common/gitignore.tmpl", filepath.Join(projectDir, ".gitignore"), config); err != nil {
		return err
	}
	if err := renderTemplate("common/readme.md.tmpl", filepath.Join(projectDir, "README.md"), config); err != nil {
		return err
	}

	// Initialize Git repository if requested
	if config.InitGit {
		if err := utils.InitGit(projectDir); err != nil {
			return err
		}
	}

	return nil
}

func scaffoldAPI(projectDir string, config *prompts.ProjectConfig) error {
	// Create API-specific structure
	dirs := []string{
		filepath.Join(projectDir, "cmd", config.ProjectName),
		filepath.Join(projectDir, "internal", "api"),
		filepath.Join(projectDir, "internal", "handlers"),
		filepath.Join(projectDir, "internal", "models"),
		filepath.Join(projectDir, "internal", "config"),
	}
	if config.Database != constants.DatabaseNone {
		dirs = append(dirs, filepath.Join(projectDir, "internal", "db"))
	}
	for _, dir := range dirs {
		if err := utils.CreateDir(dir); err != nil {
			return err
		}
	}

	// Write API-specific files
	if err := renderTemplate("api/main.go.tmpl", filepath.Join(projectDir, "cmd", config.ProjectName, "main.go"), config); err != nil {
		return err
	}
	if err := renderTemplate("api/health.go.tmpl", filepath.Join(projectDir, "internal", "handlers", "health.go"), config); err != nil {
		return err
	}
	if err := renderTemplate("api/config.go.tmpl", filepath.Join(projectDir, "internal", "config", "config.go"), config); err != nil {
		return err
	}
	if err := renderTemplate("api/env.tmpl", filepath.Join(projectDir, ".env"), config); err != nil {
		return err
	}
	if config.Database != constants.DatabaseNone {
		if err := renderTemplate("api/db.go.tmpl", filepath.Join(projectDir, "internal", "db", "db.go"), config); err != nil {
			return err
		}
	}

	// Write Dockerfile if database is selected
	if config.Database != constants.DatabaseNone {
		dockerfileContent := fmt.Sprintf(`FROM golang:1.23
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o %s ./cmd/%s
EXPOSE 8080
CMD ["./%s"]
`, config.ProjectName, config.ProjectName, config.ProjectName)
		if err := utils.WriteFile(filepath.Join(projectDir, "Dockerfile"), dockerfileContent); err != nil {
			return err
		}

		// Write docker-compose.yml
		if err := renderTemplate("api/docker-compose.yml.tmpl", filepath.Join(projectDir, "docker-compose.yml"), config); err != nil {
			return err
		}
	}

	return nil
}

func scaffoldCLI(projectDir string, config *prompts.ProjectConfig) error {
	// Create CLI-specific structure
	if err := utils.CreateDir(filepath.Join(projectDir, "cmd", config.ProjectName)); err != nil {
		return err
	}

	// Write main.go
	if err := renderTemplate("cli/main.go.tmpl", filepath.Join(projectDir, "cmd", config.ProjectName, "main.go"), config); err != nil {
		return err
	}

	return nil
}

func scaffoldLibrary(projectDir string, config *prompts.ProjectConfig) error {
	// Write main library file
	if err := renderTemplate("lib/lib.go.tmpl", filepath.Join(projectDir, fmt.Sprintf("%s.go", config.ProjectName)), config); err != nil {
		return err
	}

	return nil
}

func renderTemplate(templateName, outputPath string, config *prompts.ProjectConfig) error {
	tmpl, err := template.New(filepath.Base(templateName)).Funcs(template.FuncMap{
		"replace": func(s string) string {
			return strings.ReplaceAll(s, "-", "_")
		},
	}).ParseFiles(filepath.Join("internal", "templates", templateName))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %v", templateName, err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", outputPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, config); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", templateName, err)
	}
	return nil
}

func generateGoSum(projectDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run 'go mod tidy' in %s: %v\nOutput: %s", projectDir, err, string(output))
	}
	return nil
}
