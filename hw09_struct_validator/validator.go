package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	ValidationTag     = "validate"
	CombineTag        = "|"
	LenTag            = "len"
	RegexpTag         = "regexp"
	InTag             = "in"
	MinValueTag       = "min"
	MaxValueTag       = "max"
	ValuesSeparator   = ","
	KeyValueSeparator = ":"
)

type ValidationRules struct {
	rules []Rule
}
type Rule struct {
	Name, Value string
}

var (
	ErrInputIsNotStruct = errors.New("input argument is not struct")

	ErrInvalidFieldTag = errors.New("invalid tag format")
	ErrInvalidRegexp   = errors.New("invalid regular expression")

	ErrFailedRegexp = errors.New("failed regexp validation")
	ErrFailedMinMax = errors.New("failed min or max value validation")
	ErrFailedLen    = errors.New("failed len validation")
	ErrFailedIn     = errors.New("failed in validation")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "Validation succeeded."
	}

	errMessage := strings.Builder{}
	errMessage.WriteString("Validation failed. Details:")

	for _, err := range v {
		s := fmt.Sprintf("Field %s validation failed. Reason: %s\n", err.Field, err.Err)
		errMessage.WriteString(s)
	}

	return errMessage.String()
}

func Validate(v interface{}) error {
	if v == nil {
		return nil
	}

	validationErrors := ValidationErrors{}

	vType := reflect.TypeOf(v)
	vValue := reflect.ValueOf(&v).Elem()

	if vType.Kind() != reflect.Struct {
		return ErrInputIsNotStruct
	}

	for i := 0; i < vType.NumField(); i++ {
		validationErr, ok := validateFieldByIndex(i, vValue, vType)
		if !ok {
			continue
		}

		if validationErr.Err != nil {
			validationErrors = append(validationErrors, validationErr)
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateFieldByIndex(index int, structValue reflect.Value, structType reflect.Type) (ValidationError, bool) {
	field := structType.Field(index)
	validationErr := ValidationError{Field: field.Name}

	tagValue, ok := field.Tag.Lookup(ValidationTag)
	if !ok {
		return validationErr, false
	}

	vr, err := getValidationRules(tagValue)
	if err != nil {
		validationErr.Err = err
		return validationErr, true
	}

	fieldElement := structValue.Elem().Field(index)
	fieldValue := fieldElement.Interface()

	switch assertedValue := fieldValue.(type) {
	case int:
		err = vr.checkFieldIntValue(assertedValue)
	case []int:
		for _, itemValue := range assertedValue {
			err = vr.checkFieldIntValue(itemValue)
			if err != nil {
				validationErr.Err = err
				break
			}
		}

	case string:
		err = vr.checkFieldStringValue(assertedValue)
	case []string:
		for _, itemValue := range assertedValue {
			err = vr.checkFieldStringValue(itemValue)
			if err != nil {
				validationErr.Err = err
				break
			}
		}

	default:
		err = vr.validateNested(fieldElement, structType.Field(index))
	}

	validationErr.Err = err

	return validationErr, true
}

func (vr ValidationRules) validateNested(structElement reflect.Value, structField reflect.StructField) error {
	structValue := structElement.Interface()
	structKind := structElement.Type().Kind()

	if structKind == reflect.String {
		return vr.checkFieldStringValue(fmt.Sprintf("%v", structValue))
	} else if structKind == reflect.Int {
		sValue := fmt.Sprintf("%v", structValue)
		v, err := strconv.Atoi(sValue)
		if err != nil {
			return err
		}
		return vr.checkFieldIntValue(v)
	}

	return Validate(structField)
}

func getValidationRules(tagValue string) (ValidationRules, error) {
	vr := ValidationRules{rules: make([]Rule, 0)}
	tags := strings.Split(tagValue, CombineTag)
	for _, tag := range tags {
		rule, err := getRuleFromTag(tag)
		if err != nil {
			return vr, err
		}
		vr.rules = append(vr.rules, rule)
	}
	return vr, nil
}

func getRuleFromTag(tag string) (Rule, error) {
	rule := Rule{}

	keyValue := strings.Split(tag, KeyValueSeparator)
	if len(keyValue) != 2 {
		return rule, ErrInvalidFieldTag
	}

	rule.Name = keyValue[0]
	rule.Value = keyValue[1]

	return rule, nil
}

func (vr ValidationRules) checkFieldIntValue(value int) error {
	for _, rule := range vr.rules {
		ruleValueNumber, err := strconv.Atoi(rule.Value)
		if err != nil && rule.Name != InTag {
			return ErrInvalidFieldTag
		}

		switch rule.Name {
		case MinValueTag:
			if value < ruleValueNumber {
				return ErrFailedMinMax
			}

		case MaxValueTag:
			if ruleValueNumber < value {
				return ErrFailedMinMax
			}

		case InTag:
			if !isValueBetween(rule.Value, strconv.Itoa(value)) {
				return ErrFailedIn
			}
		}
	}
	return nil
}

func (vr ValidationRules) checkFieldStringValue(value string) error {
	for _, rule := range vr.rules {
		switch rule.Name {
		case LenTag:
			ruleValueNumber, err := strconv.Atoi(rule.Value)
			if err != nil {
				return ErrInvalidFieldTag
			}
			if len(value) != ruleValueNumber {
				return ErrFailedLen
			}

		case RegexpTag:
			ok, err := validateRegexp(rule.Value, value)
			if err != nil {
				return ErrInvalidRegexp
			}
			if !ok {
				return ErrFailedRegexp
			}

		case InTag:
			if !isValueBetween(rule.Value, value) {
				return ErrFailedIn
			}
		}
	}
	return nil
}

func validateRegexp(exp, value string) (bool, error) {
	regex, err := regexp.Compile(exp)
	if err != nil {
		return false, err
	}

	return regex.MatchString(value), nil
}

func isValueBetween(between string, value string) bool {
	values := strings.Split(between, ValuesSeparator)
	for _, s := range values {
		if value == s {
			return true
		}
	}
	return false
}
