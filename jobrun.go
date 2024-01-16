package goflow

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/philippgille/gokv"
)

type jobRun struct {
	ID        string    `json:"id"`
	JobName   string    `json:"job"`
	StartedAt string    `json:"submitted"`
	JobState  *jobState `json:"state"`
}

type jobRunIndex struct {
	JobRunIDs []string
}

func (j *Job) newJobRun() *jobRun {
	return &jobRun{
		ID:        uuid.New().String(),
		JobName:   j.Name,
		StartedAt: time.Now().UTC().Format(time.RFC3339Nano),
		JobState:  j.jobState}
}

// Persist a new jobrun.
func persistNewJobRun(store gokv.Store, jobrun *jobRun) error {
	key := jobrun.ID
	err := store.Set(key, jobrun)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Index the job runs
func indexJobRuns(store gokv.Store, jobrun *jobRun) error {
	index := jobRunIndex{}
	store.Get(jobrun.JobName, &index)

	// add the jobrun ID to the index
	index.JobRunIDs = append(index.JobRunIDs, jobrun.ID)
	return store.Set(jobrun.JobName, index)
}

// Read all the persisted jobruns for a given job.
func readJobRuns(store gokv.Store, jobName string) ([]*jobRun, error) {
	index := jobRunIndex{}
	store.Get(jobName, &index)

	jobRuns := make([]*jobRun, 0)
	for _, key := range index.JobRunIDs {
		value := jobRun{}
		store.Get(key, &value)
		jobRuns = append(jobRuns, &value)
	}

	return jobRuns, nil
}

// Sync the current jobstate to the persisted jobrun.
func updateJobState(store gokv.Store, jobrun *jobRun, jobstate *jobState) error {

	// Get the key
	key := jobrun.ID

	// Get the lock
	jobstate.TaskState.RLock()

	// Update the jobrun state
	jobrun.JobState = jobstate

	// Persist it
	err := store.Set(key, jobrun)

	// Release lock
	jobstate.TaskState.RUnlock()

	return err
}
