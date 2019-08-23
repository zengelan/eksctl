// +build integration

package integration

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestSuite(t *testing.T) {
	_ = flag.String("eksctl.path", "../eksctl", "Path to eksctl")

	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) All Suites", func() {
	var (
		suite   *IntergrationTests
		cluster *IntergrationCluster
	)

	BeforeSuite(func() {
		cluster = new(IntergrationCluster)
		suite = &IntergrationTests{
			Clusters: NewOnDemandClusterPool("us-west-2"),
		}
	})

	BeforeEach(func() {
		suite.UseCluster(cluster)
		Expect(cluster).ToNot(BeNil())
		suite.LogClusterPoolStats()
	})

	AfterSuite(func() {
		suite.FreeAllClusters()
	})

	It("TODO1", func() {
		Expect(cluster.Name).ToNot(BeEmpty())
	})

	Describe("MORE", func() {
		It("TODO2", func() {
			Expect(cluster.Name).ToNot(BeEmpty())
		})

		It("TODO3", func() {
			Expect(cluster.Name).ToNot(BeEmpty())
		})

		It("TODO4", func() {
			Expect(cluster.Name).ToNot(BeEmpty())
		})
		It("TODO5", func() {
			Expect(cluster.Name).ToNot(BeEmpty())
		})
		It("TODO6", func() {
			Expect(cluster.Name).ToNot(BeEmpty())
		})
		It("TODO7", func() {
			Expect(cluster.Name).ToNot(BeEmpty())
		})
	})

	It("TODO8", func() {
		Expect(cluster.Name).ToNot(BeEmpty())
	})

	It("TODO9", func() {
		Expect(cluster.Name).ToNot(BeEmpty())
	})
})
