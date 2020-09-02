package domain2

import (
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/db/store"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/x/tool2"
)

func HandleDomain2Route1(s *store.Store) {
	tool2.Tool2()
	s.GetAllProducts()
}

func HandleDomain2Route2(s *store.Store) {
	ps := s.GetAllProducts()
	s.SaveProduct(ps[0])
}
