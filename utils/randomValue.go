package utils

import (
	"fmt"
	"go/types"
	"math"
	"math/rand"
	"reflect"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var MAX_MUTATE_ITER = 3

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Generate_random_bool() bool {
	return seededRand.Intn(2) != 0
}

func generate_edge_value_with_type(number interface{}) (interface{}, interface{}) {
	switch number.(type) {
	case uint:
		return uint(0), uint(math.MaxUint32)
	case uint8:
		return uint8(0), uint8(math.MaxUint8)
	case uint16:
		return uint16(0), uint16(math.MaxUint16)
	case uint32:
		return uint32(0), uint32(math.MaxUint32)
	case uint64:
		return uint64(0), uint64(math.MaxUint64)
	case int:
		return int(math.MinInt32), int(math.MaxInt32)
	case int8:
		return int8(math.MinInt8), int8(math.MaxInt8)
	case int16:
		return int16(math.MinInt16), int16(math.MaxInt16)
	case int32:
		return int32(math.MinInt32), int32(math.MaxInt32)
	case int64:
		return int64(math.MinInt64), int64(math.MaxInt64)
	case float32:
		return -math.MaxFloat32, math.MaxFloat32
	case float64:
		return float64(-math.MaxFloat64), float64(math.MaxFloat64)
	default:
		fmt.Printf("%s is an unknown number type!\n", reflect.TypeOf(number).String())
		return false, false
	}
}

func small_change_fundation_number[T Numeric](number T) T {
	changeRange := seededRand.Intn(5) + 1
	randChoice := Generate_random_bool()

	minNumber, maxNumber := generate_edge_value_with_type(number)

	if randChoice {
		if number > maxNumber.(T)-T(changeRange) {
			return maxNumber.(T)
		}
		return number + T(changeRange)
	}
	if number < minNumber.(T)+T(changeRange) {
		return minNumber.(T)
	}
	return number - T(changeRange)
}

/* random mutation for string */
func random_mutate_string(str string) string {
	length := len(str)

	if length <= 0 {
		return str
	}

	mutate_times := seededRand.Intn(MAX_MUTATE_ITER + 1)

	res := str
	for i := 0; i < mutate_times; i++ {
		pos := seededRand.Intn(length)
		res = res[:pos] + string(charset[seededRand.Intn(len(charset))]) + res[pos+1:]
	}
	return res
}

func random_mutate_bytes(bytes []byte) []byte {
	length := len(bytes)

	if length <= 0 {
		return bytes
	}

	mutate_times := rand.Intn(MAX_MUTATE_ITER + 1)

	res := bytes

	for i := 0; i < mutate_times; i++ {
		pos := seededRand.Intn(length)
		res[pos] = byte(rand.Intn(math.MaxUint8))
	}
	return res
}

func handle_number_mutate(number interface{}) interface{} {
	if seededRand.Intn(100) < 80 {
		return number
	}

	reflect.TypeOf(number)
	switch number := number.(type) {
	case uint:
		return small_change_fundation_number(number)
	case uint8:
		return small_change_fundation_number(number)
	case uint16:
		return small_change_fundation_number(number)
	case uint32:
		return small_change_fundation_number(number)
	case uint64:
		return small_change_fundation_number(number)
	case int:
		return small_change_fundation_number(number)
	case int8:
		return small_change_fundation_number(number)
	case int16:
		return small_change_fundation_number(number)
	case int32:
		return small_change_fundation_number(number)
	case int64:
		return small_change_fundation_number(number)
	case float32:
		return small_change_fundation_number(number)
	case float64:
		return small_change_fundation_number(number)
	default:
		return number
	}
}

// 生成与原始value不同的新值
func GenerateDiffValue(value interface{}) interface{} {
	if value == nil {
		fmt.Println("nil value!")
		return nil
	}
	for {
		var result interface{}
		ident := reflect.ValueOf(value).Kind()
		switch ident {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			result = handle_number_mutate(value)
		case reflect.String:
			result = random_mutate_string(value.(string))
		case reflect.Slice:
			if value, ok := value.([]byte); ok {
				result = random_mutate_bytes(value)
			}
			result = random_mutate_bytes([]byte("0"))
		default:
			return random_mutate_string("string")
		}

		if result != value {
			return result
		}
	}
}

// 根据类型返回对应的值
func ParseType(t types.Type) interface{} {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int:
			return int(0)
		case types.Int8:
			return int8(0)
		case types.Int16:
			return int16(0)
		case types.Int32:
			return int32(0)
		case types.Int64:
			return int64(0)
		case types.Uint:
			return uint(0)
		case types.Uint8:
			return uint8(0)
		case types.Uint16:
			return uint16(0)
		case types.Uint32:
			return uint32(0)
		case types.Uint64:
			return uint64(0)
		case types.Float32:
			return float32(0)
		case types.Float64:
			return float64(0)
		case types.String:
			return "string"
		case types.Bool:
			return false
		}
	case *types.Struct:
		result := make(map[string]interface{})
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			fieldName := field.Name()

			result[fieldName] = ParseType(field.Type().Underlying())
		}
		return result

	// 暂时不接受任何map类型的输入
	/*
		case *types.Map:
			// maybe bugs present
			result := make(map[interface{}]interface{})
			key := ParseType(t.Key())
			result[key] = ParseType(t.Elem())

			return result
	*/

	case *types.Pointer:
		return ParseType(t.Elem())
	case *types.Slice:
		result := make([]interface{}, 1)
		result[0] = ParseType(t.Elem())
		return result
	case *types.Named:
		return ParseType(t.Obj().Type().Underlying())
	default:
		return "string"
	}

	return "string"
}
