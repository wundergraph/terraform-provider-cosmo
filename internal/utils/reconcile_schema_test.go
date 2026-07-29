package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestReconcileSchema(t *testing.T) {
	testCases := []struct {
		name         string
		configured   types.String
		serverSchema string
		expected     types.String
	}{
		{
			name:         "keeps configured value when server only strips trailing whitespace",
			configured:   types.StringValue("type Query {\n  hello: String\n}\n\n"),
			serverSchema: "type Query {\n  hello: String\n}",
			expected:     types.StringValue("type Query {\n  hello: String\n}\n\n"),
		},
		{
			name:         "keeps configured value when server only strips surrounding whitespace",
			configured:   types.StringValue("\ntype Query {\n  hello: String\n}\n"),
			serverSchema: "type Query {\n  hello: String\n}",
			expected:     types.StringValue("\ntype Query {\n  hello: String\n}\n"),
		},
		{
			name:         "adopts server value when the document differs",
			configured:   types.StringValue("type Query {\n  hello: String\n}"),
			serverSchema: "type Query {\n  world: String\n}",
			expected:     types.StringValue("type Query {\n  world: String\n}"),
		},
		{
			name:         "adopts server value when configured is null (import)",
			configured:   types.StringNull(),
			serverSchema: "type Query {\n  hello: String\n}",
			expected:     types.StringValue("type Query {\n  hello: String\n}"),
		},
		{
			name:         "keeps configured value when server schema is empty",
			configured:   types.StringValue("type Query {\n  hello: String\n}"),
			serverSchema: "",
			expected:     types.StringValue("type Query {\n  hello: String\n}"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ReconcileSchema(tc.configured, tc.serverSchema)
			if !got.Equal(tc.expected) {
				t.Errorf("reconcileSchema() = %v, want %v", got, tc.expected)
			}
		})
	}
}
