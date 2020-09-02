package domain1

import (
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/db/store"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/x/tool"
)

func HandleDomain1Route1(s *store.Store) {
	tool.Tool()
	s.GetAllProducts()
}

func HandleDomain1Route2(s *store.Store) {
	s.GetShoppingCart()
}
