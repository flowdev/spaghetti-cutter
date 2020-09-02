package store

import (
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/db/model"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/db/store/substore"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/x/tool"
	"github.com/flowdev/spaghetti-cutter/deps/testdata/half-pkgs-proj/x/tool2"
)

type Store struct {
	// no real DB so no real fields either
}

func (s *Store) GetAllProducts() []model.Product {
	tool.Tool()
	return []model.Product{
		{Name: "prod1", Price: 12.34},
		{Name: "prod2", Price: 56.78},
	}
}

func (s *Store) SaveProduct(prod model.Product) error {
	tool2.Tool2()
	return nil
}

func (s *Store) GetShoppingCart() model.ShoppingCart {
	substore.GetTechnicalThing()
	return model.ShoppingCart{
		Content: []model.Product{
			{Name: "prod1", Price: 12.34},
			{Name: "prod2", Price: 56.78},
		},
		Discount: 0.9,
	}
}
