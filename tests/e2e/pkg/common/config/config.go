// Copyright © 2023 Cisco Systems, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type TestStrategy struct{}

const (
	TestStrategyMinimal          = "minimal"
	TestStrategyVersionComplete  = "version"
	TestStrategyProviderComplete = "provider"
	TestStrategyComplete         = "complete"
)

const (
	defaultReportDir               = "reports"
	defaultCreateTestReportFile    = "true"
	defaultMaxTimeout              = "1h"
	defaultAllowedOverrunDuration  = "10m"
	defaultTestStrategy            = TestStrategyVersionComplete
	defaultKubeConfigDirectoryPath = "kubeconfigs"
)

type TestsType struct {
	ReportDir               string
	CreateTestReportFile    string
	MaxTimeout              string
	AllowedOverrunDuration  string
	TestStrategy            string
	LabelFilter             string
	KubeConfigDirectoryPath string
}

func (t TestsType) String() string {
	return fmt.Sprintf(`
ReportDir: %s
CreateTestReportFile: %s
MaxTimeout: %s
AllowedOverrunDuration: %s
TestStrategy: %s
LabelFilter: %s
KubeConfigDirectoryPath: %s
`, viper.GetString(t.ReportDir), viper.GetString(t.CreateTestReportFile),
		viper.GetString(t.MaxTimeout), viper.GetString(t.AllowedOverrunDuration),
		viper.GetString(t.TestStrategy), viper.GetString(t.LabelFilter),
		viper.GetString(t.KubeConfigDirectoryPath))
}

var Tests = TestsType{
	CreateTestReportFile:    "tests.CreateTestReportFile",
	MaxTimeout:              "tests.MaxTimeout",
	AllowedOverrunDuration:  "tests.AllowedOverrunDuration",
	TestStrategy:            "tests.TestStrategy",
	ReportDir:               "tests.ReportDir",
	LabelFilter:             "tests.LabelFilter",
	KubeConfigDirectoryPath: "tests.KubeConfigDirectoryPath",
}

func init() {
	viper.AutomaticEnv()

	viper.BindEnv(Tests.CreateTestReportFile, "CREATE_TEST_REPORT_FILE")
	viper.SetDefault(Tests.CreateTestReportFile, defaultCreateTestReportFile)

	viper.BindEnv(Tests.ReportDir, "REPORT_DIR")
	viper.SetDefault(Tests.ReportDir, defaultReportDir)

	viper.BindEnv(Tests.MaxTimeout, "MAX_TIMEOUT")
	viper.SetDefault(Tests.MaxTimeout, defaultMaxTimeout)

	viper.BindEnv(Tests.AllowedOverrunDuration, "ALLOWED_OVERRUN_DURATION")
	viper.SetDefault(Tests.AllowedOverrunDuration, defaultAllowedOverrunDuration)

	viper.BindEnv(Tests.TestStrategy, "TEST_STRATEGY")
	viper.SetDefault(Tests.TestStrategy, defaultTestStrategy)

	viper.BindEnv(Tests.KubeConfigDirectoryPath, "KUBECONFIG_DIR")
	viper.SetDefault(Tests.KubeConfigDirectoryPath, defaultKubeConfigDirectoryPath)

	viper.BindEnv(Tests.LabelFilter, "LABEL_FILTER")
}