package maestro

import (
	"log/slog"
	"strconv"
	"strings"
)

const (
	qaseIDKey     = "qase.id"
	qaseSuiteKey  = "qase.suite"
	qasePrefix    = "qase."
	tagsFieldName = "tags"
)

func parseQaseIDs(v string) []int64 {
	var ids []int64
	for _, part := range strings.Split(v, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		raw := part[strings.LastIndex(part, "-")+1:]
		id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			slog.Warn("failed to parse qase.id element", "value", part)
			continue
		}
		ids = append(ids, id)
	}
	return ids
}

func parseSuites(v string) []string {
	var parts []string
	for _, p := range strings.Split(v, ".") {
		if p = strings.TrimSpace(p); p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func fieldsFor(f flow) map[string]string {
	fields := make(map[string]string, len(f.Properties)+1)
	for k, v := range f.Properties {
		if strings.HasPrefix(k, qasePrefix) {
			continue
		}
		fields[k] = v
	}
	if len(f.Tags) > 0 {
		fields[tagsFieldName] = strings.Join(f.Tags, ", ")
	}
	return fields
}
