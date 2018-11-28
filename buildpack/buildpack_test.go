/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package buildpack_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	buildpackPkg "github.com/buildpack/libbuildpack/buildpack"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", testBuildpack, spec.Report(report.Terminal{}))
}

func testBuildpack(t *testing.T, when spec.G, it spec.S) {

	it("unmarshals default from buildpack.toml", func() {
		root := internal.ScratchDir(t, "buildpack")
		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"))()

		in := strings.NewReader(`[buildpack]
id = "buildpack-id"
name = "buildpack-name"
version = "buildpack-version"

[[stacks]]
id = 'stack-id'
build-images = ["build-image-tag"]
run-images = ["run-image-tag"]

[metadata]
test-key = "test-value"
`)

		if err := internal.WriteToFile(in, filepath.Join(root, "buildpack.toml"), 0644); err != nil {
			t.Fatal(err)
		}

		buildpack, err := buildpackPkg.DefaultBuildpack(logger.Logger{})
		if err != nil {
			t.Fatal(err)
		}

		expected := buildpackPkg.Buildpack{
			Info: buildpackPkg.Info{
				ID:      "buildpack-id",
				Name:    "buildpack-name",
				Version: "buildpack-version",
			},
			Stacks: []buildpackPkg.Stack{
				{
					ID:          "stack-id",
					BuildImages: []buildpackPkg.BuildImages{"build-image-tag"},
					RunImages:   []buildpackPkg.RunImages{"run-image-tag"},
				},
			},
			Metadata: buildpackPkg.Metadata{
				"test-key": "test-value",
			},
			Root: root,
		}

		if !reflect.DeepEqual(buildpack, expected) {
			t.Errorf("Buildpack = %s, wanted %s", buildpack, expected)
		}
	})
}
