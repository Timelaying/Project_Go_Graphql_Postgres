package jobs

import "time"

type Status string

// StatusApplied represents the initial status when a candidate has submitted their job application.
// This status indicates that the application has been received and is in the system,
// but no further action has been taken yet.
const (
	StatusApplied   Status = "APPLIED"
	StatusInterview Status = "INTERVIEW"
	StatusOffer     Status = "OFFER"
	StatusRejected  Status = "REJECTED"
)

type Job struct {
	ID        string
	Company   string
	Role      string
	Link      *string
	Status    Status
	CreatedAt time.Time
}
