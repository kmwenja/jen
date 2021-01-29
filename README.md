Jen
===

Generate a html page using Golang templates.
Jen also allows for JSON data to be applied to the template during generation.

Usage
-----

Help: `jen --help`.

Generate HTML with JSON data:
```html
# template.html
Hello {{.name}}!
```

```bash
echo '{"name": "World!"}' | jen g template.html -
```

Using markdown:
```markdown
# content.md
---
title: Hello World
---
# Hello World

Some kind words.
```

```html
# template.html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.yaml.title}}</title>
</head>
<body>
   {{.markdown}}
</body>
</html>
```

```bash
jen m content.md | jen g template.html -
```
