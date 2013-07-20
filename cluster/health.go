package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/meanpath/elastigo/api"
	"strings"
)

// The cluster health API allows to get a very simple status on the health of the cluster.
// see http://www.elasticsearch.org/guide/reference/api/admin-cluster-health.html
// TODO: implement wait_for_status, timeout, wait_for_relocating_shards, wait_for_nodes
// TODO: implement level (Can be one of cluster, indices or shards. Controls the details level of the health 
// information returned. Defaults to cluster.)
func Health(pretty bool, indices ...string) (ClusterHealthResponse, error) {
	var url string
	var retval ClusterHealthResponse
	if len(indices) > 0 {
		url = fmt.Sprintf("/_cluster/health/%s?%s", strings.Join(indices, ","), api.Pretty(pretty))
	} else {
		url = fmt.Sprintf("/_cluster/health?%s", api.Pretty(pretty))
	}
	body, err := api.DoCommand("GET", url, nil)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	//fmt.Println(body)
	return retval, err
}

type ClusterHealthResponse struct {
	ClusterName         string `json:"cluster_name"`
	Status              string `json:"status"`
	TimedOut            bool   `json:"timed_out"`
	NumberOfNodes       int    `json:"number_of_nodes"`
	NumberOfDataNodes   int    `json:"number_of_data_nodes"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}
