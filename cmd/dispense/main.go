package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ardanlabs/conf"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pkg/errors"
	models "github.com/robrohan/dispense/internals"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	ast.Leaf
	Data map[interface{}]interface{}
}

///////////////////////////////////////////////////////

func parseFrontMatter(data []byte) (ast.Node, []byte, int) {
	var frontMark = []byte("---\n")

	if !bytes.HasPrefix(data, frontMark) || bytes.HasPrefix(data, []byte("---\n\n")) {
		return nil, nil, 0
	}

	i := len(frontMark)
	end := bytes.Index(data[i:], []byte("---\n\n"))
	if end < 0 {
		return nil, data, 0
	}
	end = end + i

	matter := make(map[interface{}]interface{})
	err := yaml.Unmarshal(data[i:end], &matter)
	if err != nil {
		panic(err)
	}

	res := &FrontMatter{
		Data: matter,
	}
	return res, nil, end
}

func parserHook(data []byte) (ast.Node, []byte, int) {
	if node, d, n := parseFrontMatter(data); node != nil {
		return node, d, n
	}
	return nil, nil, 0
}

func frontMatterRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if _, ok := node.(*FrontMatter); ok {
		if entering {
			// fmt.Printf("%v\n", node.(*FrontMatter).Data)
			io.WriteString(w, "\n\n")
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

///////////////////////////////////////////////////////

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() && path.Ext(p) == ".md" {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}

func mdToHTML(md []byte) ([]byte, string, map[interface{}]interface{}) {
	template := "post"
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.Footnotes
	p := parser.NewWithExtensions(extensions)
	p.Opts.ParserHook = parserHook

	doc := p.Parse(md)

	var fm = make(map[interface{}]interface{})

	// if there is frontmatter is should be the first node
	root := doc.GetChildren()
	if root != nil {
		fm = root[0].(*FrontMatter).Data
		if fm["template"] != nil {
			template = fm["template"].(string)
		} else {
			template = "post"
		}
	}

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: frontMatterRenderHook,
	}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer), template, fm
}

// func renderTemplate(cfg *models.Config, log *log.Logger) {
// 	templateFile := cfg.Template.Directory + "/" + cfg.Template.Listing + "." + cfg.Template.Extension
// 	template, err := template.ParseFiles(templateFile)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	listings := map[string]string{
// 		"basic":               "basic.html",
// 		"Français":            "Français.html",
// 		"french_visual_cards": "french_visual_cards.html",
// 	}
// 	fo, err := os.Create(cfg.Base.Output + "/" + cfg.Template.Listing + ".html")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer fo.Close()
// 	template.Execute(fo, listings)
// }

func renderAllMarkdown(cfg *models.Config, log *log.Logger) {
	root := cfg.Base.Input
	files, err := FilePathWalkDir(root)
	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range files {
		fileName := path.Base(e)
		fileExt := path.Ext(e)
		fileName = strings.Replace(fileName, fileExt, "", 1)

		b, err := os.ReadFile(e)
		if err != nil {
			log.Fatalln(err)
		}

		htmlBytes, templateTitle, fm := mdToHTML(b)

		// log.Printf("%v", fm)
		templateFile := cfg.Template.Directory + "/" + templateTitle + "." + cfg.Template.Extension
		log.Printf("using template file %s\n", templateFile)
		template, err := template.ParseFiles(templateFile)
		if err != nil {
			log.Fatalln(err)
		}

		// Add in the HTML rendered text
		fm["postData"] = string(htmlBytes)

		output := cfg.Base.Output + "/" + fileName + ".html"
		log.Printf("output file %s\n", output)
		fo, err := os.Create(output)
		if err != nil {
			log.Fatalln(err)
		}
		defer fo.Close()

		template.Execute(fo, fm)
	}
}

func run() error {
	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "DP : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	cfg := models.Config{}

	if err := conf.Parse(os.Args[1:], "DP", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DP", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	renderAllMarkdown(&cfg, log)
	// renderTemplate(&cfg, log)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}
