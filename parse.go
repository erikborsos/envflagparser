// Package envflagparser provides functionality to parse configuration values from both
// environment variables and command-line flags into a provided struct.
// It offers the flexibility to prioritize environment variables over flag values
// The package leverages reflection to dynamically set field values based on their types,
// making it convenient for configuring applications via flags or environment variables.
package envflagparser

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// PrioritiseEnv defines whether environment variables take precedence over flag values.
var PrioritiseEnv = true

// PrintErrorUsage defines whether error messages should include usage information. (flags)
var PrintErrorUsage = false

// ParseConfig parses configuration values from flags and environment variables into the provided struct.
func ParseConfig(configStruct interface{}) (err error) {
	// flag.Parse() panics
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// Panic instead of exit
	flag.CommandLine.Init("envflagparser", flag.PanicOnError)

	// If PrintErrorUsage is false, discard usage information.
	if !PrintErrorUsage {
		flag.CommandLine.SetOutput(io.Discard)
	}

	elem := reflect.ValueOf(configStruct).Elem()
	typ := elem.Type()

	flagValues := make(map[string]interface{})

	// Iterate over fields in the provided struct.
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		// Get flag and environment variable names, default value, and usage information.
		envKey := fieldType.Tag.Get("env")
		flagArgs := fieldType.Tag.Get("flag")
		flagName, defaultValue, usage := parseFlagArgs(flagArgs, envKey)

		// Check if environment variable exists and set the field accordingly.
		envValue, envExists := os.LookupEnv(envKey)
		if envExists {
			setValue(field, envValue)
		}

		// Get flag value based on field type.
		flagSetValue, err := getFlagSetValue(field, flagName, defaultValue, usage)
		if err != nil {
			return err
		}

		flagValues[flagName] = flagSetValue
	}

	// Parse command-line flags.
	flag.Parse()

	// Set field values based on flag values.
	for flagName, flagValue := range flagValues {
		fieldIndex := getFieldIndexByFlagName(typ, flagName)
		if fieldIndex != -1 {
			field := elem.Field(fieldIndex)
			// Check if the field is already set
			// Also if PrioritiseEnv is false, overwrite it
			if !PrioritiseEnv || field.IsZero() {
				if err := setFieldValueByFlagValue(field, flagValue); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// parseFlagArgs splits flag arguments into name, defaultValue and usage strings
func parseFlagArgs(flagArgs string, envKey string) (name, defaultValue, usage string) {
	parts := strings.Split(flagArgs, ";")
	return parts[0], parts[1], parts[2]
}

// getFieldIndexByFlagName retrieves the index of a field by its flag name.
func getFieldIndexByFlagName(typ reflect.Type, flagName string) int {
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		if fieldType.Tag.Get("flag") != "" && strings.Split(fieldType.Tag.Get("flag"), ";")[0] == flagName {
			return i
		}
	}
	return -1
}

// unclean code :(
// TODO: A map with the conversion function

// setValue sets the value of a field based on its type.
func setValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Int, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			// Convert string to duration and set field value.
			durationValue, err := time.ParseDuration(value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(durationValue))
		} else {
			// Convert string to int64 and set field value.
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		}

	case reflect.Uint:
		// Convert string to uint64 and set field value.
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
	case reflect.Uint64:
		// Convert string to uint64 and set field value.
		uint64Value, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uint64Value)
	case reflect.Float64:
		// Convert string to float64 and set field value.
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.String:
		// Set string field value.
		field.SetString(value)
	case reflect.Bool:
		// Convert string to bool and set field value.
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	}
	return nil
}

// getFlagSetValue gets the flag value corresponding to the field type.
func getFlagSetValue(field reflect.Value, flagName, defaultValue, usage string) (interface{}, error) {
	switch field.Kind() {
	case reflect.Int:
		// Convert default value to int and create an Int flag.
		defaultIntValue, err := strconv.Atoi(defaultValue)
		if err != nil {
			return nil, err
		}
		return flag.Int(flagName, defaultIntValue, usage), nil
	case reflect.String:
		// Create a String flag with default value.
		return flag.String(flagName, defaultValue, usage), nil
	case reflect.Bool:
		// Convert default value to bool and create a Bool flag.
		defaultBoolValue, err := strconv.ParseBool(defaultValue)
		if err != nil {
			return nil, err
		}
		return flag.Bool(flagName, defaultBoolValue, usage), nil
	case reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			// Parse default duration value and create a Duration flag.
			defaultDurationValue, err := time.ParseDuration(defaultValue)
			if err != nil {
				return nil, err
			}
			return flag.Duration(flagName, defaultDurationValue, usage), nil
		} else {
			// Convert default value to int64 and create an Int64 flag.
			defaultInt64Value, err := strconv.ParseInt(defaultValue, 10, 64)
			if err != nil {
				return nil, err
			}
			return flag.Int64(flagName, defaultInt64Value, usage), nil
		}
	case reflect.Uint:
		// Convert default value to uint64 and create a Uint flag.
		defaultUintValue, err := strconv.ParseUint(defaultValue, 10, 64)
		if err != nil {
			return nil, err
		}
		return flag.Uint(flagName, uint(defaultUintValue), usage), nil
	case reflect.Uint64:
		// Convert default value to uint64 and create a Uint64 flag.
		defaultUint64Value, err := strconv.ParseUint(defaultValue, 10, 64)
		if err != nil {
			return nil, err
		}
		return flag.Uint64(flagName, defaultUint64Value, usage), nil
	case reflect.Float64:
		// Convert default value to float64 and create a Float64 flag.
		defaultFloatValue, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil {
			return nil, err
		}
		return flag.Float64(flagName, defaultFloatValue, usage), nil
	}
	return nil, nil
}

// setFieldValueByFlagValue sets the value of a field based on the provided flag value.
func setFieldValueByFlagValue(field reflect.Value, flagValue interface{}) error {
	switch fv := flagValue.(type) {
	case *int:
		// Set field value with int.
		setValue(field, strconv.Itoa(*fv))
	case *string:
		// Set field value with string.
		setValue(field, *fv)
	case *bool:
		// Set field value with bool.
		setValue(field, strconv.FormatBool(*fv))
	case *int64:
		// Set field value with int64.
		setValue(field, strconv.FormatInt(*fv, 10))
	case *uint:
		// Set field value with uint.
		setValue(field, strconv.FormatUint(uint64(*fv), 10))
	case *uint64:
		// Set field value with uint64.
		setValue(field, strconv.FormatUint(*fv, 10))
	case *float64:
		// Set field value with float64.
		setValue(field, strconv.FormatFloat(*fv, 'f', -1, 64))
	case *time.Duration:
		// Set field value with duration string.
		setValue(field, (*fv).String())
	default:
		return fmt.Errorf("unsupported flag value type: %T", flagValue)
	}
	return nil
}
