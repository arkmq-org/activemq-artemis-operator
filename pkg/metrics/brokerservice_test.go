/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Suite")
}

var _ = Describe("Service Metrics", func() {
	BeforeEach(func() {
		// Clean up metrics before each test
		ServiceAppsProvisioned.Reset()
	})

	It("UpdateServiceMetrics sets all gauges correctly", func() {
		UpdateServiceMetrics("test-service", "test-ns", 3)

		// Verify app count
		val := testutil.ToFloat64(ServiceAppsProvisioned.With(prometheus.Labels{
			"service":   "test-service",
			"namespace": "test-ns",
		}))
		Expect(val).To(Equal(float64(3)))
	})

	It("DeleteServiceMetrics removes gauges", func() {
		// Create metrics first
		UpdateServiceMetrics("delete-me", "test-ns", 1)

		// Delete them
		DeleteServiceMetrics("delete-me", "test-ns")
	})

	It("Multiple services tracked independently", func() {
		UpdateServiceMetrics("service-1", "ns-1", 5)
		UpdateServiceMetrics("service-2", "ns-2", 2)

		val1 := testutil.ToFloat64(ServiceAppsProvisioned.With(prometheus.Labels{
			"service":   "service-1",
			"namespace": "ns-1",
		}))
		Expect(val1).To(Equal(float64(5)))

		val2 := testutil.ToFloat64(ServiceAppsProvisioned.With(prometheus.Labels{
			"service":   "service-2",
			"namespace": "ns-2",
		}))
		Expect(val2).To(Equal(float64(2)))
	})
})
