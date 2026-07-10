package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

// NamespaceOrDefault resolves an optional namespace attribute to the "default"
// namespace. Data source schemas cannot declare defaults, so data sources
// resolve the namespace at read time with this helper.
func NamespaceOrDefault(namespace types.String) types.String {
	if namespace.IsNull() || namespace.ValueString() == "" {
		return types.StringValue("default")
	}
	return namespace
}
