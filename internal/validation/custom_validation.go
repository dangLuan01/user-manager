package validation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("minInt", func(fl validator.FieldLevel) bool {
		min, _ := strconv.ParseInt(fl.Param(), 10, 64)
		if fl.Field().Int() < min {
			return false
		}
		return true
	})
	v.RegisterValidation("maxInt", func(fl validator.FieldLevel) bool {
		max, _ := strconv.ParseInt(fl.Param(), 10, 64)
		if fl.Field().Int() > max {
			return false
		}
		return true
	})
	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		fileName := fl.Field().String()
		ext := fileName[strings.LastIndex(fileName, ".")+1:]
		validExts := strings.Fields(fl.Param())
		for _, validExt := range validExts {
			if strings.EqualFold(ext, validExt) {
				return true
			}
		}

		return false
	})
}