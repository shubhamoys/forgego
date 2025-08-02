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

	// Generate folder structure for API project
	if err := scaffoldAPI(projectDir, config); err != nil {
		return err
	}

	// Generate go.mod and go.sum
	if err := generateGoSum(projectDir, config); err != nil {
		return fmt.Errorf("failed to generate go.mod and go.sum: %v", err)
	}

	// Write common files
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

	// Write Docker files if requested
	if config.IncludeDocker {
		dockerfileContent := fmt.Sprintf(`FROM golang:%s
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o %s ./cmd/%s
EXPOSE 8080
CMD ["./%s"]
`, config.GoVersion, config.ProjectName, config.ProjectName, config.ProjectName)
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

func generateGoSum(projectDir string, config *prompts.ProjectConfig) error {
	// Initialize go.mod
	cmd := exec.Command("go", "mod", "init", config.PackageName)
	cmd.Dir = projectDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run 'go mod init %s' in %s: %v\nOutput: %s", config.PackageName, projectDir, err, string(output))
	}

	// Add API dependencies
	cmds := [][]string{
		{"go", "get", "github.com/gorilla/mux"},
		{"go", "get", "github.com/joho/godotenv"},
	}
	if config.Database == constants.DatabaseMongoDB {
		cmds = append(cmds, []string{"go", "get", "go.mongodb.org/mongo-driver/mongo"})
	} else if config.Database == constants.DatabasePostgreSQL {
		cmds = append(cmds, []string{"go", "get", "github.com/lib/pq"})
	}

	// Run go get commands
	for _, cmdArgs := range cmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = projectDir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to run '%s' in %s: %v\nOutput: %s", strings.Join(cmdArgs, " "), projectDir, err, string(output))
		}
	}

	// Run go mod tidy to generate go.sum
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy' in %s: %v\nOutput: %s", projectDir, err, string(output))
	}

	return nil
}
