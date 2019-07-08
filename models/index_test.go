package models

import (
	"test_services/util"
	"testing"
)

func init() {
	util.ConfigInit("../conf")
	util.LoggerInit()
	util.MysqlInit()
}
func Test_DeleteIndexByWhere(t *testing.T) {
	err := DeleteIndexByWhere("abt_usdt", 1562495400)
	if err != nil {
		t.Fatal(err)
	}
}
