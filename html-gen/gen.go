package gen

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"

	"github.com/karlssonerik/go-service-doc/core"
)

type Gen struct {
	api         string
	pages       core.Pages
	doc         string
	searchLink  string
	queryString string
	basepath    string
	faviconHref string
}

func New() *Gen {
	return &Gen{}
}

func (g *Gen) WithAPITitle(apiTitle string) *Gen {
	g.api = apiTitle
	return g
}

func (g *Gen) WithPages(pages core.Pages) *Gen {
	g.pages = pages
	return g
}

func (g *Gen) WithDocument(doc string) *Gen {
	g.doc = doc
	return g
}

func (g *Gen) WithSearchLink(searchLink string) *Gen {
	g.searchLink = searchLink
	return g
}

func (g *Gen) WithBasepath(basepath string) *Gen {
	g.basepath = basepath
	return g
}

func (g *Gen) WithFavicon(href string) *Gen {
	g.faviconHref = href
	return g
}

func (g *Gen) Build() (_ []byte, err error) {
	templateInfo := struct {
		API         string
		Pages       core.Pages
		Doc         string
		SearchLink  string
		QueryString string
		Basepath    string
		FaviconHref string
	}{
		API:         g.api,
		Pages:       g.pages,
		Doc:         g.doc,
		SearchLink:  g.searchLink,
		QueryString: g.queryString,
		Basepath:    g.basepath,
		FaviconHref: g.faviconHref,
	}

	generator, err := template.New("html_page").Parse(htmlPageTemplate)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse HTML page template")
		return
	}

	buffer := &bytes.Buffer{}
	if err = generator.Execute(buffer, templateInfo); err != nil {
		err = errors.Wrapf(err, "failed to execute generator")
		return
	}

	return buffer.Bytes(), nil
}

func (g *Gen) BuildSearchPageTemplate() (_ []byte, err error) {
	g.doc = `<div><h1>Search result for: "{{.QueryString}}"</h1>
      {{ range .SearchResults }}
      <div class=search-result-card onclick="location.href='.Link';"><h2>{{ join .Context " > "}}</h2><div class=search-result-content>{{ .HTML }}</div></div>
      {{ end }}
    </div>`
	g.queryString = "{{.QueryString}}"

	return g.Build()
}

func GetMarkdownCSS() []byte {
	return []byte(markdownCSS)
}

// nolint: lll
const htmlPageTemplate = `<!DOCTYPE html>
<html lang=en>
<head>
  <title>{{.API}}</title>
  <meta name='generator' content='github.com/karlssonerik/go-service-doc'>
  <link rel="stylesheet" href="{{.Basepath}}/markdown.css">
  {{if .FaviconHref}}<link rel="icon" href="{{.FaviconHref}}">{{end}}
</head>
<body class="markdown-body">
  <div class="flex-container">
    <div class="menu-container">
      <div class=menu-header>
        <h1>{{.API}}</h1>
        <form class=menu-search action="{{.Basepath}}{{.SearchLink}}" method="get">
          <input type="text" placeholder="Search.." name="q" value="{{.QueryString}}" onfocus="var temp_value=this.value; this.value=''; this.value=temp_value" autofocus />
          <button type="submit">Search</button>
        </form>
      </div>
      <div class=menu-content>
        <ul>{{range .Pages}}{{range .Headers}}
          <li><a href="{{.Link}}">{{.Title}}</a>{{if not .Headers}}</li>{{else}}
            <ul>{{end}}{{range .Headers}}
              <li><a href="{{.Link}}">{{.Title}}</a></li>{{end}}{{if .Headers}}
            </ul>
          </li>{{end}}{{end}}{{end}}
        </ul>
      </div>
    </div>
    <div class="doc-container">
      {{.Doc}}
    </div>
  </div>
</body>
</html>`
