package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var root string
	flag.StringVar(&root, "dir", ".", "项目根目录")
	flag.Parse()

	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if err := updateFile(file); err != nil {
			panic(fmt.Sprintf("处理文件 %s 失败: %v", file, err))
		}
	}
}

func updateFile(path string) error {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	changed := false

	ast.Inspect(fileAst, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		desired := buildComment(fd)
		if fd.Doc == nil || len(fd.Doc.List) == 0 {
			fd.Doc = &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Slash: fd.Pos(),
						Text:  desired,
					},
				},
			}
			fileAst.Comments = append(fileAst.Comments, fd.Doc)
			changed = true
			return true
		}

		text := strings.TrimSpace(fd.Doc.Text())
		if text != strings.TrimPrefix(desired, "// ") || len(fd.Doc.List) != 1 || fd.Doc.List[0].Text != desired {
			fd.Doc.List = []*ast.Comment{
				{
					Slash: fd.Doc.List[0].Slash,
					Text:  desired,
				},
			}
			changed = true
		}

		return true
	})

	if !changed {
		return nil
	}

	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	if err := cfg.Fprint(&buf, fset, fileAst); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		return err
	}

	fmt.Println("updated", path)
	return nil
}

func buildComment(fd *ast.FuncDecl) string {
	name := fd.Name.Name
	entity := "函数"
	if fd.Recv != nil {
		entity = "方法"
	}
	return fmt.Sprintf("// %s %s用于处理%s相关逻辑。", name, entity, name)
}
