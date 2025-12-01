package helpers

import "gocms/app/helpers/html"

type Html struct {
	Href interface{ Ahref() }
}

func (h Html) Ahref(href string, text string) string {
	return html.Ahref(href, text)
}

type Grammar struct {
	UPluralS interface {
		UpperIfPluralS(l int) string
	}
	LPluralS interface {
		UpperIfPluralS(l int) string
	}
}

// Returns an s if not equal to 1
func (h Grammar) UpperIfPluralS(number int) string {
	if number != 1 {
		return "S"
	}
	return ""
}

// Returns an s if not equal to 1
func (h Grammar) LowerIfPluralS(number int) string {
	if number != 1 {
		return "s"
	}
	return ""
}

/*
type Presets struct {
	PluralS interface {
		GetPluralS(href string, text string) string
	}
}


	s := func(l int) string {
				if l != 1 {
					return "s"
				}
				return ""
			}(l)

*/
