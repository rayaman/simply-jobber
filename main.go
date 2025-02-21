package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rayaman/simply-jobber/pkg/api/queues"
	"github.com/rayaman/simply-jobber/pkg/models"
)

type Data struct {
	Message string `json:"message"`
}

type Data2 struct {
	Msg string `json:"message"`
}

var (
	Job1 models.JobType = "Test 1"
	Job2 models.JobType = "Test 2"
)

func main() {
	var queue models.Queue = queues.NewSimple(0)

	myJob := models.Job{
		Type: Job1,
		Data: Data{
			Message: "Job message :D",
		},
	}

	myJob2 := models.Job{
		Type: Job2,
		Data: Data2{
			Msg: "Job message :D",
		},
	}

	err := queue.AddHandle(Job1, func(j models.Job) (models.JobResponse, error) {
		data, err := models.GetDataFromStruct(j, Data{})

		if err != nil {
			return models.JobResponse{}, fmt.Errorf("cannot get data from job")
		}

		fmt.Printf("Got data 1 %v\n", data.Message)
		time.Sleep(time.Second)

		return models.JobResponse{}, nil
	})

	if err != nil {
		panic(err)
	}

	err = queue.AddHandle(Job2, func(j models.Job) (models.JobResponse, error) {
		data, err := models.GetDataFromStruct(j, Data2{})

		if err != nil {
			return models.JobResponse{}, fmt.Errorf("cannot get data from job")
		}

		fmt.Printf("Got data 2 %v\n", data.Msg)
		time.Sleep(time.Second * 2)

		return models.JobResponse{}, nil
	})

	if err != nil {
		panic(err)
	}

	// queue.SetProcessor(func(j models.Job) (models.JobResponse, error) {
	// 	switch j.Type {
	// 	case "test":
	// 		data, err := models.GetDataFromStruct(j, Data{})
	// 		if err != nil {
	// 			return models.JobResponse{}, fmt.Errorf("cannot get data from job")
	// 		}
	// 		include := "N/A"
	// 		if data.Message == "Job message :D" {
	// 			include = "Special+"
	// 		}
	// 		return models.JobResponse{
	// 			Type: "test",
	// 			Data: Data{
	// 				Message: "Completed job :D " + include,
	// 			},
	// 		}, nil
	// 	default:
	// 		return models.JobResponse{}, fmt.Errorf("cannot process job")
	// 	}
	// })

	// queue.OnJobResponse(func(jr models.JobResponse, i int) {
	// 	data, err := models.GetDataFromStruct(models.Job(jr), Data{})
	// 	if err != nil {
	// 		return
	// 	}
	// 	fmt.Printf("JobID: %v Message: %v\n", i, data.Message)
	// })
	go func() {
		time.Sleep(time.Second * 10)
		queue.Scale(1)
		time.Sleep(time.Second * 5)
		queue.Scale(5)
		time.Sleep(time.Second * 5)
		queue.Scale(0)
		time.Sleep(time.Second * 5)
		queue.Scale(5)
		time.Sleep(time.Second * 5)
		os.Exit(0)
	}()

	for {
		time.Sleep(time.Second * 1)
		queue.Send(myJob)
		time.Sleep(time.Second * 1)
		queue.Send(myJob2)
	}
}
