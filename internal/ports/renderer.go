// internal/ports/renderer.go
package ports

type Renderer interface {
	Render(code string, data map[string]any) (string, error)
}
