package templates

type Engine struct {
	Ext    string
	Engine Renderer
}
type Renderer interface {
	Render(string, any) (string, error)
}

type Engines struct {
	List []Engine
}
