package glowflake

import (
	"fmt"
	"testing"
)

func TestNextId(t *testing.T) {
	gf, _ := NewGlowFlake(1, 1)

	id, _ := gf.NextId()
	fmt.Printf("id generated:%d\n", id)
}

func BenchmarkNextId(b *testing.B) {
	gf, _ := NewGlowFlake(1, 1)
	for n := 0; n < b.N; n++ {
		gf.NextId()
	}
}
