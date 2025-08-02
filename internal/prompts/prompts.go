package prompts

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/shubhamoys/forgego/internal/pkg/constants"
	"github.com/shubhamoys/forgego/internal/pkg/errors"
)

type ProjectConfig struct {
	ProjectType string
	ProjectName string
	PackageName string
	InitGit     bool
	Database    string
	GoVersion   string
}

func PromptNewProject() (*ProjectConfig, error) {
	config := &ProjectConfig{}

	// Custom select template for better UX
	selectTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ . | cyan }}",
		Inactive: "  {{ . | white }}",
		Selected: "Selected: {{ . | green }}",
	}

	// Prompt for project type
	projectTypePrompt := promptui.Select{
		Label: "Select project type",
		Items: []string{
			fmt.Sprintf("%s - %s", constants.ProjectTypeAPI, constants.ProjectTypeDescriptions[constants.ProjectTypeAPI]),
			fmt.Sprintf("%s - %s", constants.ProjectTypeCLI, constants.ProjectTypeDescriptions[constants.ProjectTypeCLI]),
			fmt.Sprintf("%s - %s", constants.ProjectTypeLibrary, constants.ProjectTypeDescriptions[constants.ProjectTypeLibrary]),
		},
		Templates: selectTemplate,
	}
	_, projectType, err := projectTypePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrPromptFailed, err)
	}
	config.ProjectType = strings.Split(projectType, " - ")[0]

	// Prompt for project name
	projectNamePrompt := promptui.Prompt{
		Label:    "Project name",
		Default:  "my-go-project",
		Validate: validateProjectName,
		Stdout:   os.Stdout,
		Stdin:    os.Stdin,
	}
	config.ProjectName, err = projectNamePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrPromptFailed, err)
	}

	// Prompt for package name
	packageNamePrompt := promptui.Prompt{
		Label:    "Go module package name",
		Default:  fmt.Sprintf("github.com/username/%s", config.ProjectName),
		Validate: validatePackageName,
		Stdout:   os.Stdout,
		Stdin:    os.Stdin,
	}
	config.PackageName, err = packageNamePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrPromptFailed, err)
	}

	config.GoVersion = strings.TrimPrefix(runtime.Version(), "go")

	// Prompt for Git initialization
	gitPrompt := promptui.Select{
		Label:     "Initialize Git repository?",
		Items:     []string{"Yes", "No"},
		Templates: selectTemplate,
	}
	_, gitChoice, err := gitPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrPromptFailed, err)
	}
	config.InitGit = gitChoice == "Yes"

	// Prompt for database (API projects only)
	if config.ProjectType == constants.ProjectTypeAPI {
		dbPrompt := promptui.Select{
			Label: "Select database",
			Items: []string{
				fmt.Sprintf("%s - %s", constants.DatabaseMongoDB, constants.DatabaseDescriptions[constants.DatabaseMongoDB]),
				fmt.Sprintf("%s - %s", constants.DatabasePostgreSQL, constants.DatabaseDescriptions[constants.DatabasePostgreSQL]),
				fmt.Sprintf("%s - %s", constants.DatabaseNone, constants.DatabaseDescriptions[constants.DatabaseNone]),
			},
			Templates: selectTemplate,
		}
		_, dbChoice, err := dbPrompt.Run()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrPromptFailed, err)
		}
		config.Database = strings.Split(dbChoice, " - ")[0]
	} else {
		config.Database = constants.DatabaseNone
	}

	return config, nil
}

func validateProjectName(name string) error {
	if len(name) == 0 {
		return errors.ErrInvalidProjectName
	}
	if strings.ContainsAny(name, " /\\") {
		return errors.ErrInvalidProjectName
	}
	return nil
}

func validatePackageName(name string) error {
	if len(name) == 0 {
		return errors.ErrInvalidPackageName
	}
	if !strings.HasPrefix(name, "github.com/") {
		return errors.ErrInvalidPackageName
	}
	return nil
}
