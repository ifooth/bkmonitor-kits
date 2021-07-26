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
package gokit

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bkmonitor-kits/logger"
)

func TestGoKitLogger(t *testing.T) {
	l := logger.New(logger.Options{Format: "logfmt", Level: "debug", Stdout: true})
	kitLog := NewLogger(l)

	assert.Empty(t, kitLog.Log("hello", "world"))

	assert.Empty(t, level.Debug(kitLog).Log("msg", "debug_msg", "missing"))
	assert.Empty(t, level.Info(kitLog).Log("msg", "exiting"))

	WithValKitLog := log.With(kitLog, "componnet", "api")
	WithValKitLog.Log("msg", "world")
	assert.Empty(t, level.Warn(WithValKitLog).Log("msg", "debug_msg", "missing"))
	assert.Empty(t, level.Info(WithValKitLog).Log("msg", "exiting"))
}
