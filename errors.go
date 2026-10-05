package zca

import "errors"

// ZaloAPIError is returned when Zalo answers with a non-zero error code or a request fails.
type ZaloAPIError struct {
	Message string
	Code    int // 0 when Zalo did not provide one
}

func (e *ZaloAPIError) Error() string { return e.Message }

func newError(msg string) error { return &ZaloAPIError{Message: msg} }

var (
	ErrLoginQRAborted             = errors.New("operation aborted")
	ErrLoginQRDeclined            = errors.New("login QR request declined")
	ErrMissingImageMetadataGetter = errors.New("missing ImageMetadataGetter, please provide it in Options")
)
