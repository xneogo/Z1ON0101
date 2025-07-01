/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/10/28 -- 11:51
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: encoder.go
*/

package xlog

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// For JSON-escaping; see textEncoder.safeAddString below.
const _hex = "0123456789abcdef"

var bufferPool = buffer.NewPool()

var _textPool = sync.Pool{New: func() interface{} {
	return &textEncoder{}
}}

func getTextEncoder() *textEncoder {
	return _textPool.Get().(*textEncoder)
}

func putTextEncoder(enc *textEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.TextEncoderConfig = nil
	enc.buf = nil
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_textPool.Put(enc)
}

type textEncoder struct {
	*TextEncoderConfig
	buf            *buffer.Buffer
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc *json.Encoder
}

// TextEncoderConfig custom config with extra info
type TextEncoderConfig struct {
	*zapcore.EncoderConfig
	FieldSeparator string `json:"fieldSeparator" yaml:"fieldSeparator"`
}

// NewAccessLogEncoder constructor
func NewAccessLogEncoder(cfg TextEncoderConfig) zapcore.Encoder {
	return newTextEncoder(cfg, false)
}

func newTextEncoder(cfg TextEncoderConfig, spaced bool) *textEncoder {
	return &textEncoder{
		TextEncoderConfig: &cfg,
		buf:               bufferPool.Get(),
	}
}

func (enc *textEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *textEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *textEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *textEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *textEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *textEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *textEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *textEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *textEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *textEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferPool.Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

func (enc *textEncoder) AddReflected(key string, obj interface{}) error {
	enc.resetReflectBuf()
	err := enc.reflectEnc.Encode(obj)
	if err != nil {
		return err
	}
	enc.reflectBuf.TrimNewline()
	enc.addKey(key)
	_, err = enc.buf.Write(enc.reflectBuf.Bytes())
	return err
}

func (enc *textEncoder) OpenNamespace(key string) {
	// enc.addKey(key)
	// enc.buf.AppendByte('{')
	// enc.openNamespaces++
}

func (enc *textEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *textEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *textEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *textEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *textEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *textEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *textEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	// enc.buf.AppendByte('"')
	// enc.safeAddByteString(val)
	// enc.buf.AppendByte('"')
	enc.buf.Write(val)
}

func (enc *textEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	// enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	// enc.buf.AppendByte('"')
}

func (enc *textEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *textEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *textEncoder) AppendReflected(val interface{}) error {
	enc.resetReflectBuf()
	err := enc.reflectEnc.Encode(val)
	if err != nil {
		return err
	}
	enc.reflectBuf.TrimNewline()
	enc.addElementSeparator()
	_, err = enc.buf.Write(enc.reflectBuf.Bytes())
	return err
}

func (enc *textEncoder) AppendString(val string) {
	enc.addElementSeparator()
	// enc.buf.AppendByte('"')
	// enc.safeAddString(val)
	// enc.buf.AppendByte('"')
	enc.buf.AppendString(val)
}

func (enc *textEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *textEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *textEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *textEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *textEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *textEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *textEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *textEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *textEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *textEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *textEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *textEncoder) clone() *textEncoder {
	clone := getTextEncoder()
	clone.TextEncoderConfig = enc.TextEncoderConfig
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferPool.Get()
	return clone
}

func (enc *textEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	// final.buf.AppendByte('{')

	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
	}

	if final.LevelKey != "" {
		final.addKey(final.LevelKey)
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeLevel was a no-op. Fall back to strings to keep
			// output JSON valid.
			final.AppendString(ent.Level.String())
		}
	}

	if ent.LoggerName != "" && final.NameKey != "" {
		final.addKey(final.NameKey)
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output JSON valid.
			final.AppendString(ent.LoggerName)
		}
	}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.addKey(final.CallerKey)
		cur := final.buf.Len()
		final.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep output JSON valid.
			final.AppendString(ent.Caller.String())
		}
	}
	if final.MessageKey != "" {
		final.addKey(enc.MessageKey)
		final.AppendString(ent.Message)
	}
	if enc.buf.Len() > 0 {
		final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes())
	}
	addFields(final, final, fields)
	final.closeOpenNamespaces()
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}
	// final.buf.AppendByte('}')
	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(zapcore.DefaultLineEnding)
	}

	ret := final.buf
	putTextEncoder(final)
	return ret, nil
}

func (enc *textEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *textEncoder) closeOpenNamespaces() {
	// for i := 0; i < enc.openNamespaces; i++ {
	// 	enc.buf.AppendByte('}')
	// }
}

func (enc *textEncoder) addKey(key string) {
}

func (enc *textEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	enc.buf.AppendString(enc.FieldSeparator)
}

func (enc *textEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *textEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *textEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *textEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *textEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addFields(textEnc *textEncoder, enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
