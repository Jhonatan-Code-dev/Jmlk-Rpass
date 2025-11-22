package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/app"
	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/domain"
)

func TestSendReset_Success(t *testing.T) {
	repo := NewRepoMock()
	smtp := &SMTPMock{}
	render := &RendererMock{}

	cfg := app.Config{
		Username:         "from@test",
		AppName:          "MiApp",
		Title:            "Reset",
		CodeLength:       6,
		CodeValidMinutes: 10,
		MaxResetAttempts: 3,
		RestrictionHours: 24,
		AllowOverride:    true,
		EmailTimeout:     5 * time.Second,
	}

	svc := app.NewService(cfg, repo, smtp, render)

	ctx := context.Background()
	if err := svc.SendReset(ctx, "u@test"); err != nil {
		t.Fatalf("expected success, got err=%v", err)
	}
	if smtp.Last.To != "u@test" {
		t.Fatalf("smtp not called with expected to")
	}
	if render.Last == "" {
		t.Fatalf("render not called")
	}
}

func TestSendReset_Blocked(t *testing.T) {
	repo := NewRepoMock()
	repo.SaveCode(context.Background(), domain.CodeEntry{
		Email:    "u2@test",
		Code:     "111111",
		Attempts: 5,
		Used:     false,
		ExpireAt: time.Now().Add(-1 * time.Hour),
	})
	smtp := &SMTPMock{}
	render := &RendererMock{}

	cfg := app.Config{
		Username:         "from@test",
		CodeLength:       6,
		CodeValidMinutes: 15,
		MaxResetAttempts: 3,
		RestrictionHours: 24,
		AllowOverride:    false,
		EmailTimeout:     5 * time.Second,
	}

	svc := app.NewService(cfg, repo, smtp, render)
	err := svc.SendReset(context.Background(), "u2@test")
	if err == nil {
		t.Fatalf("expected error due policy, got nil")
	}
}
