# File Manager
[![Go](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org/)

**File Manager** is a CLI tool for analyzing and managing files and directories. It helps find duplicate files, analyze disk space usage, and perform other useful operations.
---
## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
    - [Find Duplicate Files](#find-duplicate-files)
    - [Analyze Disk Space Usage](#analyze-disk-space-usage)
    - [Search Files by Pattern](#search-files-by-pattern)
- [Flags](#flags)
- [Examples](#examples)
---
## Installation
### Method 1: Building from Source
1. Ensure that Go (version 1.20+) is installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/SHCDevelops/file-manager.git
   ```
3. Navigate to the project directory:
   ```bash
   cd file-manager
   ```
4. Build the project:
   ```bash
   go build -o file-manager
   ```
5. Add the executable to your PATH:
   ```bash
   mv file-manager /usr/local/bin/
   ```
Now you can use the `file-manager` command globally.
---
### Method 2: Using `go install`
The project is available via GitHub, and you can install it directly:
```bash
go install github.com/SHCDevelops/file-manager@latest
```
---
## Usage
### Find Duplicate Files
This command scans the specified directory and finds files with identical content.
```bash
file-manager find-duplicates [directory] [flags]
```
#### Example:
```bash
file-manager find-duplicates /path/to/directory --ignore ".git,temp"
```
---
### Analyze Disk Space Usage
This command shows the largest files in the specified directory.
```bash
file-manager analyze-space [directory] [flags]
```
#### Example:
```bash
file-manager analyze-space /path/to/directory --top 10 --ignore "*.tmp"
```
---
### Search Files by Pattern
This command searches for files matching the given pattern.
```bash
file-manager search [pattern] [directory]
```
#### Example:
```bash
file-manager search "*.txt" /path/to/directory
```
---
### Code Statistics Analysis

This command shows detailed code statistics for supported programming languages.

```bash
file-manager code-stats [directory] [flags]
```

**Supported Languages:**
- Go (.go)
- HTML (.html, .htm)
- CSS (.css)
- JavaScript (.js)
- TypeScript (.ts, .tsx)

**Displayed Metrics:**
- Total lines of code
- Comment lines count
- Pure code lines (total - comments)
- Percentage ratio

#### Example:
```bash
file-manager code-stats ./myproject --ignore "vendor,node_modules"
```
---
## Flags

| Command           | Flag                | Description                                                  |
|-------------------|---------------------|--------------------------------------------------------------|
| `find-duplicates` | `--ignore`          | List of directories or patterns to ignore (comma-separated). |
| `analyze-space`   | `--top`             | Number of files to display (default: 10).                    |
| `analyze-space`   | `--ignore`          | List of directories or patterns to ignore (comma-separated). |
| `code-stats`      | `--ignore`          | List of directories or patterns to ignore (comma-separated). |
| `code-stats`      | `--ignore-language` | List of languages to ignore (comma-separated).               |
---
## Examples
### 1. Find duplicate files, ignoring `.git` and `temp` directories:
```bash
file-manager find-duplicates /path/to/directory --ignore ".git,temp"
```
### 2. Show the top-5 largest files, ignoring temporary files:
```bash
file-manager analyze-space /path/to/directory --top 5 --ignore "*.tmp"
```
### 3. Find all text files in a directory:
```bash
file-manager search "*.txt" /path/to/directory
```
### 4. Analyze project code statistics:
```bash
file-manager code-stats ./src --ignore "tests,dist"
```

**Sample output:**
```
Code Statistics:

Go:
  Total lines: 1520
  Comments:    320 (21.1%)
  Code lines:  1200 (78.9%)

JavaScript:
  Total lines: 890
  Comments:    178 (20.0%)
  Code lines:  712 (80.0%)
```
---
## Requirements
- **Go**: Version 1.20 or higher.
- **Operating System**: Linux, macOS, Windows.
---
## Author
Created with ❤️ by [SHCDevelops](https://github.com/SHCDevelops)
If you have any questions or suggestions, feel free to open an issue or send a pull request!
---
### Notes
- This tool is intended for personal use and can be adapted to meet your needs.
- For working with large directories, it is recommended to use the `--ignore` flag to speed up the process.
---