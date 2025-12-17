package printers

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

// TranslateMicroTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func TranslateMicroTimestampSince(timestamp metav1.MicroTime) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

// TranslateTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func TranslateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

// TranslateTimestampUntil returns the elapsed time until timestamp in
// human-readable approximation.
func TranslateTimestampUntil(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Until(timestamp.Time))
}
