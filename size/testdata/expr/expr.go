package expr

import (
	"os"
	"regexp"

	"github.com/flowdev/spaghetti-cutter/x/config"
)

const ( // literals, unary and ident
	simpleInt_1   = 1234567890123456789012345678901
	longInt_2     = 12345678901234567890123456789012
	negativeInt_1 = -longInt_2
	shortString_1 = "Hello, world!"
	longString_9  = `Here I can finally write my letter to the universe.
I would love to have the newest and biggest macBook Pro.
And a nice house in the inner city plus a garage with a Posche inside would be nice, too.
I have got wife and children already. So I will be humble here.`
	identString_1 = longString_9
)

var ( // complex expressions, arrays, selectors, struct type
	stringSlice_5 = []string{
		"Hi there!", "I am here.", "Welcome to the world.", "Good bye!",
	}

	structType_1 = struct{}{}

	configTool_16 = config.PatternList([]config.Pattern{
		{Pattern: "x/*", Regexp: regexp.MustCompile("x/*")},
		{Pattern: "almosttool", Regexp: regexp.MustCompile("almosttool")},
	})
	configDB_10 = config.PatternList([]config.Pattern{
		{Pattern: "db/*", Regexp: regexp.MustCompile("db/*")},
	})
	configGod_10 = config.PatternList([]config.Pattern{
		{Pattern: "main", Regexp: regexp.MustCompile("main")},
	})

	configLiteral_16 = config.Config{
		Allow:        nil,
		Tool:         &configTool_16,
		DB:           &configDB_10,
		God:          &configGod_10,
		Size:         123,
		NoGod:        true,
		IgnoreVendor: true,
	}
)

var ( // ellipsis?
	ellipsisExpr_3 = append([]string{}, stringSlice_5...)
)

var ( // struct type
	structType_13 = struct {
		i    int `tag: "bla"`
		j    int
		s, t string
	}{
		i: 1, j: 2,
		s: "I", t: "am",
	}
)

var ( // map type
	mapType_10 = map[int]string{1: "one", 2: "two", 3: "three", 4: "many"}
)

var ( // chan
	channel_3  = make(chan<- int, 100)
	channel2_3 = make(<-chan string, 10)
)

var ( // interface, func type, func lit
	interfaceType_2  = interface{}(1)
	interfaceType_13 = interface {
		Read([]byte) (int, error)
		Write(b []byte) (n int, err error)
		Close() error
	}(os.Stdout)

	funcLit_7 = func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
)

var ( // assert, star, slice, index
	assertedInt_2 = interfaceType_2.(int)
	start_2       = *configLiteral_16.DB
	sliceExpr_3   = stringSlice_5[1:2]
	indexExpr_5   = stringSlice_5[1+2+3-4]
)

var ( // binary and paren expressions
	binary_2 = 1 + 2
	paren_3  = (1 + (3 - 2))

	binary_4 = 1 + 2 - 3 + 4
	paren_4  = (1 + (2 - (3 + 4)))
)
