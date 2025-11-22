package tests

import (
	"context"
	"sync"

	"github.com/Jhonatan-Code-dev/Jmlk-Rpass/internal/domain"
)

type RepoMock struct {
	mu    sync.Mutex
	store map[string]domain.CodeEntry
	Err   error
}

func NewRepoMock() *RepoMock { return &RepoMock{store: map[string]domain.CodeEntry{}} }

func (r *RepoMock) SaveCode(ctx context.Context, e domain.CodeEntry) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[e.Email] = e
	return nil
}

func (r *RepoMock) GetCodeEntry(ctx context.Context, email string) (*domain.CodeEntry, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.store[email]
	if !ok {
		return nil, nil
	}
	return &v, nil
}

func (r *RepoMock) Close() error { return nil }

type SMTPMock struct {
	Last struct{ From, To, Subject, Body string }
	Err  error
}

func (m *SMTPMock) Send(ctx context.Context, from, to, subject, htmlBody string) error {
	if m.Err != nil {
		return m.Err
	}
	m.Last.From = from
	m.Last.To = to
	m.Last.Subject = subject
	m.Last.Body = htmlBody
	return nil
}

type RendererMock struct {
	Last string
	Err  error
}

func (r *RendererMock) Render(code string, data map[string]any) (string, error) {
	if r.Err != nil {
		return "", r.Err
	}
	out := "HTML:" + code
	r.Last = out
	return out, nil
}
