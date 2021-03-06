/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package builds

import (
	"fmt"

	g "github.com/onsi/ginkgo"
	o "github.com/onsi/gomega"

	exutil "github.com/openshift/origin/test/extended/util"
)

var _ = g.Describe("[Feature:Builds][Conformance] build without output image", func() {
	defer g.GinkgoRecover()
	var (
		dockerImageFixture = exutil.FixturePath("testdata", "builds", "test-docker-no-outputname.json")
		s2iImageFixture    = exutil.FixturePath("testdata", "builds", "test-s2i-no-outputname.json")
		oc                 = exutil.NewCLI("build-no-outputname", exutil.KubeConfigPath())
	)

	g.Context("", func() {

		g.BeforeEach(func() {
			exutil.DumpDockerInfo()
		})

		g.AfterEach(func() {
			if g.CurrentGinkgoTestDescription().Failed {
				exutil.DumpPodStates(oc)
				exutil.DumpPodLogsStartingWith("", oc)
			}
		})

		g.Describe("building from templates", func() {
			oc.SetOutputDir(exutil.TestContext.OutputDir)

			g.It(fmt.Sprintf("should create an image from a docker template without an output image reference defined"), func() {
				err := oc.Run("create").Args("-f", dockerImageFixture).Execute()
				o.Expect(err).NotTo(o.HaveOccurred())

				g.By("expecting build to pass without an output image reference specified")
				br, err := exutil.StartBuildAndWait(oc, "test-docker")
				br.AssertSuccess()

				g.By("verifying the build test-docker-1 output")
				buildLog, err := br.Logs()
				fmt.Fprintf(g.GinkgoWriter, "\nBuild log:\n%s\n", buildLog)
				o.Expect(err).NotTo(o.HaveOccurred())
				o.Expect(buildLog).Should(o.ContainSubstring(`Build complete, no image push requested`))
			})

			g.It(fmt.Sprintf("should create an image from a S2i template without an output image reference defined"), func() {
				err := oc.Run("create").Args("-f", s2iImageFixture).Execute()
				o.Expect(err).NotTo(o.HaveOccurred())

				g.By("expecting build to pass without an output image reference specified")
				br, err := exutil.StartBuildAndWait(oc, "test-sti")
				o.Expect(err).NotTo(o.HaveOccurred())
				br.AssertSuccess()

				g.By("verifying the build test-sti-1 output")
				buildLog, err := br.Logs()
				fmt.Fprintf(g.GinkgoWriter, "\nBuild log:\n%s\n", buildLog)
				o.Expect(err).NotTo(o.HaveOccurred())

				o.Expect(buildLog).Should(o.ContainSubstring(`Build complete, no image push requested`))
			})
		})
	})
})
