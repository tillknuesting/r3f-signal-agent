package entity

import (
	"time"
)

type CollectionJob struct {
	id          string
	profile     string
	sourceIDs   []string
	status      CollectionStatus
	startedAt   time.Time
	completedAt time.Time
	itemsCount  int
	errors      []string
}

type CollectionStatus string

const (
	CollectionStatusPending   CollectionStatus = "pending"
	CollectionStatusRunning   CollectionStatus = "running"
	CollectionStatusCompleted CollectionStatus = "completed"
	CollectionStatusFailed    CollectionStatus = "failed"
)

func NewCollectionJob(id, profile string, sourceIDs []string) *CollectionJob {
	return &CollectionJob{
		id:        id,
		profile:   profile,
		sourceIDs: sourceIDs,
		status:    CollectionStatusPending,
		errors:    []string{},
	}
}

func (j *CollectionJob) ID() string               { return j.id }
func (j *CollectionJob) Profile() string          { return j.profile }
func (j *CollectionJob) SourceIDs() []string      { return j.sourceIDs }
func (j *CollectionJob) Status() CollectionStatus { return j.status }
func (j *CollectionJob) StartedAt() time.Time     { return j.startedAt }
func (j *CollectionJob) CompletedAt() time.Time   { return j.completedAt }
func (j *CollectionJob) ItemsCount() int          { return j.itemsCount }
func (j *CollectionJob) Errors() []string         { return j.errors }

func (j *CollectionJob) Start() {
	j.status = CollectionStatusRunning
	j.startedAt = time.Now()
}

func (j *CollectionJob) Complete(count int) {
	j.status = CollectionStatusCompleted
	j.completedAt = time.Now()
	j.itemsCount = count
}

func (j *CollectionJob) Fail(err string) {
	j.status = CollectionStatusFailed
	j.completedAt = time.Now()
	j.errors = append(j.errors, err)
}

func (j *CollectionJob) AddError(err string) {
	j.errors = append(j.errors, err)
}

type CollectionResult struct {
	JobID    string
	Trends   []*Trend
	Errors   []string
	Duration time.Duration
}

func (j *CollectionJob) ToDTO() *CollectionJobDTO {
	return &CollectionJobDTO{
		ID:          j.id,
		Profile:     j.profile,
		SourceIDs:   j.sourceIDs,
		Status:      string(j.status),
		StartedAt:   j.startedAt,
		CompletedAt: j.completedAt,
		ItemsCount:  j.itemsCount,
		Errors:      j.errors,
	}
}

type CollectionJobDTO struct {
	ID          string    `json:"id"`
	Profile     string    `json:"profile"`
	SourceIDs   []string  `json:"source_ids"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	ItemsCount  int       `json:"items_count"`
	Errors      []string  `json:"errors"`
}
