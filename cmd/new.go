package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spf13/cobra"
)

const (
	baseRepoURL     = "https://github.com/q1mi/gin-base-layout.git"
	advancedRepoURL = "https://github.com/q1mi/gin-advanced-layout.git"
)

// Project 结构体，用于存储项目信息
type Project struct {
	ProjectName string `survey:"name"`
	FolderName  string // 文件夹名称，例如：github.com/xxx/xx -> xx
	RepoURL     string // 模板仓库URL
}

var (
	advanced bool
	repoURL  string
)

var NewCmd = &cobra.Command{
	Use:     "new",
	Example: "iaa new project-name [--advanced] [--repo <url>]",
	Short:   "create a new project.",
	Long:    `create a new project with gin-base-layout or gin-advanced-layout.`,
	Run:     run,
}

func init() {
	NewCmd.Flags().BoolVar(&advanced, "advanced", false, "use advanced template (gin-advanced-layout)")
	NewCmd.Flags().StringVar(&repoURL, "repo", "", "specify custom template repository URL")
}

func NewProject(projectName, templateRepo string) *Project {
	return &Project{
		ProjectName: projectName,
		FolderName:  filepath.Base(filepath.Clean(projectName)), // get xx from github.com/xxx/xx
		RepoURL:     templateRepo,
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("need project name")
		return
	}

	// 确定使用的模板仓库
	var templateRepo string
	switch {
	case repoURL != "":
		// 自定义仓库优先级最高
		templateRepo = repoURL
	case advanced:
		// 使用进阶模板
		templateRepo = advancedRepoURL
	default:
		// 默认使用基础模板
		templateRepo = baseRepoURL
	}

	p := NewProject(args[0], templateRepo)
	fmt.Printf("🚀 Start to create project \u001B[36m%s\u001B[0m...\n", p.ProjectName)
	// clone repo
	yes, err := p.cloneRepo()
	if err != nil || !yes {
		return
	}

	// replace package name
	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	// go mod tidy
	err = p.modTidy()
	if err != nil || !yes {
		return
	}
	p.rmGit()
	fmt.Printf("🎉 Project \u001B[36m%s\u001B[0m created successfully!\n\n", p.ProjectName)
	fmt.Printf("Now run:\n\n")
	fmt.Printf("› \033[36mcd %s \033[0m\n", p.FolderName)
	fmt.Printf("› \033[36mgo run cmd/server/main.go\033[0m\n\n")
}

func (p *Project) cloneRepo() (bool, error) {
	// 1.检查目录是否已存在
	stat, _ := os.Stat(p.FolderName)
	if stat != nil {
		var overwrite = false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Folder %s already exists, do you want to overwrite it?", p.FolderName),
			Help:    "Remove the old project and create a new one.",
		}
		err := survey.AskOne(prompt, &overwrite)
		if err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
		err = os.RemoveAll(p.FolderName)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return false, err
		}
	}

	fmt.Println("git clone ", p.RepoURL)
	cmd := exec.Command("git", "clone", p.RepoURL, p.FolderName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", p.RepoURL, err)
		return false, err
	}
	return true, nil
}

func (p *Project) replacePackageName() error {
	moduleName := p.getModuleName()
	if len(moduleName) == 0 {
		return fmt.Errorf("get module name error")
	}
	err := p.replaceFiles(moduleName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ProjectName)
	cmd.Dir = p.FolderName
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod edit error: ", err)
		return err
	}
	return nil
}
func (p *Project) modTidy() error {
	fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.FolderName
	if err := cmd.Run(); err != nil {
		fmt.Println("go mod tidy error: ", err)
		return err
	}
	return nil
}
func (p *Project) rmGit() {
	os.RemoveAll(filepath.Join(p.FolderName, ".git"))
}

func (p *Project) replaceFiles(old string) error {
	err := filepath.Walk(p.FolderName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newData := bytes.ReplaceAll(data, []byte(old), []byte(p.ProjectName))
		if err := os.WriteFile(path, newData, 0644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("walk file do replace error: ", err)
		return err
	}
	return nil
}

// getModuleName 从 go.mod 中获取 module name
func (p *Project) getModuleName() string {
	modFile, err := os.Open(filepath.Join(p.FolderName, "go.mod"))
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}
