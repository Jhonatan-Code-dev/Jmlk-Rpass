// internal/domain/generator.go
package domain

import (
	"math/rand"
	"sync"
	"time"
)

const digits = "0123456789"

type CodeGenerator struct {
	length int
	pool   sync.Pool
}

func NewCodeGenerator(length int) *CodeGenerator {
	g := &CodeGenerator{length: length}
	g.pool.New = func() any { return rand.New(rand.NewSource(time.Now().UnixNano())) }
	return g
}

func (g *CodeGenerator) Generate() string {
	r := g.pool.Get().(*rand.Rand)
	defer g.pool.Put(r)
	buf := make([]byte, g.length)
	for i := 0; i < g.length; i++ {
		buf[i] = digits[r.Intn(len(digits))]
	}
	return string(buf)
}
