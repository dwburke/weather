package validate

import (
	"sync"

	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

var (
	validatorInstance *validator.Validate
	validatorOnce     sync.Once
)

// getValidator returns a cached validator instance
func getValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validatorInstance = validator.New()
	})
	return validatorInstance
}

func validate(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {

		if scope.HasError() {
			return
		}

		scope.CallMethod("Validate")
		if scope.HasError() {
			return
		}

		if scope.Value != nil {
			resource := scope.IndirectValue().Interface()

			err := getValidator().Struct(resource)

			if err != nil {
				if verr, ok := err.(validator.ValidationErrors); ok && verr != nil {
					scope.DB().AddError(verr)
				} else {
					scope.DB().AddError(err)
				}
			}
		}

	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("validate:validate") == nil {
		callback.Create().Before("gorm:before_create").Register("validate:validate", validate)
	}
	if callback.Update().Get("validate:validate") == nil {
		callback.Update().Before("gorm:before_update").Register("validate:validate", validate)
	}
}
