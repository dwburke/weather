package validate

import (
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

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

			err := validate.Struct(resource)

			if verr := err.(validator.ValidationErrors); verr != nil {
				scope.DB().AddError(verr)
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
