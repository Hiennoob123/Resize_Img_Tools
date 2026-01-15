package main

import (
	"fmt"
	"sync"
)

type Job_pool struct {
	job_pending *Queue[*Job]
	job_finished *Queue[*Job]
	job_failed *Queue[*Job]
	total int
	wg sync.WaitGroup
}

func NewJob_pool(paths []string) *Job_pool {
	job_pool := new(Job_pool)
	job_pool.total = len(paths)
	jobs := make([]*Job, 0)
	for id, path := range paths {
		jobs = append(jobs, NewJob(path, id))
	}
	job_pool.job_pending = newQueue(jobs)
	job_pool.job_finished = newQueue(make([]*Job, 0, len(paths)))
	job_pool.job_failed = newQueue(make([]*Job, 0, len(paths)))
	job_pool.wg.Add(len(paths))
	go func() {
		for _ = range len(paths) {
			go job_pool.start_newjob()
		}
	}()
	return job_pool
}

func (job_pool *Job_pool) start_newjob() {
	defer job_pool.wg.Done()
	job, err := job_pool.job_pending.dequeue()
	if err != nil {
		fmt.Printf("Job %d failed: Dequeuing Job failed: %s \n", job.id, err)
		return
	}
	err = job.run()
	if err != nil {
		fmt.Printf("Job %d %s failed: %s\n", job.id, job.path, err)
		job.err = err
		job_pool.job_failed.enqueue(job)
	} else {
		//fmt.Printf("Job %d success\n", job.id)
		job_pool.job_finished.enqueue(job)
	}
}

func (job_pool *Job_pool) stats() ([]*Job, []*Job, []*Job) {
	return job_pool.job_pending.items, job_pool.job_finished.items, job_pool.job_failed.items
}

func (job_pool *Job_pool) wait_empty() ([]*Job, []*Job, []*Job) {
	job_pool.wg.Wait()
	return job_pool.stats()
}
