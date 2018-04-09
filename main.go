package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Damnever/goqueue"
)

var Wg = &sync.WaitGroup{}
var maxConcurrency = runtime.NumCPU() * 2
var errored = goqueue.New(0)
var complete = goqueue.New(0)
var queue = goqueue.New(0)
var badDNS []Host
var timeoutHosts []Host
var authHosts []Host
var generalError []Host

type Report struct {
	DNS_Hosts     []Host `json:"dns_hosts"`
	Auth_Hosts    []Host `json:"auth_hosts"`
	Timeout_Hosts []Host `json:"timeout_hosts"`
}

type Host struct {
	Hostname string `json:"hostname"`
	Domain   string `json:"domain"`
	ChefEnv  string `json:"chefenv"`
	RunList  string `json:"runlist"`
}

func host_validate(hosts []Host) {
	for _, element := range hosts {
		if element.Hostname == " " {
			log.Fatal("Please ensure the following entry contains a hostname:", element)
		}
		if element.Domain == " " {
			log.Fatal("Please ensure the following entry contains a domain:", element)
		}
		if element.ChefEnv == " " {
			log.Fatal("Please ensure the following entry contains a chef environment:", element)
		}
		if element.RunList == " " {
			log.Fatal("Please ensure the following entry contains a run list:", element)
		}
	}
	return
}

func csv_to_hosts(csv_filename string) (hosts []Host) {
	file, err := os.Open(csv_filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '	'
	// read all records into memory
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		hosts = append(hosts, Host{
			Hostname: line[0],
			Domain:   line[1],
			ChefEnv:  line[2],
			RunList:  line[3]},
		)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	host_validate(hosts)
	return
}

func handle_bootstrap_error(out []byte, host Host, exit_code int) (bootstrap_success bool) {
	o := string(out)
	if strings.Contains(o, "Authentication failed") {
		authHosts = append(authHosts, host)
		return false
	}
	if strings.Contains(o, "ConnectionTimeout") {
		timeoutHosts = append(timeoutHosts, host)
		return false
	}
	if strings.Contains(o, "nodename nor servname provided") {
		badDNS = append(badDNS, host)
		return false
	}
	if exit_code != 0 {
		generalError = append(generalError, host)
	}
	return true
}

func run_command(cmd string) (out []byte, exit_code int) {
	c := exec.Command("sh", "-c", cmd)
	cmdOutput := &bytes.Buffer{}
	cmdErrorOutput := &bytes.Buffer{}
	c.Stdout = cmdOutput
	c.Stderr = cmdErrorOutput
	if err := c.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}
	if err := c.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exit_code = status.ExitStatus()
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	var output_array []byte
	output_array = append(output_array, cmdErrorOutput.Bytes()...)
	output_array = append(output_array, cmdOutput.Bytes()...)
	return output_array, exit_code
}

func bootstrap(host Host) {
	cmd := generate_command(host)
	cmd_out, exit_code := run_command(cmd)
	filename := strings.Join([]string{"./logs/", host.Hostname, ".txt"}, "")
	ioutil.WriteFile(filename, cmd_out, 0644)
	handle_bootstrap_error(cmd_out, host, exit_code)
	return
}

func error_report() (report Report) {

	for i := range badDNS {
		report.DNS_Hosts = append(report.DNS_Hosts, badDNS[i])
	}
	for i := range authHosts {
		report.Auth_Hosts = append(report.Auth_Hosts, authHosts[i])
	}
	for i := range timeoutHosts {
		report.Timeout_Hosts = append(report.Timeout_Hosts, timeoutHosts[i])
	}
	return report
}

func generate_command(host Host) (cmd string) {
	fqdn := strings.Join([]string{host.Hostname, host.Domain}, ".")
	superuser_name := os.Getenv("SUPERUSER_NAME")
	superuser_pw := os.Getenv("SUPERUSER_PW")
	//sudo_value := os.Getenv("USE_SUDO")
	cmd = strings.Join([]string{"knife bootstrap ", fqdn, " -N ", host.Hostname, " -E ", host.ChefEnv, " --sudo", " --ssh-user ", superuser_name, " --ssh-password ", superuser_pw, " -r ", host.RunList}, "")
	return
}
func worker(queue *goqueue.Queue) {
	for !queue.IsEmpty() {
		//Get queue with 2 second timeout
		val, err := queue.Get(2)
		item := val.(Host)
		if err != nil {
			fmt.Println("Unexpect Error: \n", err)
		}
		bootstrap(item)
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

	// Read in the csv and populate queue for workers
	var hosts []Host
	var csv_filename string

	flag.StringVar(&csv_filename, "file", "./sample.tsv", "file containing hosts to be bootstrapped")
	flag.Parse()
	hosts = csv_to_hosts(csv_filename)
	// Queue all records
	for i := range hosts {
		record := hosts[i]
		recordJson, _ := json.Marshal(record)
		fmt.Println("Queueing:", string(recordJson))
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
	if len(badDNS) > 0 || len(timeoutHosts) > 0 || len(authHosts) > 0 {
		report := error_report()
		r, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println("Error Report:")
		fmt.Printf("%s/n", r)
	}
}
