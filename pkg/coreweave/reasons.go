package coreweave

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var hiddenTaints map[string]bool

type Reason struct {
	Reason   string
	Priority int
	Count    int
}

type Reasons []Reason

// init loads the environment variables into memory
func init() {
	hiddenTaints = make(map[string]bool)
	taintsEnv := os.Getenv("CW_HIDDEN_TAINTS")
	taints := strings.Split(taintsEnv, ",")
	for _, v := range taints {
		hiddenTaints[v] = true
	}
}

func (r Reasons) Len() int {
	return len(r)
}
func (r Reasons) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r Reasons) Less(i, j int) bool {
	return r[i].Priority < r[j].Priority
}

// HiddenTaint accepts a namespace and taint key to return if reason message should be visible
func HiddenTaint(namespace string, key string) bool {
	if !strings.Contains(namespace, "tenant") {
		return false
	}
	return hiddenTaints[key]
}

// AsStrings on Reasons will return a priority sorted list of reasons
func (r Reasons) AsStrings() []string {
	sort.Sort(r)
	var reasons []string
	for _, v := range r {
		reasons = append(reasons, fmt.Sprintf("%v %v", v.Count, v.Reason))
	}
	return reasons
}

// reasonPriority accepts a string to give the error priority 0-100
func reasonPriority(reason string) int {
	if strings.Contains(reason, "volume") {
		return 1
	} else if strings.Contains(reason, "Insufficient memory") {
		return 5
	} else if strings.Contains(reason, "Insufficient cpu") {
		return 6
	} else if strings.Contains(reason, "node affinity") {
		return 10
	} else if strings.Contains(reason, "taints") {
		return 20
	}
	return 99
}

// SanitizeTenantReasons accepts a map of reasons to return a complete list of Reasons used in this package
func SanitizeTenantReasons(reasons map[string]int) Reasons {
	var newReasons Reasons
	for k, v := range reasons {
		newReasons = append(newReasons, Reason{
			Reason:   k,
			Priority: reasonPriority(k),
			Count:    v,
		})
	}
	return newReasons
}
