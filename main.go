package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Damnever/goqueue"
)

var Wg = &sync.WaitGroup{}
var maxConcurrency = runtime.NumCPU() * 2
var errored = goqueue.New(0)
var complete = goqueue.New(0)
var queue = goqueue.New(0)

func bootstrap(host []string) (err error) {
	hostname := host[0]
	domain := host[1]
	chefEnv := host[2]
	fqdn := strings.Join([]string{hostname, domain}, ".")
	runlist := host[3]
	superuser_name := os.Getenv("SUPERUSER_NAME")
	superuser_pw := os.Getenv("SUPERUSER_PW")
	//sudo_value := os.Getenv("USE_SUDO")
	cmd := strings.Join([]string{"knife bootstrap ", fqdn, " -N ", hostname, " -E ", chefEnv, " --sudo", " --ssh-user ", superuser_name, " --ssh-password ", superuser_pw, " -r ", runlist}, "")
	fmt.Println("bootstrap command: ", cmd)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		filename := strings.Join([]string{"./logs/", hostname, ".txt"}, "")
		ioutil.WriteFile(filename, out, 0644)
	}
	return err
}

func worker(queue *goqueue.Queue) {
	for !queue.IsEmpty() {
		val, err := queue.Get(2)
		item := val.([]string)
		if err != nil {
			fmt.Println("Unexpect Error: %v\n", err)
		}
		bootstrap(item)
		fmt.Println("finished bootstrapping")
		if err != nil {
			errored.PutNoWait(val)
		} else {
			complete.PutNoWait(val)
		}
	}
	defer Wg.Done()
}

func main() {
	os.Mkdir("./logs", 0777)
	//in_progress := goqueue.New(0)
	//badauth := goqueue.New(0)
	//timeout := goqueue.New(0)
	//baddns := goqueue.New(0)

	// Read in the csv and populate queue for workers
	file, err := os.Open("./sample.tsv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '	'
	// read all records into memory
	result, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Queue all records
	for i := range result {
		record := result[i]
		fmt.Println("Queueing:", record)
		queue.PutNoWait(record)
	}
	// Start worker pool
	for i := 0; i < maxConcurrency && !queue.IsEmpty(); i++ {
		Wg.Add(1)
		go worker(queue)
		// Sleep 50 Milliseconds to give worker time to start
		time.Sleep(50 * time.Millisecond)
	}
	Wg.Wait()
}
