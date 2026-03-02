package job

import (
	v1 "k8s.io/api/core/v1"

	batch "github.com/LiangNing7/minerx/pkg/apis/batch/v1beta1"
)

// FinishedCondition returns true if a job is finished as well as the condition type indicating that.
// Returns false and no condition type otherwise
func FinishedCondition(j *batch.Job) (bool, batch.JobConditionType) {
	for _, c := range j.Status.Conditions {
		if (c.Type == batch.JobComplete || c.Type == batch.JobFailed) && c.Status == v1.ConditionTrue {
			return true, c.Type
		}
	}
	return false, ""
}

// IsJobFinished checks whether the given Job has finished execution.
// It does not discriminate between successful and failed terminations.
func IsJobFinished(j *batch.Job) bool {
	isFinished, _ := FinishedCondition(j)
	return isFinished
}

// IsJobSucceeded returns whether a job has completed successfully.
func IsJobSucceeded(j *batch.Job) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == batch.JobComplete && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
