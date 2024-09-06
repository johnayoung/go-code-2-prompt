package fileutils

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobwas/glob"
	"github.com/johnayoung/go-code-2-prompt/pkg/config"
	"github.com/karrick/godirwalk"
	"github.com/spf13/afero"
)

func TraverseDirectory(fs afero.Fs, cfg *config.Config) ([]string, error) {
	var files []string

	err := godirwalk.Walk(cfg.RootDir, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.IsDir() {
				return nil
			}
			if ShouldIncludeFile(path, cfg.IncludePatterns, cfg.ExcludePatterns) {
				files = append(files, path)
			}
			return nil
		},
		Unsorted: true,
	})

	return files, err
}

func ShouldIncludeFile(path string, includePatterns, excludePatterns []string) bool {
	// Check exclude patterns first
	for _, pattern := range excludePatterns {
		g := glob.MustCompile(pattern)
		if g.Match(path) {
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
		if g.Match(path) {
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
	textExtensions := []string{".txt", ".md", ".go", ".py", ".js", ".html", ".css", ".json", ".yaml", ".yml", ".toml"}
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
