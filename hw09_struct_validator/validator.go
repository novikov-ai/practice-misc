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

type Rule struct {
	Name, Value string
}

var (
	ErrInvalidFields    = errors.New("invalid struct values")
	ErrValidationFailed = errors.New("validation failed")

	ErrInputIsNotStruct = errors.New("input argument is not struct")

	ErrInvalidFieldTag = fmt.Errorf("%w: field tag", ErrInvalidFields)
	ErrInvalidRegexp   = fmt.Errorf("%w: invalid regular expression", ErrInvalidFields)

	ErrFailedRegexp = fmt.Errorf("%w: invalid regexp", ErrValidationFailed)
	ErrFailedMinMax = fmt.Errorf("%w: invalid min or max value", ErrValidationFailed)
	ErrFailedLen    = fmt.Errorf("%w: invalid len", ErrValidationFailed)
	ErrFailedIn     = fmt.Errorf("%w: invalid in value", ErrValidationFailed)
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
		field := vType.Field(i)

		tagValue, ok := getTagAndValidate(field)
		if !ok {
			continue
		}

		structValue := getStructValueElement(vValue, i)

		validationErr := validateFieldByIndex(field, structValue, tagValue)
		if validationErr == nil {
			continue
		}

		validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: validationErr})
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func getTagAndValidate(field reflect.StructField) (string, bool) {
	tag, ok := field.Tag.Lookup(ValidationTag)
	return tag, ok
}

func getStructValueElement(structValue reflect.Value, index int) reflect.Value {
	return structValue.Elem().Field(index)
}

func validateFieldByIndex(fieldType reflect.StructField, structValue reflect.Value, tagValue string) error {
	validationRules, err := getValidationRules(tagValue)
	if err != nil {
		return err
	}

	fieldValue := structValue.Interface()
	switch assertedValue := fieldValue.(type) {
	case int:
		err = checkFieldIntValue(validationRules, assertedValue)
	case []int:
		for _, itemValue := range assertedValue {
			err = checkFieldIntValue(validationRules, itemValue)
			if err != nil {
				break
			}
		}

	case string:
		err = checkFieldStringValue(validationRules, assertedValue)
	case []string:
		for _, itemValue := range assertedValue {
			err = checkFieldStringValue(validationRules, itemValue)
			if err != nil {
				break
			}
		}

	default:
		err = validateNested(validationRules, structValue, fieldType)
	}

	return err
}

func validateNested(vr []Rule, structElement reflect.Value, structField reflect.StructField) error {
	structValue := structElement.Interface()
	structKind := structElement.Type().Kind()

	if structKind == reflect.String {
		return checkFieldStringValue(vr, fmt.Sprintf("%v", structValue))
	} else if structKind == reflect.Int {
		sValue := fmt.Sprintf("%v", structValue)
		v, err := strconv.Atoi(sValue)
		if err != nil {
			return err
		}
		return checkFieldIntValue(vr, v)
	}

	return Validate(structField)
}

func getValidationRules(tagValue string) ([]Rule, error) {
	rules := make([]Rule, 0)
	tags := strings.Split(tagValue, CombineTag)
	for _, tag := range tags {
		rule, err := getRuleFromTag(tag)
		if err != nil {
			return rules, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
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

func checkFieldIntValue(rules []Rule, value int) error {
	for _, rule := range rules {
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

func checkFieldStringValue(rules []Rule, value string) error {
	for _, rule := range rules {
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
