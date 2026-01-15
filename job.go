package main 

import (
	"os"
	"github.com/davidbyttow/govips/v2/vips"
	"fmt"
	"path/filepath"
)
type RES int
const (
	p720 RES = iota
	p1080
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
	err error
}

func NewJob(path string, id int) *Job {
	newjob := new(Job)
	newjob.path = path
	newjob.id = id 
	newjob.status = PENDING
	newjob.cnt = 0
	return newjob
}

func (job  Job) run() error {
	job.start()
	err := job.execute()
	if err == nil {
		job.finish()
	}	else {
		job.failed()
	}
	return err
}

func (job Job) start() {

	fmt.Println("Job", job.id, "run")
	job.status = RUNNING
	job.cnt += 1;

} 


func (job Job) finish() {
	job.status = FINISHED
}

func (job Job) failed() {
	job.status = FAILED
}

// Decode and resize the Image
func (job Job) execute() error {
	img, err := vips.NewImageFromFile(job.path)
	if (err != nil) {
		return err
	}

	// Use the copy to resize
	imgcopy, err := img.Copy()
	if err != nil {
		return err
	}



	switch Res {
	case p720:
		err = imgcopy.Thumbnail(1280, 720, vips.InterestingCentre)
		if err != nil {
			return err
		}
	case p1080:
		err = imgcopy.Thumbnail(1920, 1080, vips.InterestingCentre)
		if err != nil {
			return err
		}
	}

	buf, _, err := imgcopy.Export(vips.NewDefaultJPEGExportParams())

	if err != nil {
		return err
	}

	dst := dst_path(job.path)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		fmt.Println("Failed to create directory")
		return err
	}
	
	err = os.WriteFile(dst, buf, 0644)

	return err

}

