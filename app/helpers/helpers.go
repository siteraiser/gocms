package helpers

import "gocms/app/helpers/html"

type Html struct {
	Href interface{ Ahref() }
}

func (h Html) Ahref(href string, text string) string {
	return html.Ahref(href, text)
}

type Grammar struct {
	PluralS interface {
		PluralS(href string, text string) string
	}
}

func (h Html) PluralS(href string, text string) string {
	return html.Ahref(href, text)
}
