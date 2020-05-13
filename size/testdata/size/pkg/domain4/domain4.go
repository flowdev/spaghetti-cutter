package domain4

import (
	"fmt"

	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/db/store"
	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/domain3"
	"github.com/flowdev/spaghetti-cutter/testdata/good-proj/pkg/x/tool2"
)

func HandleDomain4Route1(s *store.Store) {
	tool2.Tool2()
	s.GetAllProducts()
}

func HandleDomain4Route2(s *store.Store) {
	domain3.HandleDomain3Route1(s)
	foo(1, 2, 3, 4, "s1", "s2", "s3")
}

func foo(i1, i2, i3, i4 int, s1, s2, s3 string) {
	if i1 > 0 {
		fmt.Println("I sum:", i1+i2+i3+i4)
	} else {
		fmt.Println("S sum:", s1+s2+s3)
	}

	for i := 0; i < i2; i++ {
		go bar(i, s3)
	}
}

func bar(i int, s string) {
	switch i % 2 {
	case 0:
		fmt.Println("You are even:", i, s)
	case 1:
		fmt.Println("That is odd:", i)
	default:
		fmt.Println("What's this???", i)
	}
}
