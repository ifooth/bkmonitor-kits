/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云-监控平台 (Blueking - Monitor) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package logger

import (
	"sync"
	"time"

	"github.com/go-logfmt/logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var (
	_logfmtPool = sync.Pool{New: func() interface{} {
		var enc logfmtEncoder
		enc.Encoder = logfmt.NewEncoder(enc.buf)
		return &enc
	}}

	bufferpool = buffer.NewPool()
)

func getEncoder() *logfmtEncoder {
	return _logfmtPool.Get().(*logfmtEncoder)
}

func putEncoder(enc *logfmtEncoder) {
	enc.EncoderConfig = nil
	enc.buf = nil
	_logfmtPool.Put(enc)
}

type logfmtEncoder struct {
	*zapcore.EncoderConfig
	Encoder    *logfmt.Encoder
	buf        *buffer.Buffer
	namespaces []string
}

func NewLogfmtEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	enc := &logfmtEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
	}
	enc.Encoder = logfmt.NewEncoder(enc.buf)

	return enc
}

func (enc *logfmtEncoder) Reset() {
	enc.Encoder.Reset()
	enc.buf.Reset()
	enc.namespaces = nil
}

func (enc *logfmtEncoder) AddArray(k string, marshaler zapcore.ArrayMarshaler) error {
	return enc.Encoder.EncodeKeyval(k, marshaler)
}
func (enc *logfmtEncoder) AddObject(k string, marshaler zapcore.ObjectMarshaler) error {
	return enc.Encoder.EncodeKeyval(k, marshaler)
}

func (enc *logfmtEncoder) AddReflected(k string, value interface{}) error {
	return enc.Encoder.EncodeKeyval(k, value)
}

func (enc *logfmtEncoder) OpenNamespace(key string) {
	enc.namespaces = append(enc.namespaces, key)
}

func (enc *logfmtEncoder) AddBinary(k string, v []byte)          { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddByteString(k string, v []byte)      { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddBool(k string, v bool)              { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddComplex128(k string, v complex128)  { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddComplex64(k string, v complex64)    { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddDuration(k string, v time.Duration) { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddFloat64(k string, v float64)        { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddFloat32(k string, v float32)        { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddInt(k string, v int)                { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddInt64(k string, v int64)            { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddInt32(k string, v int32)            { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddInt16(k string, v int16)            { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddInt8(k string, v int8)              { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddString(k, v string)                 { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddTime(k string, v time.Time)         { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddUint(k string, v uint)              { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddUint64(k string, v uint64)          { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddUint32(k string, v uint32)          { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddUint16(k string, v uint16)          { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddUint8(k string, v uint8)            { enc.Encoder.EncodeKeyval(k, v) }
func (enc *logfmtEncoder) AddUintptr(k string, v uintptr)        { enc.Encoder.EncodeKeyval(k, v) }

func (enc *logfmtEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *logfmtEncoder) clone() *logfmtEncoder {
	clone := getEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.buf = bufferpool.Get()
	clone.Encoder = logfmt.NewEncoder(clone.buf)
	clone.namespaces = enc.namespaces
	return clone
}

func (enc *logfmtEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()

	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
	}

	if final.LevelKey != "" {
		if err := final.Encoder.EncodeKeyval(final.LevelKey, ent.Level); err != nil {
			return nil, err
		}
	}

	if ent.Caller.Defined {
		if err := final.Encoder.EncodeKeyval(final.CallerKey, ent.Caller.TrimmedPath()); err != nil {
			return nil, err
		}
	}

	addFields(final, fields)

	// add endline
	if err := final.Encoder.EndRecord(); err != nil {
		return nil, err
	}

	ret := final.buf
	putEncoder(final)
	return ret, nil
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func init() {
	zap.RegisterEncoder("logfmt", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		enc := NewLogfmtEncoder(cfg)
		return enc, nil
	})
}
