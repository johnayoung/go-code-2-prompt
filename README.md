# go-code-2-prompt

go-code-2-prompt is a powerful command-line tool designed to generate prompts for Large Language Models (LLMs) from codebases. It offers a range of features to help developers create comprehensive and context-rich prompts for various LLM-based tasks such as code analysis, documentation generation, and more.

## Features

- **Directory Traversal**: Efficiently traverses specified directories, building a tree structure of the codebase.
- **File Filtering**: Allows users to specify include and exclude patterns using glob syntax.
- **Code Extraction**: Reads the content of files that match the include criteria, supporting various file types and programming languages.
- **Source Tree Generation**: Creates a visual representation of the codebase structure.
- **Git Integration**: 
  - Includes git diff output for staged changes.
  - Supports generating diffs between two specified branches.
  - Can retrieve git log information between two branches.
- **Template System**: Uses Go's text/template package for customizable output.
- **Output Formatting**: Generates a formatted output with a directory tree, file contents, and optional Git information.
- **Flexible Output**: Supports writing to a file or standard output.
- **Token Counting**: Counts the number of tokens in the generated prompt using various tokenization schemes (cl100k_base, p50k_base, r50k_base) and displays the count in the terminal.
- **Execution Time**: Displays the total execution time of the program in the terminal.
- **Automatic .git Exclusion**: The tool automatically excludes the .git folder and its contents from the generated prompt and source tree.

## Behavior

- The tool automatically excludes the .git folder and its contents from both the file traversal and the source tree generation. This helps to reduce the output size and remove unnecessary information from the generated prompt.
- If you need to include git-related information, use the `-git-diff` and `-git-log` flags. These will include relevant git information without including the entire .git folder contents.

## Installation

To install go-code-2-prompt, make sure you have Go 1.22.3 or later installed, then run:

```
go install github.com/johnayoung/go-code-2-prompt@latest
```

Alternatively, you can clone the repository and build it manually:

```
git clone https://github.com/johnayoung/go-code-2-prompt.git
cd go-code-2-prompt
go build ./cmd/go-code-2-prompt
```

## Usage

Here's the basic usage of go-code-2-prompt:

```
go-code-2-prompt [flags]
```

### Flags

- `-dir string`: Root directory to traverse (default ".")
- `-output string`: Output file (default: stdout)
- `-tokenizer string`: Tokenizer to use (default "cl100k_base", options: "cl100k_base", "p50k_base", "r50k_base")
- `-template string`: Custom template file
- `-include string`: Include patterns (comma-separated)
- `-exclude string`: Exclude patterns (comma-separated)
- `-git-diff`: Include git diff of staged changes
- `-branch1 string`: First branch for git diff/log (default: current branch)
- `-branch2 string`: Second branch for git diff/log
- `-git-log`: Include git log between branches
- `-show-high-token-folders`: Show folders with high token counts
- `-high-token-folder-count int`: Number of high token folders to show (default 5)

### Examples

1. Generate a prompt for the current directory:
   ```
   go-code-2-prompt
   ```

2. Generate a prompt for a specific directory, including only Go files:
   ```
   go-code-2-prompt -dir /path/to/project -include "*.go"
   ```

3. Generate a prompt including files from multiple specific directories:
   ```
   go-code-2-prompt -dir /path/to/project -include "path/dir1/*,path/dir2/*"
   ```

4. Generate a prompt with Git diff and log information:
   ```
   go-code-2-prompt -git-diff -git-log -branch1 main -branch2 feature-branch
   ```

5. Generate a prompt and save it to a file:
   ```
   go-code-2-prompt -output prompt.txt
   ```

6. Generate a prompt for the current directory and show high token folders:
   ```
   go-code-2-prompt -show-high-token-folders
   ```

7. Generate a prompt for a specific directory, including only Go files, and show top 10 high token folders:
   ```
   go-code-2-prompt -dir /path/to/project -include "*.go" -show-high-token-folders -high-token-folder-count 10
   ```

Note: When specifying multiple include or exclude patterns, separate them with commas. Do not use multiple `-include` or `-exclude` flags.

## Output

The generated prompt includes:

1. A source tree representation of the codebase structure.
2. The contents of each included file.
3. Git diff information (if requested).
4. Git log information (if requested).

The token count and execution time are displayed in the terminal but not included in the generated prompt.

The generated prompt includes:

1. A source tree representation of the codebase structure.
2. The contents of each included file.
3. Git diff information (if requested).
4. Git log information (if requested).

## Customization

You can customize the output by providing your own template file using the `-template` flag. The template uses Go's text/template syntax and has access to the following data:

- `.SourceTree`: The generated source tree string.
- `.Files`: A slice of FileInfo structures containing RelativePath and Content for each file.
- `.GitDiff`: The git diff string (if requested).
- `.GitLog`: The git log string (if requested).
- `.GitBranch`: The current git branch.
- `.Config`: The configuration object containing all command-line options.

## Contributing

Contributions to go-code-2-prompt are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- This project uses several open-source libraries, including `github.com/karrick/godirwalk` for efficient directory traversal and `github.com/spf13/afero` for a flexible filesystem interface.

## Future Enhancements

- ~~Token counting for various LLM models.~~ (Implemented)
- ~~Progress indication for long-running operations.~~ (Implemented)
- Clipboard support for easy copying of generated prompts.
- JSON output support for integration with other tools.

For any questions or issues, please open an issue on the GitHub repository.