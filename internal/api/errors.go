package api

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/wundergraph/cosmo/connect-go/gen/proto/wg/cosmo/common"
)

var (
	ErrUnknown                   = errors.New("ErrUnknown")
	ErrGeneral                   = errors.New("ErrGeneral")
	ErrNotFound                  = errors.New("ErrNotFound")
	ErrLimitReached              = errors.New("ErrLimitReached")
	ErrInvalidLabels             = errors.New("ErrInvalidLabels")
	ErrSubgraphCompositionFailed = errors.New("ErrSubgraphCompositionFailed")
	ErrEmptyMsg                  = fmt.Errorf("ErrEmptyMsg")
	ErrContractCompositionFailed = fmt.Errorf("ErrContractCompositionFailed")
	ErrInvalidSubgraphSchema     = fmt.Errorf("ErrInvalidSubgraphSchema")
)

const (
	ContractCompositionFailedReason = "A contract can only be created if its respective source graph has composed successfully"
)

func IsNotFoundError(err *ApiError) bool {
	return errors.Is(err.Err, ErrNotFound)
}

func IsSubgraphCompositionFailedError(err *ApiError) bool {
	return errors.Is(err.Err, ErrSubgraphCompositionFailed)
}

func IsInvalidSubgraphSchemaError(err *ApiError) bool {
	return errors.Is(err.Err, ErrInvalidSubgraphSchema)
}

func IsContractCompositionFailedError(err *ApiError) bool {
	return errors.Is(err.Err, ErrContractCompositionFailed)
}

type ApiError struct {
	Err    error
	Reason string
	Status common.EnumStatusCode
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s: %s (status: %s)", e.Err.Error(), e.Reason, e.Status.String())
}

func NewApiErrorWithErr(statusCode common.EnumStatusCode, reason string, err error) *ApiError {
	return &ApiError{Err: err, Reason: reason, Status: statusCode}
}

func handleErrorCodes(statusCode common.EnumStatusCode, reason string) *ApiError {
	if strings.Contains(reason, ContractCompositionFailedReason) {
		return &ApiError{Err: ErrContractCompositionFailed, Reason: reason, Status: statusCode}
	}

	switch statusCode {
	case common.EnumStatusCode_OK:
		return nil
	case common.EnumStatusCode_ERR_SUBGRAPH_COMPOSITION_FAILED:
		return &ApiError{Err: ErrSubgraphCompositionFailed, Reason: reason, Status: statusCode}
	case common.EnumStatusCode_ERR_NOT_FOUND:
		return &ApiError{Err: ErrNotFound, Reason: reason, Status: statusCode}
	case common.EnumStatusCode_ERR_INVALID_SUBGRAPH_SCHEMA:
		return &ApiError{Err: ErrInvalidSubgraphSchema, Reason: reason, Status: statusCode}
	case common.EnumStatusCode_ERR:
		return &ApiError{Err: ErrGeneral, Reason: reason, Status: statusCode}
	case common.EnumStatusCode_ERR_LIMIT_REACHED:
		return &ApiError{Err: ErrLimitReached, Reason: reason, Status: statusCode}
	case common.EnumStatusCode_ERR_INVALID_LABELS:
		return &ApiError{Err: ErrInvalidLabels, Reason: reason, Status: statusCode}
	default:
		return &ApiError{Err: ErrUnknown, Reason: reason, Status: statusCode}
	}
}
