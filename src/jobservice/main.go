// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/goharbor/harbor/src/adminserver/client"
	"github.com/goharbor/harbor/src/jobservice/config"
	"github.com/goharbor/harbor/src/jobservice/env"
	"github.com/goharbor/harbor/src/jobservice/job/impl"
	"github.com/goharbor/harbor/src/jobservice/logger"
	"github.com/goharbor/harbor/src/jobservice/runtime"
	"github.com/goharbor/harbor/src/jobservice/utils"
)

func main() {
	// Get parameters
	configPath := flag.String("c", "", "Specify the yaml config file path")
	flag.Parse()

	// Missing config file
	if configPath == nil || utils.IsEmptyStr(*configPath) {
		flag.Usage()
		panic("no config file is specified")
	}

	// Load configurations
	if err := config.DefaultConfig.Load(*configPath, true); err != nil {
		panic(fmt.Sprintf("load configurations error: %s\n", err))
	}

	// Create the root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	if err := logger.Init(ctx); err != nil {
		panic(err)
	}

	// Set job context initializer
	runtime.JobService.SetJobContextInitializer(func(ctx *env.Context) (env.JobContext, error) {
		secret := config.GetAuthSecret()
		if utils.IsEmptyStr(secret) {
			return nil, errors.New("empty auth secret")
		}

		adminClient := client.NewClient(config.GetAdminServerEndpoint(), &client.Config{Secret: secret})
		jobCtx := impl.NewContext(ctx.SystemContext, adminClient)

		if err := jobCtx.Init(); err != nil {
			return nil, err
		}

		return jobCtx, nil
	})

	// Start
	runtime.JobService.LoadAndRun(ctx, cancel)
}
