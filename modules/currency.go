package modules

import (
	"fmt"
	"gate/models"
	"gate/util"
	"net/http"
)

type currency struct {
	host string
}

func NewCurrency() *currency {
	return &currency{host: "https://www.gateio.io/json_svr/get_leftbar"}
}
func (c *currency) getCurrenciesFromGate() ([]*models.JsonCurrency, error) {
	cli := util.NewClient()
	err := cli.ParseUrl(c.host).Query("u", "128").Query("c", "800264").Do(http.MethodGet)
	if err != nil {
		return nil, err
	}
	var currencies map[string]map[string]*models.JsonCurrency
	if err := cli.Json(&currencies); err != nil {
		return nil, err
	}
	var list []*models.JsonCurrency
	for _, v := range currencies["USDT"] {
		//volume, _ := strconv.ParseFloat(v.Volume, 64)
		// if volume < 10000.00 {
		// 	continue
		// }
		v.FormatName()
		list = append(list, v)
	}
	return list, nil
}
func getCurrenciesFromDB() (map[string]*models.Currency, error) {
	var list []*models.Currency
	if err := models.AllIDs(&list); err != nil {
		return nil, err
	}
	m := make(map[string]*models.Currency)
	for _, v := range list {
		m[v.Name] = v
	}
	return m, nil
}
func setTDBEvent(listDB map[string]*models.Currency, list []*models.JsonCurrency) {
	for _, jc := range list {
		if _, ok := listDB[jc.Name]; ok {
			listDB[jc.Name] = jc.ToCurrency(models.UPDATE)
		} else {
			listDB[jc.Name] = jc.ToCurrency(models.INSERT)
		}
	}
}
func (c *currency) Sync() []*models.JsonCurrency {
	// listDB, err := getCurrenciesFromDB()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	list, err := c.getCurrenciesFromGate()
	if err != nil {
		fmt.Println(err)
	}
	// setTDBEvent(listDB, list)
	// for _, v := range list {
	// 	switch v.DBEvent {
	// 	case models.INSERT:
	// 		if err := v.Insert(); err != nil {
	// 			fmt.Printf("insert curreny error: %s, data: %+v\n", err.Error(), v)
	// 		}
	// 		break
	// 	case models.UPDATE:
	// 		if err := v.Update(); err != nil {
	// 			fmt.Printf("update curreny error: %s, data: %+v\n", err.Error(), v)
	// 		}
	// 		break
	// 	case models.DELETE:
	// 		if err := v.Delete(); err != nil {
	// 			fmt.Printf("delete curreny error: %s, data: %+v\n", err.Error(), v)
	// 		}
	// 		break
	// 	}
	// }
	return list
}
