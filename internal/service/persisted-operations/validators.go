package persisted_operations

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = notJSONValidator{}

// notJSONValidator rejects operation contents that parse as JSON. wgc routes valid JSON to its
// manifest parsers (Apollo, Relay), a strategy this resource does not implement, so JSON contents
// would otherwise be silently hashed and registered as a single bogus operation.
type notJSONValidator struct{}

func (v notJSONValidator) Description(ctx context.Context) string {
	return "value must be a plain GraphQL document, not JSON"
}

func (v notJSONValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v notJSONValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if json.Valid([]byte(req.ConfigValue.ValueString())) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			ErrInvalidOperationContents,
			"operation contents must be a plain GraphQL document; JSON manifests (e.g. Apollo persisted query manifests, Relay query maps) are not supported by this resource yet. Extract the operation bodies or push them with `wgc operations push`.",
		)
	}
}
