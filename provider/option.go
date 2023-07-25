package provider

import (
	"errors"
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
	"strings"

	"github.com/malt3/abstractfs-core/sri"
)

func SetOptions(val any, opts map[string]string) error {
	availableOpts := Options(val)
	for k, v := range opts {
		if isUnknownOpt(k, availableOpts) {
			return fmt.Errorf("unknown option %q (valid options are %v)", k, availableOpts)
		}
		converted, err := valueFromString(v, OptionType(val, k))
		if err != nil {
			return err
		}
		if err := OptionSet(val, k, converted); err != nil {
			return err
		}
	}
	return nil
}

// Options uses struct field tags to determine the options of a type.
func Options(v any) []string {
	opts := []string{}
	t := reflect.TypeOf(v).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		opt := field.Tag.Get("abstractfs")
		if opt == "" {
			continue
		}
		opts = append(opts, opt)
	}
	return opts
}

// OptionType returns the type of the option.
func OptionType(v any, opt string) reflect.Type {
	t := reflect.TypeOf(v).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("abstractfs") == opt {
			return field.Type
		}
	}
	return nil
}

// OptionGet returns the value of the option as any.
func OptionGet(v any, opt string) any {
	t := reflect.TypeOf(v).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("abstractfs") == opt {
			return reflect.ValueOf(v).Elem().Field(i).Addr().Interface()
		}
	}
	return nil
}

// OptionSet sets the value of the option.
func OptionSet(v any, opt string, value any) error {
	t := reflect.TypeOf(v).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("abstractfs") == opt {
			if !reflect.ValueOf(v).Elem().Field(i).CanSet() {
				return errors.New("cannot set option")
			}
			reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(value))
			return nil
		}
	}
	return errors.New("option not found")
}

func isUnknownOpt(opt string, opts []string) bool {
	for _, o := range opts {
		if o == opt {
			return false
		}
	}
	return true
}

func valueFromString(s string, t reflect.Type) (any, error) {
	switch t {
	case typeSRIAlgorithm:
		return sri.AlgorithmFromString(s)
	}
	switch t.Kind() {
	case reflect.String:
		return s, nil
	case reflect.Bool:
		return boolFromString(s)
	case reflect.Int:
		return intFromString[int](s)
	case reflect.Int8:
		return intFromString[int8](s)
	case reflect.Int16:
		return intFromString[int16](s)
	case reflect.Int32:
		return intFromString[int32](s)
	case reflect.Int64:
		return intFromString[int64](s)
	case reflect.Uint:
		return intFromString[uint](s)
	case reflect.Uint8:
		return intFromString[uint8](s)
	case reflect.Uint16:
		return intFromString[uint16](s)
	case reflect.Uint32:
		return intFromString[uint32](s)
	case reflect.Uint64:
		return intFromString[uint64](s)
	case reflect.Float32:
		return floatFromString[float32](s)
	case reflect.Float64:
		return floatFromString[float64](s)
	case reflect.Slice:
		return sliceFromString(s, t)
	}

	return nil, errors.New("invalid type")
}

func sliceFromString(s string, t reflect.Type) (any, error) {
	elements := strings.Split(s, ",")
	results := reflect.MakeSlice(t, len(elements), len(elements))
	for i, e := range elements {
		v, err := valueFromString(e, t.Elem())
		if err != nil {
			return nil, err
		}
		results.Index(i).Set(reflect.ValueOf(v))
	}
	return results.Interface(), nil
}

func boolFromString(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "1":
		return true, nil
	case "0":
		return false, nil
	case "yes":
		return true, nil
	case "no":
		return false, nil
	case "y":
		return true, nil
	case "n":
		return false, nil
	case "on":
		return true, nil
	case "off":
		return false, nil
	case "enable":
		return true, nil
	case "disable":
		return false, nil
	}
	return false, errors.New("invalid bool")
}

func intFromString[V integer](s string) (V, error) {
	bitSize := bitSize[V]()
	if isUnsigned[V]() {
		u, err := strconv.ParseUint(s, 0, bitSize)
		if err != nil {
			return 0, err
		}
		return V(u), nil
	}
	i, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}
	return V(i), nil
}

func isUnsigned[V integer]() bool {
	var v V
	switch any(v).(type) {
	case uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func bitSize[V integer | float]() int {
	var v V
	switch any(v).(type) {
	case int8, uint8:
		return 8
	case int16, uint16:
		return 16
	case int32, uint32, float32:
		return 32
	case int64, uint64, float64:
		return 64
	case int, uint:
		return bits.UintSize
	}
	// unreachable
	return 0
}

func floatFromString[V float](s string) (V, error) {
	bitSize := bitSize[V]()
	i, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return 0, err
	}
	return V(i), nil
}

var (
	typeSRIAlgorithm = reflect.TypeOf(sri.SHA256)
)

type integer interface {
	int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint
}

type float interface {
	float32 | float64
}
