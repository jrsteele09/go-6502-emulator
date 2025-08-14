package assembler

import (
	"encoding/binary"
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
)

// toUint64 converts various integer types to uint64
func toUint64(value any) (uint64, error) {
	switch v := value.(type) {
	case uint64:
		return v, nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("negative value not allowed: %d", v)
		}
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("negative value not allowed: %d", v)
		}
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("negative value not allowed: %d", v)
		}
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("negative value not allowed: %d", v)
		}
		return uint64(v), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("negative value not allowed: %d", v)
		}
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("invalid integer type: %T", value)
	}
}

// parseOperandSize determines the operand size and converts values for assembly
func parseOperandSize(negative bool, value any) (string, any, error) {
	var finalValue int64
	var naturalSize int // Natural size in bytes of the input type

	// Convert to int64 and determine natural size
	switch v := value.(type) {
	case int8:
		naturalSize = 1
		finalValue = int64(v)
	case uint8:
		naturalSize = 1
		finalValue = int64(v)
	case int16:
		naturalSize = 2
		finalValue = int64(v)
	case uint16:
		naturalSize = 2
		finalValue = int64(v)
	case int32:
		naturalSize = 4
		finalValue = int64(v)
	case uint32:
		naturalSize = 4
		finalValue = int64(v)
	case int64:
		naturalSize = 8
		finalValue = v
	case uint64:
		naturalSize = 8
		finalValue = int64(v)
	case int:
		naturalSize = 8
		finalValue = int64(v)
	default:
		return "", nil, fmt.Errorf("invalid operand type: %T", value)
	}

	// Apply negative flag
	if negative && finalValue > 0 {
		finalValue = -finalValue
	}

	// Determine required size based on value and negative flag
	var requiredSize int
	if finalValue >= -128 && finalValue <= 255 {
		requiredSize = 1
	} else if finalValue >= -32768 && finalValue <= 65535 {
		requiredSize = 2
	} else {
		return "", nil, fmt.Errorf("[parseOperandSize] Number too large: %d", finalValue)
	}

	// If negative flag caused size promotion, use the larger size
	if negative && naturalSize == 1 && (finalValue < -128 || finalValue > 127) {
		requiredSize = 2
	}

	// Generate size mask and reduced value
	var sizeMask string
	var reducedValue any

	correctSign := func() any {
		if negative {
			return finalValue
		}
		return uint64(finalValue)
	}

	switch requiredSize {
	case 1:
		sizeMask = "nn"
		reducedValue = ReduceBytes(correctSign(), 1)
	case 2:
		sizeMask = "nnnn"
		reducedValue = ReduceBytes(correctSign(), 2)
	default:
		return "", nil, fmt.Errorf("unsupported operand size: %d bytes", requiredSize)
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
