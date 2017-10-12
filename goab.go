package goab

import (
	"time"
	"context"
	"sync"
	"net/http"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
)
type RequestBuilder interface {
	buildRequest()(*http.Request,error)
}
type Goab struct {
	concurrency int
	secondLimit uint
	CancelFunc  context.CancelFunc
	deadContext context.Context
	waitSync    sync.WaitGroup
	requestBuilder RequestBuilder
	Counter     *SafeCounter
	Process     *process

	//seedRequest chan *http.Request
	//feedbackResp chan *http.Response
}

func New(requestBuilder RequestBuilder, concurrency int, secondLimit uint) *Goab {

	return &Goab{
		requestBuilder: requestBuilder,
		concurrency: concurrency,
		secondLimit: secondLimit,
		Counter:     NewCounter(),
		Process:     NewProcess(),
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
			//fmt.Println("dead...")
			break Quit
		default:
			startTime := time.Now()
			request, err := t.requestBuilder.buildRequest()
			if err != nil {
				panic(err)
			}
			t.Counter.Inc("Send")
			resp, err := ctxhttp.Do(t.deadContext, http.DefaultClient, request)
			if err != nil {
				if err == context.DeadlineExceeded {
					// this context err
					t.Counter.Inc("Deadline")
					break
				} else {
					// not context err
					panic(err)
				}
			} else {

				if resp.StatusCode == 200 {
					t.Counter.Inc("200")
				} else {
					t.Counter.Inc("not 200")
				}

				_, err := ioutil.ReadAll(resp.Body)
				processTime := time.Now().Sub(startTime)
				t.Process.Add(processTime)
				if err != nil {
					t.Counter.Inc("ReadFail")
				} else {
					t.Counter.Inc("ReadDone")
				}
				resp.Body.Close()
			}
			//t.feedbackResp <- resp
		}
	}
	t.waitSync.Done()
}
