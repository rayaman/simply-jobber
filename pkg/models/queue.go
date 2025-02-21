package models

type Queue interface {
	Send(Job) (int, error)
	// Callback function that is triggered when a job has been processed. Can view the job response
	OnJobResponse(ResponseFunc)
	// Sets a function that will be triggered by all jobs not processed by a handle
	SetProcessor(ProcessFunc, ...string) error
	// Adds a handle for a certain job type. this func is only triggered by a piticular job type
	AddHandle(JobType, ProcessFunc) error
	// Allows you to scale the queue
	Scale(uint)
	// Number of queued jobs
	Len() int
}

type ProcessFunc func(Job) (JobResponse, error)
type ResponseFunc func(JobResponse, int)
