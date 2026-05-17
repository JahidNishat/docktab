package commands

import (
	"fmt"
	"strings"
)

func ValidateAllowedValue(flagName, value string, allowed []string) error {
	for _, item := range allowed {
		if item == value {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid %s %q; valid values are: %s",
		flagName, value, strings.Join(allowed, ", "),
	)
}

func ValidateMutuallyExclusive(enabledA bool, nameA string, enabledB bool, nameB string) error {
	if enabledA && enabledB {
		return fmt.Errorf("%s and %s cannot be used together", nameA, nameB)
	}

	return nil
}
