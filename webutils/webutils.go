/*
Package webutils
Library: msgolib
Some handy functions for HTTP/HTML stuff. Early stages of development.
Offered up under GPL 3.0 but absolutely not guaranteed fit for use.
This is code created by an amateur dilettante, so use at your own risk.
Github: https://github.com/mspeculatrix
Blog: https://mansfield-devine.com/speculatrix/
*/

package webutils

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseParams - takes a slice of strings, which are assumed
// to be in the format 'key=value', and returns a map.
// This was originally written for my own query string parsing
// method, but I'm not using that now.
// Still, might be handy to have around.
func ParseParams(items []string) (params map[string]string, err error) {
	if len(items) == 0 {
		return params, fmt.Errorf("empty list")
	}
	for item := range items {
		itemArr := strings.Split(items[item], "=")
		switch len(itemArr) {
		case 1:
			// assume a key has been passed but no value
			params[itemArr[0]] = ""
		case 2:
			// this is the expected condition
			params[itemArr[0]] = itemArr[1]
		default:
			// more than 2 items, so we'll assume that
			// the value contains one or more equals sign(s). Make the
			// value a Join of the itemArr items minus the first one
			// (which is the key)
			params[itemArr[0]] = strings.Join(itemArr[1:], "=")
		}
	}
	return params, err
}

// SimpleParseQuery - uses the builtin url.ParseQuery() method to
// parse the query string, because that might be cleverer than anything
// I could do. But that returns a type Values, which is a map with
// string keys and values consisting of string slices.
// I need something simpler - a map with both keys and values as strings.
func SimpleParseQuery(query string) (params map[string]string, err error) {
	if len(query) < 3 {
		// This function is designed to parse queries, so at the very least
		// the string passed to it should contain a slash, a '?' and at
		// least one character indicating a passed query. That's 3 chars.
		// Anything less means it can't possibly be a valid query.
		// Might want a regex here?
		return params, fmt.Errorf("query too short")
	}
	var queryStr string // will be the string we ultimately parse
	// Need to get rid of everything before and including query character
	parts := strings.Split(query, "?")
	switch len(parts) {
	case 0:
		// can't see how this could happen because of the length check
		// above, but just for the sake of completeness
		return params, fmt.Errorf("weird shit happened")
	case 1:
		// This must be an error. Only one part means it came before
		// the '?', or there wasn't a '?'.
		return params, fmt.Errorf("no query")
	case 2: // this is what we're expecting
		queryStr = parts[1]
	default: // more than 2 items
		// Maybe something in the query contained a question mark.
		// Let's stitch back together the query, except for the first item
		// (the bit before the first '?').
		queryStr = strings.Join(parts[1:], "?")
	}
	// Now perform the parsing
	items, err := url.ParseQuery(queryStr)
	if err != nil {
		return params, fmt.Errorf("simpleparsequery : %v", err)
	}
	// The values in the resulting 'items' map are string slices and
	// we want simple strings. In most cases, each slice probably has
	// only one item, but we need to be sure.
	for k, v := range items {
		if len(v) > 1 {
			// The value slice has more than one item, all
			// relating to the same key. Convert to a single string
			// with spaces between values and all enclosed in
			// square brackets.
			params[k] = "[" + strings.Join(v, " ") + "]"
		} else {
			// Only one item in the slice, so use this as our value.
			params[k] = v[0]
		}
	}
	return params, err
}
