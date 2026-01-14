package errwrap

import "google.golang.org/grpc/codes"

const (
	gRPC_OK                 = "OK"
	gRPC_Canceled           = "CANCELLED"
	gRPC_Unknown            = "UNKNOWN"
	gRPC_InvalidArgument    = "INVALID_ARGUMENT"
	gRPC_DeadlineExceeded   = "DEADLINE_EXCEEDED"
	gRPC_NotFound           = "NOT_FOUND"
	gRPC_AlreadyExists      = "ALREADY_EXISTS"
	gRPC_PermissionDenied   = "PERMISSION_DENIED"
	gRPC_ResourceExhausted  = "RESOURCE_EXHAUSTED"
	gRPC_FailedPrecondition = "FAILED_PRECONDITION"
	gRPC_Aborted            = "ABORTED"
	gRPC_OutOfRange         = "OUT_OF_RANGE"
	gRPC_Unimplemented      = "UNIMPLEMENTED"
	gRPC_Internal           = "INTERNAL"
	gRPC_Unavailable        = "UNAVAILABLE"
	gRPC_DataLoss           = "DATA_LOSS"
	gRPC_Unauthenticated    = "UNAUTHENTICATED"
)

// gRPCStatusText текстовое описание gRPC статуса
func gRPCStatusText(code int) string {
	switch code {
	case int(codes.OK):
		return gRPC_OK
	case int(codes.Canceled):
		return gRPC_Canceled
	case int(codes.Unknown):
		return gRPC_Unknown
	case int(codes.InvalidArgument):
		return gRPC_InvalidArgument
	case int(codes.DeadlineExceeded):
		return gRPC_DeadlineExceeded
	case int(codes.NotFound):
		return gRPC_NotFound
	case int(codes.AlreadyExists):
		return gRPC_AlreadyExists
	case int(codes.PermissionDenied):
		return gRPC_PermissionDenied
	case int(codes.ResourceExhausted):
		return gRPC_ResourceExhausted
	case int(codes.FailedPrecondition):
		return gRPC_FailedPrecondition
	case int(codes.Aborted):
		return gRPC_Aborted
	case int(codes.OutOfRange):
		return gRPC_OutOfRange
	case int(codes.Unimplemented):
		return gRPC_Unimplemented
	case int(codes.Internal):
		return gRPC_Internal
	case int(codes.Unavailable):
		return gRPC_Unavailable
	case int(codes.DataLoss):
		return gRPC_DataLoss
	case int(codes.Unauthenticated):
		return gRPC_Unauthenticated
	default:
		return ""
	}
}
