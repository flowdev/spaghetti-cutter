package model

type Product struct {
	Name  string
	Price float64
}

type ShoppingCart struct {
	Content  []Product
	Discount float64
}
