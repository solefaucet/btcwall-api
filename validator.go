package main

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/solefaucet/btcwall-api/utils"
	validator "gopkg.in/go-playground/validator.v8"
)

func initializeGinValidator() {
	must(nil, validate.RegisterValidation("btc_addr", bitcoinAddressValidator))
	binding.Validator = extendedValidator{}
}

var validate = validator.New(&validator.Config{TagName: "binding"})

type extendedValidator struct{}

var _ binding.StructValidator = extendedValidator{}

func (v extendedValidator) ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}

func bitcoinAddressValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	return utils.ValidateBitcoinAddress(field.String())
}
