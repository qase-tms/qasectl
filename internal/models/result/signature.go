package result

import (
	"sort"
	"strconv"
	"strings"
)

// GenerateSignature mirrors the signature built by the Qase reporters. Params
// are emitted sorted by key because Go maps are unordered.
func GenerateSignature(ids []int64, suites []string, params map[string]string) string {
	var parts []string

	if len(ids) > 0 {
		strs := make([]string, len(ids))
		for i, id := range ids {
			strs[i] = strconv.FormatInt(id, 10)
		}
		parts = append(parts, strings.Join(strs, "-"))
	}

	var processed []string
	for _, suite := range suites {
		suite = strings.ReplaceAll(strings.TrimSpace(suite), `\`, "/")
		for _, byColons := range strings.Split(suite, "::") {
			for _, byTab := range strings.Split(strings.TrimSpace(byColons), "\t") {
				s := strings.ToLower(strings.Join(strings.Fields(byTab), "_"))
				if s != "" {
					processed = append(processed, s)
				}
			}
		}
	}
	if len(processed) > 0 {
		parts = append(parts, strings.Join(processed, "::"))
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts = append(parts, "{'"+strings.ToLower(k)+"':'"+strings.ToLower(params[k])+"'}")
	}

	return strings.Join(parts, "::")
}
