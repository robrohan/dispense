package main

import (
	"bufio"
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
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/pkg/errors"
	models "github.com/robrohan/dispense/internals"
)

// var mds = `# header
// Sample text.
// [link](http://example.com)
// `

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

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func renderTemplate(cfg *models.Config, log *log.Logger) {
	templateFile := cfg.Template.Directory + "/" + cfg.Template.Listing + "." + cfg.Template.Extension
	template, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatalln(err)
	}
	listings := map[string]string{
		"basic":               "basic.html",
		"Français":            "Français.html",
		"french_visual_cards": "french_visual_cards.html",
	}
	fo, err := os.Create(cfg.Base.Output + "/" + cfg.Template.Listing + ".html")
	if err != nil {
		log.Fatalln(err)
	}
	defer fo.Close()
	template.Execute(fo, listings)
}

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

		log.Printf("%s\n", fileName)

		// make a read buffer
		r := bytes.NewReader(mdToHTML(b))

		// open output file
		fo, err := os.Create(cfg.Base.Output + "/" + fileName + ".html")
		if err != nil {
			log.Fatalln(err)
		}
		defer fo.Close()

		// make a write buffer
		w := bufio.NewWriter(fo)

		// make a buffer to keep chunks that are read
		buf := make([]byte, 1024)
		for {
			// read a chunk
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatalln(err)
			}
			if n == 0 {
				break
			}

			// write a chunk
			if _, err := w.Write(buf[:n]); err != nil {
				log.Fatalln(err)
			}
		}

		if err = w.Flush(); err != nil {
			log.Fatalln(err)
		}
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
	renderTemplate(&cfg, log)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}
