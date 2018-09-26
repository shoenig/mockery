package libmockery

import (
	"go/types"
	"strings"
)

type Interface struct {
	Name string
	// Path      string
	Type      *types.Interface
	NamedType *types.Named
}

type sortableIFaceList []*Interface

func (s sortableIFaceList) Len() int {
	return len(s)
}

func (s sortableIFaceList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableIFaceList) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) == -1
}
