package repository

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/ffajarpratama/boiler-api/pkg/constant"
	custom_error "github.com/ffajarpratama/boiler-api/pkg/error"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository struct{}

func (r *BaseRepository) Create(db *gorm.DB, params interface{}) (err error) {
	err = db.Model(params).Create(params).Error
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusBadRequest,
				Code:     constant.DefaultDuplicateDataError,
			})

			return err
		}

		fmt.Printf("error-query-insert from caller : %s", getCallerFunctionName())
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusInternalServerError,
		})

		return err
	}

	return
}

func (r *BaseRepository) Update(db *gorm.DB, params interface{}) (err error) {
	rows := db.Omit(clause.Associations).Model(params).Updates(params)
	err = rows.Error
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusBadRequest,
				Code:     constant.DefaultDuplicateDataError,
			})

			return err
		}

		fmt.Printf("error-query-update from caller : %s ", getCallerFunctionName())
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusInternalServerError,
		})

		return err
	}

	if rows.RowsAffected == 0 {
		nameField := reflect.TypeOf(params).Elem().Name()
		msg := ""
		if nameField != "" {
			msg = fmt.Sprintf("%s Not found", nameField)
		}

		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusNotFound,
			Code:     constant.DefaultNotFoundError,
			Message:  msg,
		})

		return err
	}

	return
}

func (r *BaseRepository) FindOne(db *gorm.DB, result interface{}) (err error) {
	err = db.First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			nameField := reflect.TypeOf(result).Elem().Name()
			msg := ""
			if nameField != "" {
				msg = fmt.Sprintf("%s Not found", nameField)
			}

			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusNotFound,
				Code:     constant.DefaultNotFoundError,
				Message:  msg,
			})

			return err
		}

		fmt.Printf("error-query from caller : %s", getCallerFunctionName())
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusInternalServerError,
		})

		return err
	}

	return
}

func (r *BaseRepository) Delete(db *gorm.DB, params interface{}) (err error) {
	err = db.Model(params).Delete(params).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			nameField := reflect.TypeOf(params).Elem().Name()
			msg := ""
			if nameField != "" {
				msg = fmt.Sprintf("%s Not found", nameField)
			}

			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusNotFound,
				Code:     constant.DefaultNotFoundError,
				Message:  msg,
			})

			return err
		}

		fmt.Printf("error-query from caller : %s", getCallerFunctionName())
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusInternalServerError,
		})

		return err
	}

	return
}

func IsDuplicateErr(err error) bool {
	value, ok := err.(*custom_error.CustomErrors)
	if !ok {
		return false
	}

	if value.ErrorContext.Code == constant.DefaultDuplicateDataError {
		return true
	}

	return strings.Contains(err.Error(), "Duplicate")
}

func IsRecordNotfound(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}

	value, ok := err.(*custom_error.CustomErrors)
	if !ok {
		return false
	}

	if value == nil {
		return false
	}

	if value.ErrorContext == nil {
		return false
	}

	if value.ErrorContext.HTTPCode == http.StatusNotFound {
		return true
	}

	if value.ErrorContext.Code == constant.DefaultNotFoundError {
		return true
	}

	return false
}

func IsDuplicateConstraintErr(err error) bool {
	return strings.Contains(err.Error(), "command cannot affect row a second time")
}

//nolint:gomnd
func getCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(2).Function
}

//nolint:gomnd
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need

	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
