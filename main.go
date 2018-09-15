package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/convox/stdapi"
)

var categorySlugs = []string{
	"introduction",
	"deployment",
	"development",
	"management",
	"monitoring",
	"resources",
	"integrations",
	"console",
	"reference",
	"gen1",
}

func init() {
	if err := LoadCategories(categorySlugs...); err != nil {
		log.Fatal(err)
	}
}

func main() {
	s := stdapi.New("docs", "docs.convox")

	s.Router.Static("/assets/", "./assets")

	s.Route("GET", "/", index)
	s.Route("GET", "/{category}/{slug}", doc)
	s.Route("GET", "/{slug}", redirect)

	stdapi.LoadTemplates("./templates", helpers)

	if err := s.Listen("https", ":3000"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func helpers(c *stdapi.Context) template.FuncMap {
	return template.FuncMap{}
}

func index(c *stdapi.Context) error {
	return c.Redirect(302, "/introduction/overview")
}

func doc(c *stdapi.Context) error {
	cc, ok := categories.Find(c.Var("category"))
	if !ok {
		return stdapi.Errorf(404, "not found")
	}

	d, ok := cc.Documents.Find(c.Var("slug"))
	if !ok {
		return stdapi.Errorf(404, "not found")
	}

	params := map[string]interface{}{
		"Categories": categories,
		"Category":   cc.Slug,
		"Document":   template.HTML(d.Body),
		"Slug":       d.Slug,
		"Title":      d.Title,
	}

	return c.RenderTemplate("doc", params)
}

func redirect(c *stdapi.Context) error {
	for _, cc := range categories {
		if d, ok := cc.Documents.Find(c.Var("slug")); ok {
			return c.Redirect(301, fmt.Sprintf("/%s/%s", cc.Slug, d.Slug))
		}
	}

	return stdapi.Errorf(404, "not found")
}
