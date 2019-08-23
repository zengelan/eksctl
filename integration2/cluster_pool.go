package integration

import (
	"fmt"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega/gexec"

	"github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

type IntergrationTests struct {
	Clusters ClusterPool
	Commands map[string]runner.Cmd
}

type ClusterPoolStats struct {
	InUse, Available int
}

type ClusterPool interface {
	Get() *IntergrationCluster
	FreeAll()
	Stats() ClusterPoolStats
}

type IntergrationCluster struct {
	Name, Region   string
	KubeconfigPath string
}

func (t *IntergrationTests) UseCluster(useCluster *IntergrationCluster) {
	newCluster := t.Clusters.Get()
	*useCluster = *newCluster
}

func (t *IntergrationTests) LogClusterPoolStats() {
	logInfo("%#v", t.Clusters.Stats())
}

func (t *IntergrationTests) FreeAllClusters() {
	t.Clusters.FreeAll()
}

type OnDemandClusterPool struct {
	clusters []*IntergrationCluster
	mutex    *sync.Mutex
	region   string
}

func NewOnDemandClusterPool(region string) *OnDemandClusterPool {
	return &OnDemandClusterPool{
		clusters: []*IntergrationCluster{},
		mutex:    &sync.Mutex{},
		region:   region,
	}
}

func (p *OnDemandClusterPool) Get() *IntergrationCluster {
	cluster := &IntergrationCluster{
		Name:   cmdutils.ClusterName("", ""),
		Region: p.region,
	}
	logInfo("createging cluster %q...", cluster.Name)
	p.mutex.Lock()
	p.clusters = append(p.clusters, cluster)
	p.mutex.Unlock()

	Expect(true).To(BeTrue())
	return cluster
}

func (p *OnDemandClusterPool) FreeAll() {
	for _, cluster := range p.clusters {
		logInfo("deleting cluser %q...", cluster.Name)
	}
	gexec.KillAndWait()
}

func (p *OnDemandClusterPool) Stats() ClusterPoolStats {
	return ClusterPoolStats{
		InUse:     len(p.clusters),
		Available: -1,
	}
}

func logInfo(msgFmt string, args ...interface{}) {
	msgPrefix := fmt.Sprintf("\n[ginkonode:%d/%d]\t", config.GinkgoConfig.ParallelNode, config.GinkgoConfig.ParallelTotal)
	fmt.Fprintf(GinkgoWriter, msgPrefix+msgFmt+"\n", args...)
}
