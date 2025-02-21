package models

import "encoding/json"

type JobName string
type JobType string
type JobResponse Job

type Job struct {
	// Type of job
	Type JobType `json:"type"`
	// The data you are packaging with the job
	Data any `json:"data"`
}

// Takes in raw job data and a reference struct for it's data
func GetJobFromData[T any](j []byte, data T) (*Job, T, error) {
	job := &Job{}
	err := json.Unmarshal(j, job)

	if err != nil {
		return nil, data, err
	}

	dat, err := json.Marshal(job.Data)

	if err != nil {
		return nil, data, err
	}

	err = json.Unmarshal(dat, &data)

	job.Data = data

	return job, data, err
}

func GetDataFromStruct[T any](job Job, data T) (T, error) {
	dat, err := json.Marshal(job.Data)

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(dat, &data)

	job.Data = data

	return data, err
}
