package rules

import (
	"math/rand"
	"sync"
	"time"
)

const digits = "0123456789"

type CodeGenerator struct {
	codeLength int
	prngPool   sync.Pool
}

func NewCodeGenerator(length int) *CodeGenerator {
	g := &CodeGenerator{codeLength: length}
	g.prngPool.New = func() any { return rand.New(rand.NewSource(time.Now().UnixNano())) }
	return g
}

func (g *CodeGenerator) Generate() string {
	r := g.prngPool.Get().(*rand.Rand)
	var buf [64]byte
	for i := 0; i < g.codeLength; i++ {
		buf[i] = digits[r.Intn(10)]
	}
	g.prngPool.Put(r)
	return string(buf[:g.codeLength])
}
