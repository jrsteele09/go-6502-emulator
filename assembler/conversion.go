package assembler

import (
	"encoding/binary"
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// toUint64 converts various integer types to uint64
func toUint64(value any) (uint64, error) {
	i64, err := toInt64(value)
	if err != nil {
		return 0, err
	}
	if i64 < 0 {
		return 0, fmt.Errorf("negative value not allowed: %d", i64)
	}
	return uint64(i64), nil
}

// toInt64 converts various integer types to int64
func toInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case int:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("invalid integer type: %T", value)
	}
}

// minimumOperandSize converts any integer type to operand size mask and value
// Returns "nn" for 8-bit values, "nnnn" for 16-bit values, etc.
// Handles negative flag by promoting to larger size if needed
func minimumOperandSize(negative bool, value any) (string, any, error) {
	finalValue, err := toInt64(value)
	if err != nil {
		return "", nil, fmt.Errorf("invalid operand type: %w", err)
	}

	// Apply negative flag
	if negative && finalValue > 0 {
		finalValue = -finalValue
	}

	// Determine required size based on value range
	var sizeMask string
	var reducedValue any

	switch {
	case finalValue >= -128 && finalValue <= 255:
		sizeMask = oneByteOperand
		reducedValue = ReduceBytes(finalValue, 1)
	case finalValue >= -32768 && finalValue <= 65535:
		sizeMask = twoByteOperand
		reducedValue = ReduceBytes(finalValue, 2)
	default:
		return "", nil, fmt.Errorf("[parseOperandSize] Number too large: %d", finalValue)
	}

	return sizeMask, reducedValue, nil
}

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
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, littleEndianBytes...)
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
	// Convert to int64 first to determine sign
	i64, err := toInt64(value)
	if err != nil {
		return value // Return original if can't convert
	}

	if i64 < 0 {
		return reduceSigned(i64, noOfBytes)
	}
	return reduceUnsigned(uint64(i64), noOfBytes)
}

func reduceSigned(value int64, noOfBytes int) any {
	switch noOfBytes {
	case 1:
		if value >= math.MinInt8 && value <= math.MaxInt8 {
			return int8(value)
		}
	case 2:
		if value >= math.MinInt16 && value <= math.MaxInt16 {
			return int16(value)
		}
	case 3, 4: // We use 4 because Go doesn't support 3-byte types
		if value >= math.MinInt32 && value <= math.MaxInt32 {
			return int32(value)
		}
	}
	return value // Return as int64 if no smaller type fits
}

func reduceUnsigned(value uint64, noOfBytes int) any {
	switch noOfBytes {
	case 1:
		if value <= math.MaxUint8 {
			return uint8(value)
		}
	case 2:
		if value <= math.MaxUint16 {
			return uint16(value)
		}
	case 3, 4: // We use 4 because Go doesn't support 3-byte types
		if value <= math.MaxUint32 {
			return uint32(value)
		}
	}
	return value // Return as uint64 if no smaller type fits
}
