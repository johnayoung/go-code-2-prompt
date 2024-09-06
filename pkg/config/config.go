package config

import (
	"flag"
	"strings"
)

type Config struct {
	RootDir         string
	IncludePatterns []string
	ExcludePatterns []string
	OutputFile      string
	Tokenizer       string
	Template        string
	IncludeGitDiff  bool
	GitBranch1      string
	GitBranch2      string
	IncludeGitLog   bool
}

func ParseFlags() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.RootDir, "dir", ".", "Root directory to traverse")
	flag.StringVar(&config.OutputFile, "output", "", "Output file (default: stdout)")
	flag.StringVar(&config.Tokenizer, "tokenizer", "cl100k_base", "Tokenizer to use (options: cl100k_base, p50k_base, r50k_base)")
	flag.StringVar(&config.Template, "template", "", "Custom template file")
	flag.BoolVar(&config.IncludeGitDiff, "git-diff", false, "Include git diff of staged changes")
	flag.StringVar(&config.GitBranch1, "branch1", "", "First branch for git diff/log (default: current branch)")
	flag.StringVar(&config.GitBranch2, "branch2", "", "Second branch for git diff/log")
	flag.BoolVar(&config.IncludeGitLog, "git-log", false, "Include git log between branches")

	includes := flag.String("include", "", "Include patterns (comma-separated)")
	excludes := flag.String("exclude", "", "Exclude patterns (comma-separated)")

	flag.Parse()

	if *includes != "" {
		config.IncludePatterns = splitPatterns(*includes)
	}
	if *excludes != "" {
		config.ExcludePatterns = splitPatterns(*excludes)
	}

	return config, nil
}

func splitPatterns(patterns string) []string {
	return strings.Split(patterns, ",")
}
