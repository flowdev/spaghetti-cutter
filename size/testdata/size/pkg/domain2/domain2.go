package domain2

import (
	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/db/store"
	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/x/tool2"
)

func HandleDomain2Route1(s *store.Store) {
	tool2.Tool2()
	s.GetAllProducts()
}

func HandleDomain2Route2(s *store.Store) {
	ps := s.GetAllProducts()
	s.SaveProduct(ps[0])
}
