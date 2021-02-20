package cmd

import (
	"strings"
)

// CaseMap 忽略大小写key的map
type CaseMap map[string]CommandHelp

func (c *CaseMap) Find(key string) (CommandHelp, bool) {
	key = strings.TrimRight(key, ` `)
	uKey := strings.ToUpper(key)

	v, e := (*c)[uKey]
	return v, e
}
