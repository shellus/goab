package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/shellus/goab"
)

var ab *goab.Goab

func main() {
	app := cli.NewApp()
	app.Name = "go version ApacheBench"
	app.Usage = "ab [options] [http[s]://]hostname[:port]/path"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "method, m",
			Value: "GET",
			Usage: "request `method`",
		},
		cli.StringFlag{
			Name:  "concurrency, c",
			Usage: "Number of multiple requests to make at a time",
		},
		cli.StringFlag{
			Name:  "timelimit, t",
			Usage: "Seconds to max. to spend on benchmarking\n\tThis implies -n 50000",
		},
		cli.StringFlag{
			Name:  "data, d",
			Usage: "HTTP POST data",
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "Pass custom header LINE to server",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Printf("m %s\n", c.String("m"))
		fmt.Printf("c %s\n", c.String("c"))
		fmt.Printf("t %s\n", c.String("t"))
		//fmt.Printf("d %s\n", c.String("d"))

		for _, header := range c.StringSlice("H") {
			fmt.Printf("H %s\n", header)
		}
		fmt.Printf("args %s\n", c.Args().Get(0))

		headers := c.StringSlice("H")
		ab = goab.New(c.Args().Get(0), headers, c.String("m"), c.Int("c"), c.Uint("t"))
		ab.Run()
		ab.Wait()
		fmt.Println(ab.Counter.Dump())

		for i, process := range ab.Process.Dump() {
			fmt.Printf("% 3d%%\t\t%s\n", (i+1)*10, process.String())
		}
		return nil
	}

	app.Run(os.Args)
}
