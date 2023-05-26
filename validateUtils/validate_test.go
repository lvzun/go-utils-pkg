package validateUtils

import (
	"fmt"
	"testing"
)

func TestVerifyIdNo(t *testing.T) {
	no := VerifyIdNo("430105198503044317")
	fmt.Printf("no:%v", no)

}
