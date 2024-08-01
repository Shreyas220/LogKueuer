package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Process data then ->  save statistics on node then send the location and
var Aggregated_result map[string][]string

var Aggreate chan map[string][]string

func RunWorker() {
	Aggregated_result = make(map[string][]string)
	Aggreate = make(chan map[string][]string, 1000)

	file, err := os.Open("k8s_audit_logs.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const chunkSize = 256 * 1024 // 256 KB
	buf := make([]byte, chunkSize)
	leftover := make([]byte, 0, chunkSize)

	// go func() {
	for {
		bytesRead, err := file.Read(buf)
		if bytesRead > 0 {
			chunk := make([]byte, bytesRead)
			copy(chunk, buf[:bytesRead])
			validChunk, newLeftover := processChunk(chunk, leftover)
			leftover = newLeftover
			processChunkData(validChunk)
			// fmt.Println("go over valid chunk", validChunk)
		}
		if err != nil {
			break
		}
	}
	// }()
	for {
		select {
		// case <-Context().Done():
		// 	return nil
		case resp := <-Aggreate:
			fmt.Println("Aggregate", resp)
		default:
			time.Sleep(time.Millisecond * 10)

		}
	}
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
	scanner := bufio.NewScanner(strings.NewReader(string(chunk)))
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		var data map[string]interface{}
		err := json.Unmarshal(line, &data)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			continue
		}
		if user, ok := data["user"].(map[string]interface{}); ok {
			if username, ok := user["username"].(string); ok {
				fmt.Printf("Username: %s\n", username)
			} else {
				fmt.Println("Username field is missing or not a string")
			}
		} else {
			fmt.Println("User field is missing or not a map")
		}

		i++
		if i == 5 {
			break
		}
	}
}
