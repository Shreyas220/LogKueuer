package worker

import (
	"fmt"
	"os"
)

// Process data then ->  save statistics on node then send the location and

func RunWorker() {
	file, err := os.Open("location.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const chunkSize = 256 * 1024 // 256 KB
	buf := make([]byte, chunkSize)
	leftover := make([]byte, 0, chunkSize)

	go func() {
		for {
			bytesRead, err := file.Read(buf)
			if bytesRead > 0 {
				chunk := make([]byte, bytesRead)
				copy(chunk, buf[:bytesRead])
				validChunk, newLeftover := processChunk(chunk, leftover)
				leftover = newLeftover
				fmt.Println("go over valid chunk", validChunk)
			}
			if err != nil {
				break
			}
		}
	}()
	// ...
}

func processChunk(chunk, leftover []byte) (validChunk, newLeftover []byte) {
	firstNewline := -1
	lastNewline := -1
	// Find the first and last newline in the chunk.
	for i, b := range chunk {
		if b == '\n' {
			if firstNewline == -1 {
				firstNewline = i
			}
			lastNewline = i
		}
	}
	if firstNewline != -1 {
		validChunk = append(leftover, chunk[:lastNewline+1]...)
		newLeftover = make([]byte, len(chunk[lastNewline+1:]))
		copy(newLeftover, chunk[lastNewline+1:])
	} else {
		newLeftover = append(leftover, chunk...)
	}

	return validChunk, newLeftover
}

func processChunkData(chunk []byte) {
}
