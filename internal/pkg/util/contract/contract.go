package contract

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/flect"
)

// CalculateCRDName generates a CRD name based on group and kind according to
// the naming conventions in the contract.
func CalculateCRDName(group, kind string) string {
	return fmt.Sprintf("%s.%s", flect.Pluralize(strings.ToLower(kind)), group)
}
