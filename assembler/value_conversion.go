package assembler

import (
	"encoding/binary"
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

func ValuesToLittleEndianBytes(operandValues []any) ([]byte, error) {
	var bytes []byte
	for _, ov := range operandValues {
		var littleEndianBytes []byte
		var err error
		switch v := ov.(type) {
		case int8:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case uint8:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case int16:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case uint16:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case int32:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case uint32:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case int64:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		case uint64:
			littleEndianBytes, err = ToLittleEndianBytes(v)
		default:
			return nil, fmt.Errorf("[valuesToLittleEndianBytes] unsupported integer type: %v", v)
		}

		bytes = append(bytes, littleEndianBytes...)
		if err != nil {
			return nil, err
		}
	}
	return bytes, nil
}

func ToLittleEndianBytes[T constraints.Integer](value T) ([]byte, error) {
	var result []byte
	switch any(value).(type) {
	case int8, uint8:
		result = make([]byte, 1)
		result[0] = byte(value)
	case int16, uint16:
		result = make([]byte, 2)
		binary.LittleEndian.PutUint16(result, uint16(value))
	case int32, uint32:
		result = make([]byte, 4)
		binary.LittleEndian.PutUint32(result, uint32(value))
	case int64, uint64:
		result = make([]byte, 8)
		binary.LittleEndian.PutUint64(result, uint64(value))
	default:
		return nil, fmt.Errorf("[toLittleEndianBytes] unsupported integer type: %v", value)
	}
	return result, nil
}

func ReduceBytes(value any, noOfBytes int) any {
	switch v := value.(type) {
	case int:
		return reduceSigned(int64(v), noOfBytes)
	case int8:
		return reduceSigned(int64(v), noOfBytes)
	case int16:
		return reduceSigned(int64(v), noOfBytes)
	case int32:
		return reduceSigned(int64(v), noOfBytes)
	case int64:
		return reduceSigned(v, noOfBytes)
	case uint:
		return reduceUnsigned(uint64(v), noOfBytes)
	case uint8:
		return reduceUnsigned(uint64(v), noOfBytes)
	case uint16:
		return reduceUnsigned(uint64(v), noOfBytes)
	case uint32:
		return reduceUnsigned(uint64(v), noOfBytes)
	case uint64:
		return reduceUnsigned(v, noOfBytes)
	default:
		return value
	}
}

func reduceSigned(value int64, noOfBytes int) interface{} {
	switch noOfBytes {
	case 1:
		if value <= math.MaxInt8 {
			return int8(value)
		}
	case 2:
		if value <= math.MaxInt16 {
			return int16(value)
		}
	case 3, 4: // We use 4 because Go doesn’t support 3-byte unsigned types
		if value <= math.MaxInt32 {
			return int32(value)
		}
	}
	// If no matching case, default to uint64
	return int64(value)
}

func reduceUnsigned(value uint64, noOfBytes int) interface{} {
	switch noOfBytes {
	case 1:
		if value <= math.MaxUint8 {
			return uint8(value)
		}
	case 2:
		if value <= math.MaxUint16 {
			return uint16(value)
		}
	case 3, 4: // We use 4 because Go doesn’t support 3-byte unsigned types
		if value <= math.MaxUint32 {
			return uint32(value)
		}
	}
	// If no matching case, default to uint64
	return uint64(value)
}
