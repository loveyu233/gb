package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var root string
	flag.StringVar(&root, "dir", ".", "项目根目录")
	flag.Parse()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "vendor" {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		return processFile(path)
	})
}

func processFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var out []string

	changed := false
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "//") && i+2 < len(lines) {
			if strings.TrimSpace(lines[i+1]) == "" && strings.TrimSpace(lines[i+2]) == trimmed {
				changed = true
				i++
				continue
			}
		}
		out = append(out, lines[i])
	}

	if !changed {
		return nil
	}

	return os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o644)
}
