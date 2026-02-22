/*
Package randomhelper that has helper functions that doesn't belong to a particular destination!
*/
package randomhelper

import "log"

func CheckIfAllEnvValid(variables ...string) {
	for i, pass := range variables {
		if pass == "" {
			log.Fatalf("could not load variable at index: %d", i)
		}
	}
}
