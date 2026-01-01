
# Dispense

Dispense is a simple utility that takes a directory of markdown files and transforms the files into a simple static website. Dispensed is similar to applications like [Hugo](https://gohugo.io/) or [Jekyll](https://jekyllrb.com/), but Dispense is purposefully much simplier than those applications. Its main use case, for me, is to transform a directory from [Obsidian](https://obsidian.md/) and create an uploadable static website.

Dispense works by using a markdown file's front matter to define a post's metadata (template, title, date, author, etc) and then generates html files from the file. The template files are stored in a directory called _hive_ (you can think of that directory as a theme if you like).

By default Dispense uses the Pico css styling (styles basic HTML tags) and adds KaTeX to support displaying maths:

- [Pico](https://picocss.com/)
- [KaTex](https://katex.org/)

## Usage

Make sure the markdown file has something like the following at the top of the file:

```markdown
---
author: rob
date: 2026-01-01
title: 🥾🐛气!
template: page
---
```

_template_ and _title_ are the most important. Then just run the application against a directory of markdown files (Dispense only processes the top level files, and the name of the markdown file becomes the name of the html file).

A simple example of using Dispensed:

```bash
mkdir public
dispense --base-input ~/Documents/Obsidian/98\ Output/example.com \
	--base-output ./public \
	--template-directory ./hive
cp -R ./hive/assets ./public/assets
```

Which would make the folder `public` available to upload and server.

## Command

```
Usage: dispense [options] [arguments]

OPTIONS
  --base-input/$DP_BASE_INPUT                  <string>  (default: ./test_data)
  --base-output/$DP_BASE_OUTPUT                <string>  (default: ./public)
  --template-extension/$DP_TEMPLATE_EXTENSION  <string>  (default: tpl)
  --template-directory/$DP_TEMPLATE_DIRECTORY  <string>  (default: ./hive)
  --template-listing/$DP_TEMPLATE_LISTING      <string>  (default: index)
  --template-post/$DP_TEMPLATE_POST            <string>  (default: post)
  --help/-h
  display this help message
```
