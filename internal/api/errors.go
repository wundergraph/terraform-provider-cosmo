package api

import (
	"fmt"

	"github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
)

var (
	ErrEmptyMsg                  = fmt.Errorf("empty message")
	ErrNotFound                  = fmt.Errorf("resource not found")
	ErrSubgraphCompositionFailed = fmt.Errorf("subgraph composition failed")
)

func IsNotFoundError(err error) bool {
	return err == ErrNotFound
}

func IsSubgraphCompositionFailedError(err error) bool {
	return err == ErrSubgraphCompositionFailed
}

func handleErrorCodes(statusCode common.EnumStatusCode) error {
	switch statusCode {
	case common.EnumStatusCode_OK:
		return nil
	case common.EnumStatusCode_ERR_SUBGRAPH_COMPOSITION_FAILED:
		return ErrSubgraphCompositionFailed
	case common.EnumStatusCode_ERR_NOT_FOUND:
		return ErrNotFound
	default:
		return fmt.Errorf("failed to create resource: %s", statusCode.String())
	}
}
