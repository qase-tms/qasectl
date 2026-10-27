package integration

import (
	"testing"
)

type scenario struct {
	name   string   // subtest name, e.g. "junit/pytest"
	format string   // --format value
	dir    string   // testdata/<format>/<scenario>
	args   []string // extra flags
	golden string   // defaults to dir + "/expected.json"
}

var scenarios = []scenario{
	{name: "junit/pytest", format: "junit", dir: "testdata/junit/pytest"},
	{name: "junit/maven", format: "junit", dir: "testdata/junit/maven"},
	{name: "junit/jest", format: "junit", dir: "testdata/junit/jest"},
	{name: "allure/pytest", format: "allure", dir: "testdata/allure/pytest"},
	{name: "junit/maestro", format: "junit", dir: "testdata/junit/maestro"},
	{name: "maestro/run", format: "maestro", dir: "testdata/maestro/run"},
	{
		name: "maestro/run-steps-all", format: "maestro", dir: "testdata/maestro/run",
		args: []string{"--steps", "all"}, golden: "testdata/maestro/run/expected-steps-all.json",
	},
	{name: "maestro/shards", format: "maestro", dir: "testdata/maestro/shards"},
}

func TestUpload(t *testing.T) {
	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			url, rec := newFakeQase(t)
			runUpload(t, url, sc.format, sc.dir+"/input", sc.args...)
			golden := sc.golden
			if golden == "" {
				golden = sc.dir + "/expected.json"
			}
			assertGolden(t, golden, normalize(rec.Requests()))
		})
	}
}
