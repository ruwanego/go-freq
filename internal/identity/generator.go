package identity

import "fmt"

type Generator struct {
	enginePrefix string
	sequences    map[string]int64
}

func NewGenerator(prefix string, initial map[string]int64) *Generator {
	seqs := make(map[string]int64, len(initial))
	for k, v := range initial {
		seqs[k] = v
	}

	return &Generator{
		enginePrefix: prefix,
		sequences:    seqs,
	}
}

func (g *Generator) Next(strategy string, ts int64) string {
	seq := g.sequences[strategy] + 1
	g.sequences[strategy] = seq

	return fmt.Sprintf("%s-%s-%d-%04d", g.enginePrefix, strategy, ts, seq)
}
