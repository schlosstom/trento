package discovery

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trento-project/trento/agent/discovery/collector"
	"github.com/trento-project/trento/internal/cluster"
)

const ClusterDiscoveryId string = "ha_cluster_discovery"
const ClusterDiscoveryMinPeriod time.Duration = 1 * time.Second

// This Discover handles any Pacemaker Cluster type
type ClusterDiscovery struct {
	id              string
	collectorClient collector.Client
	interval        time.Duration
}

func NewClusterDiscovery(collectorClient collector.Client, config DiscoveriesConfig) Discovery {
	d := ClusterDiscovery{}
	d.collectorClient = collectorClient
	d.id = ClusterDiscoveryId
	d.interval = config.DiscoveriesPeriodsConfig.Cluster

	return d
}

func (c ClusterDiscovery) GetId() string {
	return c.id
}

func (d ClusterDiscovery) GetInterval() time.Duration {
	return d.interval
}

// Execute one iteration of a discovery and publish the results to the collector
func (d ClusterDiscovery) Discover() (string, error) {
	cluster, err := cluster.NewCluster()
	if err != nil {
		return "No HA cluster discovered on this host", nil
	}

	err = d.collectorClient.Publish(d.id, cluster)
	if err != nil {
		log.Debugf("Error while sending cluster discovery to data collector: %s", err)
		return "", err
	}

	return fmt.Sprintf("Cluster with name: %s successfully discovered", cluster.Name), nil
}
