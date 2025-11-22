package renderer

import (
	"fmt"
	"strings"
	"text/template"

	assets "github.com/Jhonatan-Code-dev/Jmlk-Rpass"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/ports"
)

type TemplateRenderer struct {
	tmplPath string
}

func NewTemplateRenderer() ports.Renderer {
	return &TemplateRenderer{tmplPath: "templates/reset_password.html"}
}

func (r *TemplateRenderer) Render(code string, data map[string]any) (string, error) {
	tmpl, err := template.ParseFS(assets.Templates, r.tmplPath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	if data == nil {
		data = map[string]any{}
	}
	local := map[string]any{}
	for k, v := range data {
		local[k] = v
	}
	local["Code"] = code
	var sb strings.Builder
	if err := tmpl.Execute(&sb, local); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return sb.String(), nil
}
