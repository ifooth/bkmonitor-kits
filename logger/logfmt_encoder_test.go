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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestEncoderObjectFields(t *testing.T) {
	tests := []struct {
		desc     string
		expected string
		f        func(zapcore.Encoder)
	}{
		{"binary", `k=ab12`, func(e zapcore.Encoder) { e.AddBinary("k", []byte("ab12")) }},
		{"bool", `k\=true`, func(e zapcore.Encoder) { e.AddBool(`k\`, true) }}, // test key escaping once
		{"bool", `k=true`, func(e zapcore.Encoder) { e.AddBool("k", true) }},
		{"bool", `k=false`, func(e zapcore.Encoder) { e.AddBool("k", false) }},
		{"byteString", `k=v\`, func(e zapcore.Encoder) { e.AddByteString(`k`, []byte(`v\`)) }},
	}

	for _, tt := range tests {
		assertOutput(t, tt.desc, tt.expected, tt.f)
	}
}

func assertOutput(t testing.TB, desc string, expected string, f func(zapcore.Encoder)) {
	encoderConfig := zapcore.EncoderConfig{
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	enc := NewLogfmtEncoder(encoderConfig)
	f(enc)
	l := enc.(*logfmtEncoder)
	assert.Equal(t, expected, l.buf.String(), "Unexpected encoder output after adding a %s.", desc)

	l.Reset()
}
