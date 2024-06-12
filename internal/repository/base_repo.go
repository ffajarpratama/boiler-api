package repository

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/ffajarpratama/boiler-api/constant"
	"github.com/ffajarpratama/boiler-api/lib/custom_error"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository struct{}

func (r *BaseRepository) Create(db *gorm.DB, params interface{}) (err error) {
	err = db.Model(params).Create(params).Error
	if err != nil {
		fmt.Printf("[error] %v\n", err.Error())

		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusConflict,
				Message:  constant.HTTPStatusText(http.StatusConflict),
			})

			return
		}

		return
	}

	return
}

func (r *BaseRepository) Update(db *gorm.DB, params interface{}) (err error) {
	rows := db.Omit(clause.Associations).Model(params).Updates(params)
	err = rows.Error
	if err != nil {
		fmt.Printf("[error] %v\n", err.Error())

		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusConflict,
				Message:  constant.HTTPStatusText(http.StatusConflict),
			})

			return
		}

		return
	}

	if rows.RowsAffected == 0 {
		nameField := reflect.TypeOf(params).Elem().Name()
		msg := ""
		if nameField != "" {
			msg = fmt.Sprintf("%s Not found", nameField)
		}

		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusNotFound,
			Message:  msg,
		})

		return
	}

	return
}

func (r *BaseRepository) FindOne(db *gorm.DB, result interface{}) (err error) {
	err = db.First(result).Error
	if err != nil {
		fmt.Printf("[error] %v\n", err.Error())

		if err == gorm.ErrRecordNotFound {
			nameField := reflect.TypeOf(result).Elem().Name()
			msg := ""
			if nameField != "" {
				msg = fmt.Sprintf("%s Not found", nameField)
			}

			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusNotFound,
				Message:  msg,
			})
		}

		return
	}

	return
}

func (r *BaseRepository) Delete(db *gorm.DB, params interface{}) (err error) {
	err = db.Model(params).Delete(params).Error
	if err != nil {
		fmt.Printf("[error] %v\n", err.Error())

		if err == gorm.ErrRecordNotFound {
			nameField := reflect.TypeOf(params).Elem().Name()
			msg := ""
			if nameField != "" {
				msg = fmt.Sprintf("%s Not found", nameField)
			}

			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusNotFound,
				Message:  msg,
			})
		}

		return
	}

	return
}

func IsDuplicateErr(err error) bool {
	if strings.Contains(err.Error(), "duplicate") {
		return true
	}

	val, ok := err.(*custom_error.CustomError)
	if !ok {
		return false
	}

	return val.ErrorContext.HTTPCode == http.StatusConflict
}

func IsRecordNotfound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
