package models

type View struct {
	RequestId string
	Args      any
	Output    string
	Render    func(source string, args any) string
}
