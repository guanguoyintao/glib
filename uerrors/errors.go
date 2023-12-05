// Package uerrors provides a way to return detailed information
// for all server and common error.
package uerrors

import (
	"context"
	"fmt"
	"git.umu.work/be/goframework/common"
	"github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/v2/errors"
	gerrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
)

var (
	// errDuplicateEntryCode 命中唯一索引
	errDuplicateEntryCode uint16 = 1062
)

const (
	COMMON_ERROR_CODE_PREFIX   int32 = 20
	UCS_ERROR_CODE_PREFIX      int32 = 21
	UCC_ERROR_CODE_PREFIX      int32 = 22
	ASR_POLY_ERROR_CODE_PREFIX int32 = 23
	USHOW_ERROR_CODE_PREFIX    int32 = 24
	UEB_ERROR_CODE_PREFIX      int32 = 25
	TMT_POLY_ERROR_CODE_PREFIX int32 = 26
)

// common
var (
	UErrorDynamicConfigTypeUnknown  = common.NewUmuError(2001, "incorrect dynamic configuration type")
	UErrorRepoRecordNotFound        = common.NewUmuError(2002, "record not found")
	UErrorRepoRecordConstraint      = common.NewUmuError(2003, "data business constraint error")
	UErrorPkgCounterNegativeNumber  = common.NewUmuError(2004, "business counter becomes negative")
	UErrorParameterError            = common.NewUmuError(2005, "internal parameter transfer error")
	UErrorNullValueError            = common.NewUmuError(2006, "null value error")
	UErrorTimeout                   = common.NewUmuError(2007, "internal timeout error")
	UErrorSystemError               = common.NewUmuError(2008, "internal system error")
	UErrorPaginationNotFoundOrder   = common.NewUmuError(2009, "pagination not found order condition")
	UErrorKafukaNotExitError        = common.NewUmuError(2010, "not exist kafuka")
	UErrorInvalidIP                 = common.NewUmuError(2012, "Invalid IPv4 address")
	UErrorRuleNotRegister           = common.NewUmuError(2013, "rule not register")
	UErrorRuleHasNoResult           = common.NewUmuError(2014, "rule has no result")
	UErrorRuleInvalidGraph          = common.NewUmuError(2015, "rule invalid graph")
	UErrorCanceled                  = common.NewUmuError(2016, "client canceled error")
	UErrorFuncOperatorTypeAssertion = common.NewUmuError(2017, "func operator input type assertion error")
)

// UErrorInterfaceNotImplemented 表示接口未被实现的错误
func UErrorInterfaceNotImplemented(interfaceName string) error {
	return common.NewUmuError(2018, fmt.Sprintf("interface '%s' has not been implemented", interfaceName))
}

func IsCloseError(err error) bool {
	if err == nil {
		return false
	}

	if gerrors.Is(err, io.EOF) {
		return true
	}
	if gerrors.Is(err, context.Canceled) {
		return true
	}
	if gerrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	statusError, ok := status.FromError(err)
	if ok {
		if statusError.Code() == codes.DeadlineExceeded ||
			statusError.Code() == codes.Canceled ||
			statusError.Code() == codes.Unavailable {
			return true
		}
	}
	microError := errors.FromError(err)
	if microError != nil {
		if microError.GetCode() == http.StatusRequestTimeout {
			return true
		}
		if microError.GetCode() == 600 && microError.Detail == "EOF" {
			return true
		}
	}

	return false
}

// IsDBDuplicateEntryError 根据mysql错误信息返回错误代码
func IsDBDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	ok := gerrors.As(err, &mysqlErr)
	if ok {
		if mysqlErr.Number == errDuplicateEntryCode {
			return true
		}
	}

	return false
}
