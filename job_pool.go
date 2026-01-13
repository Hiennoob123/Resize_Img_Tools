package main

type job_pool struct {
	cnt_pending int
	cnt_finished int
	cnt_failed int
	job_pending Queue[Job]
	job_finished Queue[Job]
	job_failed Queue[Job]
}
