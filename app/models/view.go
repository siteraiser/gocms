package models

type View struct {
	Args   any
	Output string
	Render func(source string, args any) string
}
