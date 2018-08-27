package errors

import (
	"fmt"
	"net/http"
)

const (
	ReasonGroupBase    = "resource:"
	ReasonGroupStorage = "cluster." + ReasonGroupBase

	// bad request
	ErrorReasonBadPageStartOrLimit = ReasonGroupStorage + "BadPageStartOrLimit"
	ErrorReasonBadTenantOrUser     = ReasonGroupStorage + "BadTenantOrUser"
	ErrorReasonBadRequestBody      = ReasonGroupStorage + "BadRequestBody"
	ErrorReasonObjectAlreadyExist  = ReasonGroupStorage + "ObjectAlreadyExist"
	ErrorReasonObjectNotFound      = ReasonGroupStorage + "ObjectNotFound"
	ErrorReasonObjectConflict      = ReasonGroupStorage + "Conflict"
	ErrorReasonMissParameter       = ReasonGroupStorage + "MissParameter"
	ErrorReasonClusterNotFound     = ReasonGroupStorage + "ClusterNotFound"
	ErrorReasonClusterNotReady     = ReasonGroupStorage + "ClusterNotReady"
	// other error
	ErrorReasonAuthFailed          = ReasonGroupStorage + "AuthFailed"
	ErrorReasonInternalServerError = ReasonGroupStorage + "InternalServerError"
)

var (
	// inner use
	ErrVarKubeClientNil    = fmt.Errorf("kubernetes client is nil")
	ErrVarLogsStorageNil   = fmt.Errorf("log storage is nil")
	ErrVarClusterNil       = fmt.Errorf("cluster is nil")
	ErrVarBadConfig        = fmt.Errorf("bad config")
	ErrVarDuplicatedConfig = fmt.Errorf("duplicated config")
)

// bad request

func (fe *FormatError) SetErrorBadPageStartOrLimit(start, limit string) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("bad start or limit in query parameters start=%s, limit=%s", start, limit)
	fe.Reason = ErrorReasonBadPageStartOrLimit
	fe.HttpCode = http.StatusBadRequest
	return fe
}

func (fe *FormatError) SetErrorBadTenantOrUser(xTenant, xUser string) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("bad tenant or user in header parameters [%s:%s]", xTenant, xUser)
	fe.Reason = ErrorReasonBadTenantOrUser
	fe.HttpCode = http.StatusBadRequest
	return fe
}

func (fe *FormatError) SetErrorBadRequestBody(e error) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("parse request body failed, %v", e)
	fe.Reason = ErrorReasonBadRequestBody
	fe.HttpCode = http.StatusBadRequest
	fe.SetRawError(e)
	return fe
}

func (fe *FormatError) SetErrorObjectAlreadyExist(name string, e error) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("object %s already exist", name)
	fe.Reason = ErrorReasonObjectAlreadyExist
	fe.HttpCode = http.StatusBadRequest
	fe.SetRawError(e)
	return fe
}

func (fe *FormatError) SetErrorObjectNotFound(name string, e error) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("object %s not found", name)
	fe.Reason = ErrorReasonObjectNotFound
	fe.HttpCode = http.StatusNotFound
	fe.SetRawError(e)
	return fe
}

func (fe *FormatError) SetErrorObjectConflict(name string, e error) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("object %s conflict", name)
	fe.Reason = ErrorReasonObjectConflict
	fe.HttpCode = http.StatusConflict
	fe.SetRawError(e)
	return fe
}

func (fe *FormatError) SetErrorMissParameter(name string) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("parameter %s is empty", name)
	fe.Reason = ErrorReasonMissParameter
	fe.HttpCode = http.StatusBadRequest
	return fe
}

// cluster

func (fe *FormatError) SetErrorClusterNotReady(cluster, status string) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("cluster %s is in status %s, not ready for operation", cluster, status)
	fe.Reason = ErrorReasonClusterNotReady
	fe.HttpCode = http.StatusForbidden
	return fe
}

// other error

func (fe *FormatError) SetErrorAuthFailed(e error) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("parse auth info failed")
	fe.Reason = ErrorReasonAuthFailed
	fe.HttpCode = http.StatusForbidden
	fe.SetRawError(e)
	return fe
}

func (fe *FormatError) SetErrorInternalServerError(e error) *FormatError {
	fe.ApiError.Message = fmt.Sprintf("internal server error")
	fe.Reason = ErrorReasonInternalServerError
	fe.HttpCode = http.StatusInternalServerError
	fe.SetRawError(e)
	return fe
}
