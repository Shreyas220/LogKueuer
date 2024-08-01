package main

import "LogKueueEngine/worker"

// get the data into local storage
func main() {
	worker.RunWorker()
}

// divide it into chucks and create kueue jobs using it

// schedule job in kueue
