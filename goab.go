package goab

import (
	"time"
	"context"
	"sync"
	"fmt"
	"net/http"
	"golang.org/x/net/context/ctxhttp"
)

type Goab struct {
	url         string
	headers     []string
	method      string
	concurrency int
	secondLimit uint
	CancelFunc  context.CancelFunc
	deadContext context.Context
	waitSync    sync.WaitGroup
	//seedRequest chan *http.Request
	//feedbackResp chan *http.Response
}

func New(url string, headers []string, method string, concurrency int, secondLimit uint) *Goab {

	return &Goab{
		url:         url,
		headers:     headers,
		method:      method,
		concurrency: concurrency,
		secondLimit: secondLimit,
		//feedbackResp: make(chan *http.Response),
	}
}

func (t *Goab) Run() {
	t.deadContext, t.CancelFunc = context.WithTimeout(context.Background(), time.Second*time.Duration(t.secondLimit))

	t.waitSync = sync.WaitGroup{}
	t.waitSync.Add(t.concurrency)

	for i := 0; i < t.concurrency; i++ {
		go t.requestPool()
	}
}

func (t *Goab) Cancel() {
	t.CancelFunc()
}

func (t *Goab) Wait() {
	t.waitSync.Wait()
	t.CancelFunc()
}

func (t *Goab) requestPool() {
Quit:
	for {
		select {
		case <-t.deadContext.Done():
			fmt.Println("dead...")
			break Quit
		default:
			request,err:=http.NewRequest(t.method,t.url,nil)
			if err != nil {
				panic(err)
			}

			resp,err := ctxhttp.Do(t.deadContext, http.DefaultClient, request)
			if err != nil {
				if err == context.DeadlineExceeded {

				}else {
					fmt.Println(request.Context().Err())
					fmt.Println(err)
					panic(err)
				}
			}else {
				fmt.Println(resp.Status)
			}
			//t.feedbackResp <- resp
		}
	}
	t.waitSync.Done()
}
