package goab

import (
	"net/http"
	"io"
	"bufio"
	"os"
	"fmt"
)

type PipeRequestBuilder struct {
	input         io.Reader
}

func NewPipeRequestBuilder(input io.Reader) *PipeRequestBuilder{
	return &PipeRequestBuilder{
		input: input,
	}
}
func(t *PipeRequestBuilder) buildRequest()(request *http.Request,err error){

	scanner := bufio.NewScanner(t.input)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

