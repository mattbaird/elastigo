package cluster

import (
	"github.com/mattbaird/elastigo/api"
)

// http://www.elasticsearch.org/guide/reference/api/admin-cluster-update-settings.html
func State(settingType string, key string, value int) error {
	url := "/_cluster/settings"
	m := map[string]map[string]int{settingType: map[string]int{key: value}}
	_, err := api.DoCommand("PUT", url, m)
	if err != nil {
		return err
	}
	return nil
}
