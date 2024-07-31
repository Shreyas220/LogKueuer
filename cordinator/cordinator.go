package cordinator

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

var Jobs map[string]JobStatus

type JobStatus struct {
	lock         sync.Mutex
	totalNumber  int
	jobcompleted int
}

func RunCordinator() {
	js := JobStatus{}
	chunks := js.processChunks()
	for i, chunk := range chunks {
		jobName := fmt.Sprintf("process-file-job-%d-%s", i, "jobid")

		encodedChunk := base64.StdEncoding.EncodeToString(chunk)
		fmt.Println("chunk stored on node in somefile", encodedChunk)
		js.createJob(jobName)

	}
}

func createConfigMap(clientset *kubernetes.Clientset, namespace, name string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	return err
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

func (j *JobStatus) processChunks() [][]byte {
	const chunkSize = 100 * 1024 * 1024 // 100 mb files
	buf := make([]byte, chunkSize)
	leftover := make([]byte, 0, chunkSize)

	file, err := os.Open("k8s_audit_logs.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var chunks [][]byte
	counter := 0

	reading := true
	for reading {

		bytesRead, _ := file.Read(buf)
		if bytesRead > 0 {
			// buffer will be reused in the next iteration.
			chunk := make([]byte, bytesRead)
			copy(chunk, buf[:bytesRead])
			validChunk, newLeftover := processChunk(chunk, leftover)
			leftover = newLeftover
			chunks = append(chunks, validChunk)
			counter++
		} else {
			reading = false
		}

	}

	j.totalNumber = counter
	return chunks
}

func (j *JobStatus) createJob(jobname string) {

}

func (j *JobStatus) updateJobCompletion() {

}
