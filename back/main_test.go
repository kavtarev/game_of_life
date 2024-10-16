package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestMain(t *testing.T) {
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
			f.update(init)
			f.run()
			if reg.ReplaceAllString(schemas[i][1], "") != f.getStateString() {
				t.Fatal()
			}
		})

	}

}
