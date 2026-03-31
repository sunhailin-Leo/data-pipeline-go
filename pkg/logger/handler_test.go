package logger

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/utils"
)

func TestGetEncodeConfig(t *testing.T) {
	configWithServiceName := getEncodeConfig(true)
	assert.Equal(t, "logger", configWithServiceName.NameKey)
	assert.Equal(t, "linenum", configWithServiceName.CallerKey)
	assert.Equal(t, "level", configWithServiceName.LevelKey)

	configWithoutServiceName := getEncodeConfig(false)
	assert.Equal(t, "logger", configWithoutServiceName.NameKey)
	assert.Equal(t, "linenum", configWithoutServiceName.CallerKey)
	assert.Equal(t, "level", configWithoutServiceName.LevelKey)
}

func TestGetRollingConfig(t *testing.T) {
	filename := "test.log"
	expectedLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    utils.LogMaxSize,
		MaxAge:     utils.LogMaxAge,
		MaxBackups: utils.LogMaxBackups,
		Compress:   true,
	}
	actualLogger := getRollingConfig(filename)
	assert.Equal(t, expectedLogger, actualLogger)
}

// TestNewZapLogger tests that NewZapLogger initializes the Logger correctly
func TestNewZapLogger(t *testing.T) {
	NewZapLogger()
	assert.NotNil(t, Logger)
}

// TestNewZapLogger_LogOutput tests that Logger.Info does not panic after initialization
func TestNewZapLogger_LogOutput(t *testing.T) {
	NewZapLogger()
	assert.NotNil(t, Logger)

	// This should not panic
	assert.NotPanics(t, func() {
		Logger.Info("test")
	})
}

// mockPrimitiveArrayEncoder implements zapcore.PrimitiveArrayEncoder for testing
type mockPrimitiveArrayEncoder struct {
	strings []string
	ints    []int64
}

func (m *mockPrimitiveArrayEncoder) AppendBool(bool)              {}
func (m *mockPrimitiveArrayEncoder) AppendByteString([]byte)      {}
func (m *mockPrimitiveArrayEncoder) AppendComplex128(complex128)  {}
func (m *mockPrimitiveArrayEncoder) AppendComplex64(complex64)    {}
func (m *mockPrimitiveArrayEncoder) AppendFloat64(float64)        {}
func (m *mockPrimitiveArrayEncoder) AppendFloat32(float32)        {}
func (m *mockPrimitiveArrayEncoder) AppendInt(int)                {}
func (m *mockPrimitiveArrayEncoder) AppendInt64(v int64)          { m.ints = append(m.ints, v) }
func (m *mockPrimitiveArrayEncoder) AppendInt32(int32)            {}
func (m *mockPrimitiveArrayEncoder) AppendInt16(int16)            {}
func (m *mockPrimitiveArrayEncoder) AppendInt8(int8)              {}
func (m *mockPrimitiveArrayEncoder) AppendString(v string)        { m.strings = append(m.strings, v) }
func (m *mockPrimitiveArrayEncoder) AppendUint(uint)              {}
func (m *mockPrimitiveArrayEncoder) AppendUint64(uint64)          {}
func (m *mockPrimitiveArrayEncoder) AppendUint32(uint32)          {}
func (m *mockPrimitiveArrayEncoder) AppendUint16(uint16)          {}
func (m *mockPrimitiveArrayEncoder) AppendUint8(uint8)            {}
func (m *mockPrimitiveArrayEncoder) AppendUintptr(uintptr)        {}
func (m *mockPrimitiveArrayEncoder) AppendDuration(time.Duration) {}
func (m *mockPrimitiveArrayEncoder) AppendTime(time.Time)         {}

// TestTimeEncoderFunc tests that timeEncoderFunc does not panic and encodes correctly
func TestTimeEncoderFunc(t *testing.T) {
	enc := &mockPrimitiveArrayEncoder{}
	assert.NotPanics(t, func() {
		timeEncoderFunc(time.Now(), enc)
	})
	assert.Len(t, enc.strings, 1)
}

// TestDurationEncoderFunc tests that durationEncoderFunc does not panic and encodes correctly
func TestDurationEncoderFunc(t *testing.T) {
	enc := &mockPrimitiveArrayEncoder{}
	assert.NotPanics(t, func() {
		durationEncoderFunc(time.Second, enc)
	})
	assert.Len(t, enc.ints, 1)
	assert.Equal(t, int64(1000), enc.ints[0])
}

// TestJsonEncoderFunc tests that jsonEncoderFunc returns a non-nil encoder
func TestJsonEncoderFunc(t *testing.T) {
	var buf bytes.Buffer
	encoder := jsonEncoderFunc(&buf)
	assert.NotNil(t, encoder)
}

// TestCallerEncoderFunc tests that callerEncoderFunc returns a non-nil encoder
func TestCallerEncoderFunc(t *testing.T) {
	t.Run("with service name", func(t *testing.T) {
		enc := &mockPrimitiveArrayEncoder{}
		callerFunc := callerEncoderFunc(true)
		assert.NotNil(t, callerFunc)
		assert.NotPanics(t, func() {
			callerFunc(zapcore.EntryCaller{Defined: true, File: "test.go", Line: 42}, enc)
		})
		assert.Len(t, enc.strings, 2) // service name + caller
	})

	t.Run("without service name", func(t *testing.T) {
		enc := &mockPrimitiveArrayEncoder{}
		callerFunc := callerEncoderFunc(false)
		assert.NotNil(t, callerFunc)
		assert.NotPanics(t, func() {
			callerFunc(zapcore.EntryCaller{Defined: true, File: "test.go", Line: 42}, enc)
		})
		assert.Len(t, enc.strings, 1) // only caller
	})
}
