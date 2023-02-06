// Copyright 2016-2018, Pulumi Corporation.
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

package deploytest

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
)

type ProgramFunc func(runInfo plugin.RunInfo, monitor *ResourceMonitor) error

func NewLanguageRuntime(program ProgramFunc, requiredPlugins ...workspace.PluginSpec) plugin.LanguageRuntime {
	return &languageRuntime{
		requiredPlugins: requiredPlugins,
		program:         program,
	}
}

type languageRuntime struct {
	requiredPlugins []workspace.PluginSpec
	program         ProgramFunc
}

func (p *languageRuntime) Close() error {
	return nil
}

func (p *languageRuntime) GetRequiredPlugins(info plugin.ProgInfo) ([]workspace.PluginSpec, error) {
	return p.requiredPlugins, nil
}

func (p *languageRuntime) Run(info plugin.RunInfo) (string, bool, error) {
	monitor, err := dialMonitor(context.Background(), info.MonitorAddress)
	if err != nil {
		return "", false, err
	}
	defer contract.IgnoreClose(monitor)

	// Run the program.
	done := make(chan error)
	go func() {
		done <- p.program(info, monitor)
	}()
	if progerr := <-done; progerr != nil {
		return progerr.Error(), false, nil
	}
	return "", false, nil
}

func (p *languageRuntime) GetPluginInfo() (workspace.PluginInfo, error) {
	return workspace.PluginInfo{Name: "TestLanguage"}, nil
}

func (p *languageRuntime) InstallDependencies(directory string) error {
	return nil
}

func (p *languageRuntime) About() (plugin.AboutInfo, error) {
	return plugin.AboutInfo{}, nil
}

func (p *languageRuntime) GetProgramDependencies(
	info plugin.ProgInfo, transitiveDependencies bool) ([]plugin.DependencyInfo, error) {
	return nil, nil
}

func (p *languageRuntime) RunPlugin(info plugin.RunPluginInfo) (io.Reader, io.Reader, context.CancelFunc, error) {
	return nil, nil, nil, fmt.Errorf("inline plugins are not currently supported")
}

func (p *languageRuntime) PublishPackage(ctx context.Context, artifactPath string, stderr *os.File) error {
	return nil
}

func (p *languageRuntime) PackPackage(ctx context.Context, packagePath, outPath string, stderr *os.File) (string, error) {
	return "", nil
}
