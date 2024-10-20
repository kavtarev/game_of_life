package internal

import (
	"fmt"
	"regexp"
	"testing"
)

func TestMain3x3(t *testing.T) {
	schemas := [][2]string{
		{
			`
			000
			111
			000
			`,
			`
			010
			010
			010
			`,
		},
		{
			`
			010
			010
			010
			`,
			`
			000
			111
			000
			`,
		},
		{
			`
			001
			010
			100
			`,
			`
			000
			010
			000
			`,
		},
		{
			`
			011
			011
			000
			`,
			`
			011
			011
			000
			`,
		},
		{
			`
			111
			001
			001
			`,
			`
			011
			001
			000
			`,
		},
		{
			`
			101
			010
			101
			`,
			`
			010
			101
			010
			`,
		},
	}

	for i := 0; i < len(schemas); i++ {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			f := NewField(3)
			reg := regexp.MustCompile(`\D`)
			init := reg.ReplaceAllString(schemas[i][0], "")
			f.Update(init)
			str := f.Run()
			if reg.ReplaceAllString(schemas[i][1], "") != str {
				t.Fatal()
			}
		})

	}
}

func TestMain4x4(t *testing.T) {
	schemas := [][2]string{
		{
			`
			1111
			1100
			1010
			1001
			`,
			`
			1010
			0001
			1010
			0100
			`,
		},
		{
			`
			1011
			1100
			0010
			1101
			`,
			`
			1010
			1001
			0010
			0110
			`,
		},
	}

	for i := 0; i < len(schemas); i++ {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			f := NewField(4)
			reg := regexp.MustCompile(`\D`)
			init := reg.ReplaceAllString(schemas[i][0], "")
			f.Update(init)
			str := f.Run()
			if reg.ReplaceAllString(schemas[i][1], "") != str {
				t.Fatal()
			}
		})

	}
}

func TestMain5x5(t *testing.T) {
	schemas := [][2]string{
		{
			`
			10111
			10100
			11111
			00101
			11101
			`,
			`
			00110
			10000
			10001
			00001
			01100
			`,
		},
		{
			`
			00100
			01010
			10001
			01010
			00100
			`,
			`
			00100
			01110
			11011
			01110
			00100
			`,
		},
	}

	for i := 0; i < len(schemas); i++ {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			f := NewField(5)
			reg := regexp.MustCompile(`\D`)
			init := reg.ReplaceAllString(schemas[i][0], "")
			f.Update(init)
			str := f.Run()
			if reg.ReplaceAllString(schemas[i][1], "") != str {
				t.Fatal()
			}
		})

	}
}
