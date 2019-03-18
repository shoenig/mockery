package libmockery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type GatheringVisitor struct {
	Interfaces []*Interface
}

func (gv *GatheringVisitor) VisitWalk(iface *Interface) error {
	gv.Interfaces = append(gv.Interfaces, iface)
	return nil
}

func NewGatheringVisitor() *GatheringVisitor {
	return &GatheringVisitor{
		Interfaces: make([]*Interface, 0, 1024),
	}
}

func TestWalker_gather_cwd(t *testing.T) {
	w := Walker{
		Verbose:   true,
		BaseDir:   ".",
		Interface: "WalkerVisitor",
	}

	gv := NewGatheringVisitor()

	w.Walk(gv)

	require.Equal(t, 1, len(gv.Interfaces))
	first := gv.Interfaces[0]
	require.Equal(t, "WalkerVisitor", first.Name)
}

func TestWalker_gather_subdir(t *testing.T) {
	w := Walker{
		Verbose:   true,
		BaseDir:   "fixtures",
		Interface: "AsyncProducer",
	}

	gv := NewGatheringVisitor()

	w.Walk(gv)

	require.Equal(t, 1, len(gv.Interfaces))
	first := gv.Interfaces[0]
	require.Equal(t, "AsyncProducer", first.Name)
}
