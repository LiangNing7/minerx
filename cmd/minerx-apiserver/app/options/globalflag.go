package options

import (
	"github.com/spf13/pflag"

	// ensure libs have a chance to globally register their flags.
	_ "k8s.io/apiserver/pkg/admission"
)

func AddCustomGlobalFlags(fs *pflag.FlagSet) {
	// Lookup flags in global flag set and re-register the values with our flagset.

	// Adds flags from k8s.io/apiserver/pkg/admission.
}
