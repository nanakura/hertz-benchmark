/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/cloudwego/hertz-benchmark/perf"
	"github.com/cloudwego/hertz-benchmark/runner"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/gofiber/fiber/v2"
)

const (
	port        = "0.0.0.0:8004"
	debugPort   = "0.0.0.0:18004"
	actionQuery = "action"
)

var recorder = perf.NewRecorder("Fiber@Server")

func main() {
	// start pprof server
	go func() {
		err := perf.ServeMonitor(debugPort)
		if err != nil {
			hlog.Error(err)
		}
	}()

	cfg := fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}
	r := fiber.New(cfg)

	r.Post("/", requestHandler)

	err := r.Listen(port)
	if err != nil {
		hlog.Error(err)
	}
}

func requestHandler(ctx *fiber.Ctx) error {
	runner.ProcessRequest(recorder, ctx.Query(actionQuery))

	ctx.Response().Header.SetContentType("text/plain; charset=utf8")
	return ctx.Send(ctx.Request().Body())
}
