package stmt

import "fmt"

func simpleAssignment_4() { a := 2; _ = a }

func incDecStmt_4() { a := 3; a++; a-- }

func returnStmt_3() (int, string, int) { return 3, "foo", 4 }

func exprStmt_1() { incDecStmt_4() }

func ifStmt_6() {
	if i := 5; i < 6 {
		i++
	} else {
		i--
	}
}

func forStmt_8() {
	for i := 0; i < 8; i++ {
		fmt.Println("Hi!")
	}
}

func rangeStmt_6() {
	for _, c := range "foo and bar" {
		fmt.Println(c)
	}
}

func blockStmt_13() {
	if 4 > 5 {
	} else {
		fmt.Println("Hey!")
		fmt.Println(4*5 + 6)
		fmt.Println("Ho!")
	}
}

func switchStmt_19() { // and case, too
	switch "a" {
	case "a":
		a := 6 + 7*8
		fmt.Println("Hey!", a)
	case "b":
		fmt.Println(4*5 + 6)
	default:
		fmt.Println("Ho!")

	}
}

func typeSwitchStmt_20() { // and case, too
	var v interface{}
	switch v.(type) {
	case string:
		a := 6 + 7*8
		fmt.Println("Hey!", a)
	case nil:
		fmt.Println(4*5 + 6)
	default:
		fmt.Println("Ho!")

	}
}

func typeSelectStmt_35() { // comm, decl and send, too
	var a []int
	var c1, c2, c3, c4 chan int
	var i1, i2 int
	select {
	case i1 = <-c1:
		print("received ", i1, " from c1\n")
	case c2 <- i2:
		print("sent ", i2, " to c2\n")
	case i3, ok := (<-c3):
		if ok {
			print("received ", i3, " from c3\n")
		} else {
			print("c3 is closed\n")
		}
	case a[3] = <-c4:
	default:
		print("no communication\n")
	}
}

func branchStmt_11() {
	for _, r := range "foo and bar" {
		if r == 'o' {
			continue
		}
		print("got rune\n")
		if r == 'r' {
			break
		}
	}
}

func goDeferStmt_6() {
	defer func() { print("foo\n") }()
	go typeSelectStmt_35()
}

func labeledStmt_9() {
Loop:
	for {
		if 3 > 7 {
			continue Loop
		}
		print("bar\n")
		if 5 < 6 {
			break Loop
		}
	}
}
