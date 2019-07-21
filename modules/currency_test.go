package modules

import (
	"fmt"
	"testing"
)

func init() {
	//util.ConfigInit("../conf")
	//util.LoggerInit()
}
func Test_getCurrencies(t *testing.T) {
	c := NewCurrency()
	list, err := c.getCurrenciesFromGate()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(list)
}
