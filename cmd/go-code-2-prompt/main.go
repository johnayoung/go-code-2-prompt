package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/johnayoung/go-code-2-prompt/internal/fileutils"
	"github.com/johnayoung/go-code-2-prompt/internal/gitops"
	"github.com/johnayoung/go-code-2-prompt/internal/promptgen"
	"github.com/johnayoung/go-code-2-prompt/pkg/config"
	"github.com/johnayoung/go-code-2-prompt/pkg/tokenizer"
	"github.com/spf13/afero"
)

func main() {
	log.SetFlags(log.Lshortfile)

	startTime := time.Now()
	fmt.Println("Starting go-code-2-prompt...")

	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Print configuration
	fmt.Printf("Configuration:\n")
	fmt.Printf("  RootDir: %s\n", cfg.RootDir)
	fmt.Printf("  IncludePatterns: %v\n", cfg.IncludePatterns)
	fmt.Printf("  ExcludePatterns: %v\n", cfg.ExcludePatterns)
	fmt.Printf("  OutputFile: %s\n", cfg.OutputFile)
	fmt.Printf("  Tokenizer: %s\n", cfg.Tokenizer)

	fs := afero.NewOsFs()
	files, err := fileutils.TraverseDirectory(fs, cfg)
	if err != nil {
		log.Fatalf("Error traversing directory: %v", err)
	}

	fmt.Printf("Found %d files to include in the prompt:\n", len(files))
	for _, file := range files {
		relPath, _ := filepath.Rel(cfg.RootDir, file)
		fmt.Printf("  %s\n", relPath)
	}

	if cfg.IncludeGitDiff || cfg.IncludeGitLog {
		if !gitops.IsGitRepository(cfg.RootDir) {
			log.Printf("Warning: %s is not a git repository. Git information will not be included.", cfg.RootDir)
		}
	}

	prompt, err := promptgen.GeneratePrompt(fs, files, cfg)
	if err != nil {
		log.Fatalf("Error generating prompt: %v", err)
	}

	// Count tokens
	tokenizer, err := tokenizer.GetTokenizer(cfg.Tokenizer)
	if err != nil {
		log.Fatalf("Error getting tokenizer: %v", err)
	}

	tokenCount, err := tokenizer.CountTokens(prompt)
	if err != nil {
		log.Fatalf("Error counting tokens: %v", err)
	}

	if err := outputPrompt(prompt, cfg); err != nil {
		log.Fatalf("Error outputting prompt: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("Finished in %v\n", duration)
	fmt.Printf("Token count: %d\n", tokenCount)
}

func outputPrompt(prompt string, cfg *config.Config) error {
	if cfg.OutputFile != "" {
		return os.WriteFile(cfg.OutputFile, []byte(prompt), 0644)
	}
	_, err := fmt.Print(prompt)
	return err
}
