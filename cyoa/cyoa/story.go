package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Story map[string]Chapter

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

func ParseJSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	story := Story{}

	err := d.Decode(&story)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return story, nil
}

type handler struct {
	s Story
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tmpl *template.Template

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8" />
        <title>Choose your own adventure</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
        <ul>
        {{range .Options}}
            <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
        {{end}}
        </ul>
    </body>
</html>`

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
}
