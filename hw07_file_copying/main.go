package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
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

	var wg sync.WaitGroup
	progress := make(chan int64)

	wg.Add(1)
	go func() {
		defer wg.Done()

		inFile, err := os.Open(from)
		if err != nil {
			fmt.Println("Failed to open input file:", err)
			close(progress)
			return
		}
		defer inFile.Close()

		totalSize, err := getSize(inFile) // Используем getSize
		if err != nil {
			fmt.Println("Failed to get file size:", err)
			close(progress)
			return
		}

		if limit == 0 || limit > totalSize-offset {
			limit = totalSize - offset
		}

		bar := pb.Start64(limit)
		bar.SetRefreshRate(time.Millisecond * 100)
		defer bar.Finish()
		for p := range progress {
			bar.SetCurrent(p)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Copy(from, to, offset, limit, progress)
		if err != nil {
			fmt.Println("Error during copy:", err)
		}
	}()

	wg.Wait()

	fmt.Println("Progress completed!")
}
