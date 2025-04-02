package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
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

	if from == "" || to == "" {
		fmt.Println("Both -from and -to arguments are required!")
		return
	}

	var wg sync.WaitGroup
	progress := make(chan int64)

	wg.Add(1)
	go func() {
		defer wg.Done()

		inFile, err := os.Open(from)
		if err != nil {
			fmt.Println("Failed to open input file:", err)
			return
		}
		defer inFile.Close()

		totalSize, err := getSize(inFile) // Используем getSize
		if err != nil {
			fmt.Println("Failed to get file size:", err)
			return
		}

		if limit == 0 || limit > totalSize-offset {
			limit = totalSize - offset
		}

		bar := NewProgressBar(limit)
		defer bar.Finish()
		for p := range progress {
			bar.Update(p)
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
}
