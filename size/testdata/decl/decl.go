package decl

import "fmt"

const (
	simpleInt_21 = 1234567890123456789012345678901
	longInt      = 12345678901234567890123456789012
	negativeInt  = -longInt
	shortString  = "Hello, world!"
	longString   = `Here I can finally write my letter to the universe.
			I would love to have the newest and biggest macBook Pro.
			And a nice house in the inner city plus a garage with a Posche inside would be nice, too.
			I have got wife and children already. So I will be humble here.`
	identString = longString
)

var c1_3, c2, c3 chan int

type BankAccount_5 struct {
	HolderName  string
	IBAN        string
	BIC         string
	AccountType string
}

type SpecialBankAccount_3 struct {
	BankAccount_5
	HolderTitle string
}

type mySlice_1 []SpecialBankAccount_3
type myMap_2 map[string][]mySlice_1

var (
	myAccount_18 = BankAccount_5{
		HolderName:  "Dagobert Duck",
		IBAN:        "DE07123412341234123412",
		BIC:         "MARKDEFF",
		AccountType: "credit",
	}
	i, j, k            int
	s1, s2, s3, s4, s5 string
)

func SimpleFunc_4() {
	fmt.Println("vim-go")
}

func (ba *BankAccount_5) DoAccountingMagic_15(newHolder string, newType string) (iban, bic string) {
	ba.HolderName = newHolder
	ba.AccountType = newType

	return ba.IBAN, ba.BIC
}

func (sba *SpecialBankAccount_3) DoSpecialAccountMagic_14(newHolder string, newTitle string) (iban string) {
	sba.DoAccountingMagic_15(newHolder, "special")
	sba.HolderTitle = newTitle
	return sba.IBAN
}

func funcWithEllipsis_18(i, j int, names ...string) []string {
	var addNames []string
	for k := i; k < j; k++ {
		addNames = append(addNames, names[k])
	}

	return append(names, addNames...)
}
