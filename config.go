package main

type Config struct {
	TemplateEngine ViewRenderer
}

type ViewRenderer struct {
	Name            string
	ViewRenderer    string
	ViewRendererExt string
}
