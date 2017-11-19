package libmockery

import (
	"os"
	"path/filepath"
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
		BaseDir:   ".",
		Interface: "WalkerVisitor",
	}

	gv := NewGatheringVisitor()

	w.Walk(gv)

	require.Equal(t, 1, len(gv.Interfaces))
	first := gv.Interfaces[0]
	require.Equal(t, "WalkerVisitor", first.Name)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	path := filepath.Join(cwd, "walker.go")
	require.NoError(t, err)

	require.Equal(t, path, first.Path)
}

func TestWalker_gather_subdir(t *testing.T) {
	w := Walker{
		BaseDir:   "fixtures",
		Interface: "AsyncProducer",
	}

	gv := NewGatheringVisitor()

	w.Walk(gv)

	require.Equal(t, 1, len(gv.Interfaces))
	first := gv.Interfaces[0]
	require.Equal(t, "AsyncProducer", first.Name)
	require.Equal(t, getFixturePath("async.go"), first.Path)
}
