package domain3

import (
	"github.com/flowdev/spaghetti-cutter/deps/testdata/complex-proj/pkg/db/store"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/complex-proj/pkg/x/tool"
)

func HandleDomain3Route1(s *store.Store) {
	tool.Tool()
	s.GetAllProducts()
}

func HandleDomain3Route2(s *store.Store) {
	s.GetShoppingCart()
}
