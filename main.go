package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"
)

const indexHtmlTemplate = `
<!DOCTYPE html>
<html lang=en>
	<head>
		<title>Haikud</title>
	</head>
	<body>
    {{range .}}
    <article>
      <p>
        {{index .Lines 0}}<br>
        {{index .Lines 1}}<br>
        {{index .Lines 2}}
      </p>
      <time datetime="{{.CreatedAt}}">{{formatDate .CreatedAt}}</time>
    </article>
    {{end}}
	</body>
</html>
`

func main() {
	if err := generateIndexHtml(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Haiku struct {
	Lines     []string  `json:"lines"`
	CreatedAt time.Time `json:"created_at"`
}

func generateIndexHtml() error {
	funcs := template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("02.01.2006")
		},
	}

	tmpl, err := template.New("").Funcs(funcs).Parse(indexHtmlTemplate)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile("haikus.json")
	if err != nil {
		return err
	}

	var haikus []Haiku
	if err := json.Unmarshal(b, &haikus); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, haikus); err != nil {
		return err
	}

	if err := ioutil.WriteFile("index.html", buf.Bytes(), 0664); err != nil {
		return err
	}

	return nil
}
