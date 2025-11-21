package ports

type Renderer interface {
	Render(code string, data map[string]any) (string, error)
}
