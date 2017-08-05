package cplib

import (
	"fmt"
	"net/url"
	"strings"
)

//Args pass paramaters to cpanel API
type Args map[string]interface{}

func (a Args) Values(apiVersion string) url.Values {
	vals := url.Values{}
	for k, v := range a {
		if apiVersion == "1" {
			kv := strings.SplitN(k, "=", 2)
			if len(kv) == 1 {
				vals.Add(kv[0], "")
			} else if len(kv) == 2 {
				vals.Add(kv[0], kv[1])
			}
		} else {
			vals.Add(k, fmt.Sprintf("%v", v))
		}
	}
	return vals
}