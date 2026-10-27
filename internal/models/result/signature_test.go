package result

import "testing"

func TestGenerateSignature(t *testing.T) {
	tests := []struct {
		name   string
		ids    []int64
		suites []string
		params map[string]string
		want   string
	}{
		{name: "all parameters", ids: []int64{1, 2, 3}, suites: []string{"Suite1", "Suite2"},
			params: map[string]string{"Param1": "Value1", "Param2": "Value2"},
			want:   "1-2-3::suite1::suite2::{'param1':'value1'}::{'param2':'value2'}"},
		{name: "suite values processed", suites: []string{"  Suite With Spaces  ", "Suite\tWith\tTabs"},
			want: "suite_with_spaces::suite::with::tabs"},
		{name: "complex suite names",
			suites: []string{"test/fields.spec.js::Example: Fields.spec.js\tTest Cases With Field: Layer::Layer = Unit"},
			want:   "test/fields.spec.js::example:_fields.spec.js::test_cases_with_field:_layer::layer_=_unit"},
		{name: "nil ids", suites: []string{"suite1", "suite2"}, params: map[string]string{"param1": "value1"},
			want: "suite1::suite2::{'param1':'value1'}"},
		{name: "empty ids", ids: []int64{}, suites: []string{"suite1", "suite2"}, params: map[string]string{"param1": "value1"},
			want: "suite1::suite2::{'param1':'value1'}"},
		{name: "empty suites", ids: []int64{1, 2, 3}, suites: []string{}, params: map[string]string{"param1": "value1"},
			want: "1-2-3::{'param1':'value1'}"},
		{name: "empty params", ids: []int64{1, 2, 3}, suites: []string{"suite1"}, params: map[string]string{},
			want: "1-2-3::suite1"},
		{name: "only ids", ids: []int64{1, 2, 3}, want: "1-2-3"},
		{name: "only suites", suites: []string{"suite1", "suite2"}, want: "suite1::suite2"},
		{name: "only params", params: map[string]string{"param1": "value1", "param2": "value2"},
			want: "{'param1':'value1'}::{'param2':'value2'}"},
		{name: "backslashes normalized", suites: []string{`e2e\file.spec.ts`, "Suite1"}, want: "e2e/file.spec.ts::suite1"},
		{name: "slash direction irrelevant", suites: []string{`e2e\nested\file.spec.ts`}, want: "e2e/nested/file.spec.ts"},
		{name: "ids suites params",
			ids:    []int64{12, 15},
			suites: []string{"Search", "Smoke", "03-search-article", "03 · Search for an article"},
			params: map[string]string{"device": "iPhone 16 Pro - iOS 18.4"},
			want:   "12-15::search::smoke::03-search-article::03_·_search_for_an_article::{'device':'iphone 16 pro - ios 18.4'}"},
		{name: "suite split by :: and tab, backslash", suites: []string{" a\\b :: c\td "}, want: "a/b::c::d"},
		{name: "params sorted", params: map[string]string{"b": "2", "a": "1"}, want: "{'a':'1'}::{'b':'2'}"},
		{name: "empty", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSignature(tt.ids, tt.suites, tt.params); got != tt.want {
				t.Errorf("GenerateSignature() = %q, want %q", got, tt.want)
			}
		})
	}
}
