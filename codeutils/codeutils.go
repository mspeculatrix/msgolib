/*
Package codeutils
Library: msgolib
Various handy functions meant primarily for development.
Offered up under GPL 3.0 but absolutely not guaranteed fit for use.
This is code created by an amateur dilettante, so use at your own risk.
Github: https://github.com/mspeculatrix
Blog: https://mansfield-devine.com/speculatrix/
*/

package codeutils

import (
	"fmt"
	"strconv"
)

// PrintStringMap - prints out a map where both keys and
// values are strings.
func PrintStringMap(data *map[string]string) {
	longestKey := 0
	for k := range *data {
		if len(k) > longestKey {
			longestKey = len(k)
		}
	}
	format := "%" + strconv.Itoa(longestKey) + "s : %s\n"
	for k, v := range *data {
		fmt.Printf(format, k, v)
	}
}
