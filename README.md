Jen
===

Generate a html page from a markdown file.
Can use a template to further customize the output.

Usage
-----

`jen file.md > file.html`

`jen --template-data file.yaml|file.json|file.toml file.md > file.html`

```
# file.md
---
template: jen.template
title: Hello World
---

This is a markdown file.
```

```
# jen.template
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{ .front_matter.title }}</title>
</head>
<body>
    <h1>{{ .front_matter.title }}</h1>

    {{ .content }}
</body>
</html>
```

Features
--------
- Accept multiple `template` (solves the problem of template composition)
- Accept multiple `template-data` (solves the problem of data composition)
- Support frontline yaml matter
- Support templating within the markdown

Notes
-----

ls to json: `tree -n -D --noreport -L 1 -J`
Blog Roll: `for i in blog/*.md; do jen extract-front-matter $(cat $i) | jq {.title,.date}
