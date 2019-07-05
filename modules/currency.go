package modules

import (
	"net/http"
	"test_services/models"
	"test_services/util"
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
func (c *currency) Sync() map[string]*models.Currency {
	listDB, err := getCurrenciesFromDB()
	if err != nil {
		util.Logger().Error(err)
	}
	list, err := c.getCurrenciesFromGate()
	if err != nil {
		util.Logger().Error(err)
	}
	setTDBEvent(listDB, list)
	for _, v := range listDB {
		switch v.DBEvent {
		case models.INSERT:
			if err := v.Insert(); err != nil {
				util.Logger().ErrorF("insert curreny error: %s, data: %+v", err.Error(), v)
			}
			break
		case models.UPDATE:
			if err := v.Update(); err != nil {
				util.Logger().ErrorF("update curreny error: %s, data: %+v", err.Error(), v)
			}
			break
		case models.DELETE:
			if err := v.Delete(); err != nil {
				util.Logger().ErrorF("delete curreny error: %s, data: %+v", err.Error(), v)
			}
			break
		}
	}
	return listDB
}
