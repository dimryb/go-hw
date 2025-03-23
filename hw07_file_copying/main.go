package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	fmt.Println("From:", from)
	fmt.Println("To:", to)
	fmt.Println("Limit:", limit)
	fmt.Println("Offset:", offset)

	count := 10000

	var wg sync.WaitGroup
	progress := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			progress <- 1
			time.Sleep(time.Millisecond)
		}
		close(progress)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bar := pb.StartNew(count)
		bar.SetRefreshRate(time.Millisecond * 100)
		for v := range progress {
			bar.Add(v)
		}
		bar.Finish()
		fmt.Println("Bar Finished!")
	}()

	wg.Wait()

	fmt.Println("Progress completed!")
}
