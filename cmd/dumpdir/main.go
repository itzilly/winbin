package main

import (
    "flag"
    "fmt"
    "io/fs"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"

    ignore "github.com/sabhiram/go-gitignore"
)

func main() {
    output := flag.String("o", "dir.txt", "Output file name (default: dir.txt)")
    consoleOnly := flag.Bool("c", false, "Print to console only, don't save to file")
    noGitIgnore := flag.Bool("no-gitignore", false, "Disable .gitignore rules and include .git contents")
    flag.Parse()

    cwd, err := os.Getwd()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
        os.Exit(1)
    }

    var ign *ignore.GitIgnore
    gitignorePath := filepath.Join(cwd, ".gitignore")
    if !_fileExists(gitignorePath) {
        ign = ignore.CompileIgnoreLines()
    } else {
        ign, _ = ignore.CompileIgnoreFile(gitignorePath)
    }

    var dump strings.Builder

    err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        relPath, _ := filepath.Rel(cwd, path)
        if relPath == "." {
            return nil
        }

        if filepath.Base(path) == ".gitignore" {
            dump.WriteString(fmt.Sprintf("=== %s ===\n", relPath))
            content, err := ioutil.ReadFile(path)
            if err != nil {
                dump.WriteString(fmt.Sprintf("Error reading file: %v\n\n", err))
            } else {
                dump.Write(content)
                dump.WriteString("\n\n")
            }
            return nil
        }

        if relPath == ".git" {
            dump.WriteString("=== .git ===\n<git folder present - contents skipped>\n\n")
            if *noGitIgnore {
                return nil
            }
            return fs.SkipDir
        }

        if !_shouldInclude(relPath, d, ign, *noGitIgnore) {
            if d.IsDir() {
                return fs.SkipDir
            }
            return nil
        }

        dump.WriteString(fmt.Sprintf("=== %s ===\n", relPath))
        if !d.IsDir() {
            content, err := ioutil.ReadFile(path)
            if err != nil {
                dump.WriteString(fmt.Sprintf("Error reading file: %v\n\n", err))
            } else {
                dump.Write(content)
                dump.WriteString("\n\n")
            }
        }

        return nil
    })

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
        os.Exit(1)
    }

    if *consoleOnly {
        fmt.Print(dump.String())
    } else {
        err := ioutil.WriteFile(*output, []byte(dump.String()), 0644)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Directory dumped to %s\n", *output)
    }
}

func _fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func _shouldInclude(relPath string, d fs.DirEntry, ign *ignore.GitIgnore, noIgnore bool) bool {
    if noIgnore {
        return true
    }
    return !ign.MatchesPath(relPath)
}
