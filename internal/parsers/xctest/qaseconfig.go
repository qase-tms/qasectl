package xctest

import (
	"slices"
	"strings"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

// merge folds a later Qase config of the same test into c: IDs, fields and params add up,
// a later title or suite wins.
func (c *QaseConfig) merge(o QaseConfig) {
	ids := o.IDs
	if o.CaseID != nil {
		ids = append([]int64{*o.CaseID}, ids...)
	}
	for _, id := range ids {
		if !slices.Contains(c.IDs, id) {
			c.IDs = append(c.IDs, id)
		}
	}

	if o.Title != "" {
		c.Title = o.Title
	}
	if o.Suite != "" {
		c.Suite = o.Suite
	}

	c.Fields = mergeMap(c.Fields, o.Fields)
	c.Params = mergeMap(c.Params, o.Params)
	c.Ignore = c.Ignore || o.Ignore
	c.Muted = c.Muted || o.Muted
}

func (c QaseConfig) apply(r *models.Result) {
	if len(c.IDs) > 0 {
		// API v1 takes a single case ID, so the first one is kept for it.
		id := c.IDs[0]
		r.TestOpsID = &id
	}
	if len(c.IDs) > 1 {
		ids := slices.Clone(c.IDs)
		r.TestOpsIDs = &ids
	}

	if c.Title != "" {
		r.Title = c.Title
	}

	if suites := parseSuites(c.Suite); len(suites) > 0 {
		data := make([]models.SuiteData, 0, len(suites))
		for _, s := range suites {
			data = append(data, models.SuiteData{Title: s})
		}
		r.Relations.Suite.Data = data
	}

	r.Fields = mergeMap(r.Fields, c.Fields)
	r.Params = mergeMap(r.Params, c.Params)

	if c.Muted {
		r.Muted = true
		r.Fields = mergeMap(r.Fields, map[string]string{"muted": "true"})
	}
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

func mergeMap(dst, src map[string]string) map[string]string {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = make(map[string]string, len(src))
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
