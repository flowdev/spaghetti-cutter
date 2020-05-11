package domain4

import (
	"github.com/flowdev/spaghetti-cutter/deps/testdata/complex-proj/pkg/db/store"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/complex-proj/pkg/domain3"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/complex-proj/pkg/x/tool2"
)

func HandleDomain4Route1(s *store.Store) {
	tool2.Tool2()
	s.GetAllProducts()
}

func HandleDomain4Route2(s *store.Store) {
	domain3.HandleDomain3Route1(s)
}
