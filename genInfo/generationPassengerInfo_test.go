package genInfo

import (
	"fmt"
	"testing"
)

func TestGenIdNo(t *testing.T) {

	for i := 0; i < 100; i++ {
		sex := "M"
		if i%2 == 0 {
			sex = "F"
		}
		idNo := GenIdNoWithAge(26, i%2 == 0)
		name := GenName()
		phone := GenPhoneNumber()
		birthday := ResolveBirthDayFromIdNo(idNo)
		fmt.Printf("Idno:%s,birthday:%s,phone:%s,name:%s,sex:%s\n", idNo, birthday, phone, name, sex)
	}
}
