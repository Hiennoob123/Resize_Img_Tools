package main 

import (
	"fmt"
	"github.com/nfnt/resize"
	"os"
	"image/jpeg"
	"image"
)

type STATUS int

const (
	PENDING STATUS = iota
	RUNNING
	FINISHED
	FAILED
)

type Job struct {
	id int 
	path string 
	status STATUS 
	cnt int
}

func NewJob(path string, id int) *Job {
	newjob := new(Job)
	newjob.path = path
	newjob.id = id 
	newjob.status = PENDING
	newjob.cnt = 0
	return newjob
}

func (job Job) start() {

	job.status = RUNNING
	job.cnt += 1;
	job.execute()

} 

func (job Job) execute() {

	
	file, err := os.Open(job.path)
	if err != nil {
		fmt.Println("Job", job.id, job.path, ": Error", err)
	}
	defer file.Close()
}

