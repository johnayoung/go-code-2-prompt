package fileutils

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobwas/glob"
	"github.com/johnayoung/go-code-2-prompt/pkg/config"
	"github.com/johnayoung/go-code-2-prompt/pkg/tokenizer"
	"github.com/karrick/godirwalk"
	"github.com/spf13/afero"
)

// Add this to internal/fileutils/fileutils.go

func TraverseDirectory(fs afero.Fs, cfg *config.Config, tokenizer tokenizer.Tokenizer) ([]string, map[string]*config.FolderInfo, error) {
	files := []string{}
	folderInfo := make(map[string]*config.FolderInfo)

	err := godirwalk.Walk(cfg.RootDir, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				if de.Name() == ".git" {
					return godirwalk.SkipThis
				}
				return nil
			}
			relPath, err := filepath.Rel(cfg.RootDir, path)
			if err != nil {
				return err
			}
			if ShouldIncludeFile(relPath, cfg.IncludePatterns, cfg.ExcludePatterns) {
				files = append(files, path)

				if cfg.ShowHighTokenFolders {
					content, err := ReadFileContent(fs, path)
					if err != nil {
						return err
					}
					tokenCount, err := tokenizer.CountTokens(content)
					if err != nil {
						return err
					}

					folder := filepath.Dir(relPath)
					if _, ok := folderInfo[folder]; !ok {
						folderInfo[folder] = &config.FolderInfo{Path: folder}
					}
					folderInfo[folder].TokenCount += tokenCount
					folderInfo[folder].FileCount++
				}
			}
			return nil
		},
		Unsorted: true,
	})

	return files, folderInfo, err
}

func ShouldIncludeFile(relPath string, includePatterns, excludePatterns []string) bool {
	// Always exclude files in the .git directory
	if strings.Contains(relPath, ".git"+string(filepath.Separator)) {
		return false
	}

	// Check exclude patterns first
	for _, pattern := range excludePatterns {
		g := glob.MustCompile(pattern)
		if g.Match(relPath) {
			return false
		}
	}

	// If no include patterns are specified, include all files
	if len(includePatterns) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range includePatterns {
		g := glob.MustCompile(pattern)
		if g.Match(relPath) {
			return true
		}
	}

	return false
}

func ReadFileContent(fs afero.Fs, path string) (string, error) {
	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetRelativePath(rootDir, filePath string) string {
	relPath, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		return filePath
	}
	return relPath
}

func IsTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	textExtensions := []string{
		".txt",
		".md",
		".go",
		".py",
		".js",
		".ts",
		".jsx",
		".tsx",
		".html",
		".css",
		".json",
		".yaml",
		".yml",
		".toml",
	}
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

func SortedDirEntries(entries []fs.DirEntry) []fs.DirEntry {
	sort.Slice(entries, func(i, j int) bool {
		// Directories come first
		if entries[i].IsDir() && !entries[j].IsDir() {
			return true
		}
		if !entries[i].IsDir() && entries[j].IsDir() {
			return false
		}
		// Then sort alphabetically
		return entries[i].Name() < entries[j].Name()
	})
	return entries
}

func GetHighTokenFolders(folderInfo map[string]*config.FolderInfo, count int) []config.FolderInfo {
	folders := make([]config.FolderInfo, 0, len(folderInfo))
	for _, info := range folderInfo {
		folders = append(folders, *info)
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].TokenCount > folders[j].TokenCount
	})

	if count > len(folders) {
		count = len(folders)
	}

	return folders[:count]
}
