package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type HttpCli struct {
	Cli      *http.Client
	body     []byte
	u        string
	queries  url.Values
	document *goquery.Document
}

var (
	httpCli          *http.Client
	DefaultTransport *http.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     30 * time.Second,
	}
	onceCli sync.Once
)

func NewClient() *HttpCli {
	onceCli.Do(func() {
		httpCli = &http.Client{
			Transport: DefaultTransport,
			Timeout:   60 * time.Second,
		}
	})
	return &HttpCli{
		Cli:     httpCli,
		queries: make(url.Values),
	}
}

func (hc *HttpCli) ParseUrl(uri string) *HttpCli {
	hc.u = uri
	u, err := url.Parse(uri)
	if err != nil {
		return hc
	}
	hc.u = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	for k, v := range u.Query() {
		hc.queries[k] = v
	}
	return hc
}

func (hc *HttpCli) Url() (u string) {
	if hc.u == "" {
		return
	}
	if len(hc.queries) > 0 {
		return fmt.Sprintf("%s?%s", hc.u, hc.queries.Encode())
	}
	return hc.u
}

func (hc *HttpCli) Query(k, v string) *HttpCli {
	hc.queries[k] = []string{v}
	return hc
}

func (hc *HttpCli) Do(method string) error {
	resp, err := hc.do(method)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read data from response error: %s", err.Error())
	}
	hc.body = body
	return err
}
func (hc *HttpCli) Html(method string) error {
	resp, err := hc.do(method)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("unmarshal document error: %s", err.Error())
	}
	hc.document = doc
	return nil
}
func (hc *HttpCli) do(method string) (*http.Response, error) {
	var (
		buf  bytes.Buffer
		resp *http.Response
		code int
	)
	req, err := http.NewRequest(method, hc.Url(), &buf)
	if err != nil {
		return nil, fmt.Errorf("create req error: %s", err.Error())
	}
	req.Close = true
	resp, err = hc.Cli.Do(req)
	if err != nil {
		if resp != nil {
			code = resp.StatusCode
		}
		return nil, fmt.Errorf("do request error: %s, code: %d", err.Error(), code)
	}
	return resp, nil
}
func (hc *HttpCli) Body() []byte {
	return hc.body
}

func (hc *HttpCli) Json(i interface{}) error {
	if len(hc.body) == 0 {
		return fmt.Errorf("body is empty")
	}
	return json.Unmarshal(hc.body, i)
}

func (hc *HttpCli) Document() *goquery.Document {
	return hc.document
}
