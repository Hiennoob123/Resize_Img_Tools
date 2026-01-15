package main

import "fmt"

type Job_pool struct {
	job_pending *Queue[*Job]
	job_finished *Queue[*Job]
	job_failed *Queue[*Job]
	total int
}

func NewJob_pool(paths []string) *Job_pool {
	job_pool := new(Job_pool)
	job_pool.total = len(paths)
	jobs := make([]*Job, len(paths))
	for id, path := range paths {
		jobs = append(jobs, NewJob(path, id))
	}
	job_pool.job_pending = newQueue(jobs)
	job_pool.job_finished = newQueue(make([]*Job, len(paths)))
	job_pool.job_failed = newQueue(make([]*Job, len(paths)))
	go func() {
		for job_pool.job_pending.length() != 0 {
			job_pool.start_newjob()
		}
	}()
	return job_pool
}

func (job_pool *Job_pool) start_newjob() {
	job, err := job_pool.job_pending.dequeue()
	if err != nil {
		fmt.Println("Dequeuing Job failed: ", err)
		return
	}
	err = job.run()
	if err != nil {
		fmt.Println("Job %d failed: %s", job.id, err)
		job_pool.job_failed.enqueue(job)
	} else {
		fmt.Println("Job %d success", job.id)
		job_pool.job_finished.enqueue(job)
	}
}

func (job_pool *Job_pool) stats() ([]int, []int, []int) {
	finished := make([]int, job_pool.job_finished.length())
	for _, job := range job_pool.job_finished.items {
		finished = append(finished, job.id)
	}
	failed := make([]int, job_pool.job_failed.length())
	for _, job := range job_pool.job_failed.items {
		failed = append(failed, job.id)
	}
	pending := make([]int, job_pool.job_pending.length())
	for _, job := range job_pool.job_pending.items {
		pending = append(pending, job.id)
	}
	return pending, finished, failed
}

func (job_pool *Job_pool) wait_empty() {
	for job_pool.job_finished.length() + job_pool.job_failed.length() < job_pool.total {}
}
