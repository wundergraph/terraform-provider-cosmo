package utils

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ReconcileSchema(configured types.String, serverSchema string) types.String {
	if serverSchema == "" {
		return configured
	}
	if strings.TrimSpace(configured.ValueString()) == strings.TrimSpace(serverSchema) {
		return configured
	}
	return types.StringValue(serverSchema)
}

func StringValueOrNil(s types.String) *string {
	if s.IsNull() {
		return nil
	}
	value := s.ValueString()
	return &value
}

func ConvertHeadersToStringList(headersList types.List) []string {
	var headers []string
	for _, header := range headersList.Elements() {
		if headerStr, ok := header.(types.String); ok {
			headers = append(headers, headerStr.ValueString())
		}
	}
	return headers
}

func GetValueOrDefault[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

// NamespaceOrDefault falls back to the "default" namespace. Data source schemas cannot declare defaults.
func NamespaceOrDefault(namespace types.String) types.String {
	if namespace.IsNull() || namespace.ValueString() == "" {
		return types.StringValue("default")
	}
	return namespace
}
