package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var closeErr error

	inFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}

	defer func() {
		if err := inFile.Close(); err != nil && closeErr == nil {
			closeErr = fmt.Errorf("failed to close input file: %w", err)
		}
	}()

	totalSize, err := getSize(inFile)
	if err != nil {
		return fmt.Errorf("failed to get file size: %w", err)
	}

	if offset > totalSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to offset: %w", err)
	}

	if limit == 0 || limit > totalSize-offset {
		limit = totalSize - offset
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() {
		if err := outFile.Close(); err != nil && closeErr == nil {
			closeErr = fmt.Errorf("failed to close output file: %w", err)
		}
	}()

	progress := make(chan int64)
	go func() {
		var lastProgress int64
		for p := range progress {
			fmt.Printf("Copied: %d bytes\n", p-lastProgress)
			lastProgress = p
		}
		fmt.Println("Copying completed!")
	}()

	err = copyProcess(inFile, outFile, limit, progress)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	return closeErr
}

func getSize(in io.Seeker) (int64, error) {
	size, err := in.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, fmt.Errorf("failed to determine size: %w", err)
	}

	_, err = in.Seek(0, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("failed to reset position: %w", err)
	}

	return size, nil
}

func copyProcess(in io.Reader, out io.Writer, limit int64, progress chan<- int64) error {
	bufferSize := 1024
	buf := make([]byte, bufferSize)

	var totalRead int64
	for totalRead < limit {
		bytesToRead := bufferSize
		if remaining := limit - totalRead; remaining < int64(bufferSize) {
			bytesToRead = int(remaining)
		}

		n, readErr := in.Read(buf[:bytesToRead])
		if n > 0 {
			totalRead += int64(n)

			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("failed to write: %w", writeErr)
			}

			if progress != nil {
				progress <- totalRead
			}
		}
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return fmt.Errorf("failed to read: %w", readErr)
		}
	}
	if progress != nil {
		close(progress)
	}

	if totalRead != limit {
		return errors.New("copied data size does not match expected limit")
	}

	return nil
}
