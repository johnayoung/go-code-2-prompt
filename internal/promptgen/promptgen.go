package promptgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/johnayoung/go-code-2-prompt/internal/fileutils"
	"github.com/johnayoung/go-code-2-prompt/internal/gitops"
	"github.com/johnayoung/go-code-2-prompt/pkg/config"
	"github.com/spf13/afero"
)

type FileInfo struct {
	RelativePath string
	Content      string
}

type PromptData struct {
	Files      []FileInfo
	Config     *config.Config
	GitDiff    string
	GitLog     string
	GitBranch  string
	SourceTree string
}

func GeneratePrompt(fs afero.Fs, files []string, cfg *config.Config) (string, error) {
	var fileInfos []FileInfo

	for _, file := range files {
		if !fileutils.IsTextFile(file) {
			continue
		}

		content, err := fileutils.ReadFileContent(fs, file)
		if err != nil {
			return "", fmt.Errorf("error reading file %s: %v", file, err)
		}

		relPath := fileutils.GetRelativePath(cfg.RootDir, file)
		fileInfos = append(fileInfos, FileInfo{
			RelativePath: relPath,
			Content:      content,
		})
	}

	sourceTree, err := generateSourceTree(cfg.RootDir)
	if err != nil {
		return "", fmt.Errorf("error generating source tree: %v", err)
	}

	promptData := PromptData{
		Files:      fileInfos,
		Config:     cfg,
		SourceTree: sourceTree,
	}

	if gitops.IsGitRepository(cfg.RootDir) {
		if cfg.IncludeGitDiff {
			diff, err := gitops.GetStagedDiff(cfg.RootDir)
			if err != nil {
				return "", fmt.Errorf("error getting git diff: %v", err)
			}
			promptData.GitDiff = diff
		}

		if cfg.IncludeGitLog {
			branch1 := cfg.GitBranch1
			if branch1 == "" {
				var err error
				branch1, err = gitops.GetCurrentBranch(cfg.RootDir)
				if err != nil {
					return "", fmt.Errorf("error getting current branch: %v", err)
				}
			}
			promptData.GitBranch = branch1

			if cfg.GitBranch2 != "" {
				log, err := gitops.GetGitLog(cfg.RootDir, branch1, cfg.GitBranch2)
				if err != nil {
					return "", fmt.Errorf("error getting git log: %v", err)
				}
				promptData.GitLog = log
			}
		}
	}

	return renderTemplate(promptData, cfg)
}

func generateSourceTree(rootDir string) (string, error) {
	var buffer bytes.Buffer
	buffer.WriteString("Source Tree:\n\n```\n")
	buffer.WriteString(filepath.Base(rootDir) + "\n")

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == rootDir {
			return nil
		}

		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		depth := strings.Count(rel, string(os.PathSeparator))
		if info.IsDir() {
			buffer.WriteString(strings.Repeat("│   ", depth) + "├── " + info.Name() + "\n")
		} else {
			buffer.WriteString(strings.Repeat("│   ", depth) + "└── " + info.Name() + "\n")
		}

		return nil
	})

	buffer.WriteString("```\n")

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func renderTemplate(data PromptData, cfg *config.Config) (string, error) {
	var tmpl *template.Template
	var err error

	if cfg.Template != "" {
		tmpl, err = template.ParseFiles(cfg.Template)
	} else {
		tmpl, err = template.New("default").Parse(defaultTemplate)
	}

	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}

func OutputPrompt(prompt string, cfg *config.Config) error {
	if cfg.OutputFile != "" {
		return os.WriteFile(cfg.OutputFile, []byte(prompt), 0644)
	}
	fmt.Print(prompt)
	return nil
}

const defaultTemplate = `
{{.SourceTree}}

File Contents:
{{range .Files}}
--- {{.RelativePath}} ---
{{.Content}}

{{end}}

{{if .GitDiff}}
Git Diff (Staged Changes):
{{.GitDiff}}
{{end}}

{{if .GitLog}}
Git Log ({{.GitBranch}} to {{.Config.GitBranch2}}):
{{.GitLog}}
{{end}}
`
