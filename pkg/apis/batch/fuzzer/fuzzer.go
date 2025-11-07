package fuzzer

import (
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/randfill"

	"github.com/LiangNing7/minerx/pkg/apis/batch"
)

// Funcs returns the fuzzer functions for the batch api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(j *batch.Job, c randfill.Continue) {
			c.FillNoCustom(j) // fuzz self without calling this function again

			// // match defaulting
			// if len(j.Labels) == 0 {
			// 	j.Labels = j.Spec.Template.Labels
			// }
		},
		// func(j *batch.JobSpec, c randfill.Continue) {
		// 	c.FillNoCustom(j) // fuzz self without calling this function again
		// 	completions := int32(c.Rand.Int31())
		// 	parallelism := int32(c.Rand.Int31())
		// 	backoffLimit := int32(c.Rand.Int31())
		// 	j.Completions = &completions
		// 	j.Parallelism = &parallelism
		// 	j.BackoffLimit = &backoffLimit
		// 	j.ManualSelector = pointer.Bool(c.Bool())
		// 	mode := batch.NonIndexedCompletion
		// 	if c.Bool() {
		// 		mode = batch.IndexedCompletion
		// 		j.BackoffLimitPerIndex = pointer.Int32(c.Rand.Int31())
		// 		j.MaxFailedIndexes = pointer.Int32(c.Rand.Int31())
		// 	}
		// 	if c.Bool() {
		// 		j.BackoffLimit = pointer.Int32(math.MaxInt32)
		// 	}
		// 	j.CompletionMode = &mode
		// 	// We're fuzzing the internal JobSpec type, not the v1 type, so we don't
		// 	// need to fuzz the nil value.
		// 	j.Suspend = pointer.Bool(c.Bool())
		// 	podReplacementPolicy := batch.TerminatingOrFailed
		// 	if c.Bool() {
		// 		podReplacementPolicy = batch.Failed
		// 	}
		// 	j.PodReplacementPolicy = &podReplacementPolicy
		// 	if c.Bool() {
		// 		c.Fill(j.ManagedBy)
		// 	}
		// },
		func(sj *batch.CronJobSpec, c randfill.Continue) {
			c.FillNoCustom(sj)
			suspend := c.Bool()
			sj.Suspend = &suspend
			sds := int64(c.Uint64())
			sj.StartingDeadlineSeconds = &sds
			sj.Schedule = c.String(0)
			successfulJobsHistoryLimit := int32(c.Rand.Int31())
			sj.SuccessfulJobsHistoryLimit = &successfulJobsHistoryLimit
			failedJobsHistoryLimit := int32(c.Rand.Int31())
			sj.FailedJobsHistoryLimit = &failedJobsHistoryLimit
		},
		func(cp *batch.ConcurrencyPolicy, c randfill.Continue) {
			policies := []batch.ConcurrencyPolicy{batch.AllowConcurrent, batch.ForbidConcurrent, batch.ReplaceConcurrent}
			*cp = policies[c.Rand.Intn(len(policies))]
		},
		// func(p *batch.PodFailurePolicyOnPodConditionsPattern, c randfill.Continue) {
		// 	c.FillNoCustom(p)
		// 	if p.Status == "" {
		// 		p.Status = api.ConditionTrue
		// 	}
		// },
	}
}
