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
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/TencentBlueKing/bkmonitor-kits/logger"
)

type Logger struct {
	base *logger.Logger
}

func NewLogger(l logger.Logger) Logger {
	return Logger{base: &l}

}

// Log Implementation https://github.com/go-kit/kit/blob/master/log/log.go#L10 interface
func (l Logger) Log(keyvals ...interface{}) error {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, log.ErrMissingValue)
	}

	var level string
	var msg string

	keyvals2 := make([]interface{}, 0, len(keyvals))

	for i := 0; i+2 <= len(keyvals); i += 2 {
		key := fmt.Sprintf("%s", keyvals[i])

		if key == "level" {
			level = fmt.Sprintf("%s", keyvals[i+1])
		} else if key == "msg" {
			msg = fmt.Sprintf("%s", keyvals[i+1])
		} else {
			keyvals2 = append(keyvals2, keyvals[i], keyvals[i+1])
		}
	}

	var logFn func(msg string, keyVal ...interface{})

	switch level {
	case "debug":
		logFn = l.base.Debugw
	case "info":
		logFn = l.base.Infow
	case "warn":
		logFn = l.base.Warnw
	case "error":
		logFn = l.base.Errorw
	default:
		logFn = l.base.Infow
	}

	logFn(msg, keyvals2...)

	return nil
}
