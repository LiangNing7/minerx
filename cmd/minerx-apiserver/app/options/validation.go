package options

import "fmt"

// Validate checks ServerRunOptions and return a slice of found errs.
func (o CompletedOptions) Validate() []error {
	errs := []error{}
	errs = append(errs, o.CompletedOptions.Validate()...)
	// errs = append(errs, s.CloudProvider.Validate()...)

	if o.MasterCount <= 0 {
		errs = append(errs, fmt.Errorf("--apiserver-count should be a positive number, but value '%d' provided", o.MasterCount))
	}

	return errs
}
