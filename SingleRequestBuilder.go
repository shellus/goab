package goab

import "net/http"

type SingleRequestBuilder struct {
	url         string
	headers     []string
	method      string
}

func NewSingleRequestBuilder(url string, headers []string, method string) *SingleRequestBuilder{
	return &SingleRequestBuilder{
		url:         url,
		headers:     headers,
		method:      method,
	}
}
func(t *SingleRequestBuilder) buildRequest()(request *http.Request,err error){
	request, err = http.NewRequest(t.method, t.url, nil)
	return
}