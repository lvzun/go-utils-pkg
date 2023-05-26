package algorithm

import (
	"fmt"
	"testing"
)

func TestComb(t *testing.T) {
	data := []interface{}{1, 2, 3, 4}
	comb := Comb(data, 3)
	for _, item := range comb {
		fmt.Printf("%v\n", item)
	}
}
