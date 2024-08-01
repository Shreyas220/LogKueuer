package main

import "logkueuer/worker"

// get the data into local storage
func main() {
	worker.RunWorker()
}

// divide it into chucks and create kueue jobs using it

// schedule job in kueue
