package elastigo

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/mattbaird/elastigo/fixedwidth"
)

// newCatNodeInfo returns an instance of CatNodeInfo populated with the
// the information in the cat output node line, which is passed in as a
// map of field name to value. An err is returned if a field is not known.
func newCatNodeInfo(data map[string]string) (catNode *CatNodeInfo, err error) {

	catNode = &CatNodeInfo{}

	// Populate the apropriate field in CatNodeInfo
	for field, value := range data {
		switch field {
		case "id", "nodeId":
			catNode.Id = value
		case "pid", "p":
			i, _ := strconv.Atoi(value)
			catNode.PID = int32(i)
		case "host", "h":
			catNode.Host = value
		case "ip", "i":
			catNode.IP = value
		case "port", "po":
			i, _ := strconv.Atoi(value)
			catNode.Port = int16(i)
		case "version", "v":
			catNode.Version = value
		case "build", "b":
			catNode.Build = value
		case "jdk", "j":
			catNode.JDK = value
		case "disk.avail", "d", "disk", "diskAvail":
			catNode.DiskAvail = value
		case "heap.current", "hc", "heapCurrent":
			catNode.HeapCur = value
		case "heap.percent", "hp", "heapPercent":
			i, _ := strconv.Atoi(value)
			catNode.HeapPerc = int16(i)
		case "heap.max", "hm", "heapMax":
			catNode.HeapMax = value
		case "ram.current", "rc", "ramCurrent":
			catNode.RamCur = value
		case "ram.percent", "rp", "ramPercent":
			i, _ := strconv.Atoi(value)
			catNode.RamPerc = int16(i)
		case "ram.max", "rm", "ramMax":
			catNode.RamMax = value
		case "file_desc.current", "fdc", "fileDescriptorCurrent":
			catNode.FileDescCur = value
		case "file_desc.percent", "fdp", "fileDescriptorPercent":
			i, _ := strconv.Atoi(value)
			catNode.FileDescPerc = int16(i)
		case "file_desc.max", "fdm", "fileDescriptorMax":
			catNode.FileDescMax = value
		case "load", "l":
			catNode.Load = value
		case "uptime", "u":
			catNode.UpTime = value
		case "node.role", "r", "role", "dc", "nodeRole":
			catNode.NodeRole = value
		case "master", "m":
			catNode.Master = value
		case "name", "n":
			catNode.Name = value
		case "completion.size", "cs", "completionSize":
			catNode.CmpltSize = value
		case "fielddata.memory_size", "fm", "fielddataMemory":
			catNode.FieldMem = value
		case "fielddata.evictions", "fe", "fieldataEvictions":
			i, _ := strconv.Atoi(value)
			catNode.FieldEvict = int32(i)
		case "filter_cache.memory_size", "fcm", "filterCacheMemory":
			catNode.FiltMem = value
		case "filter_cache.evictions", "fce", "filterCacheEvictions":
			i, _ := strconv.Atoi(value)
			catNode.FiltEvict = int32(i)
		case "flush.total", "ft", "flushTotal":
			i, _ := strconv.Atoi(value)
			catNode.FlushTotal = int32(i)
		case "flush.total_time", "ftt", "flushTotalTime":
			catNode.FlushTotalTime = value
		case "get.current", "gc", "getCurrent":
			i, _ := strconv.Atoi(value)
			catNode.GetCur = int32(i)
		case "get.time", "gti", "getTime":
			catNode.GetTime = value
		case "get.total", "gto", "getTotal":
			i, _ := strconv.Atoi(value)
			catNode.GetTotal = int32(i)
		case "get.exists_time", "geti", "getExistsTime":
			catNode.GetExistsTime = value
		case "get.exists_total", "geto", "getExistsTotal":
			i, _ := strconv.Atoi(value)
			catNode.GetExistsTotal = int32(i)
		case "get.missing_time", "gmti", "getMissingTime":
			catNode.GetMissingTime = value
		case "get.missing_total", "gmto", "getMissingTotal":
			i, _ := strconv.Atoi(value)
			catNode.GetMissingTotal = int32(i)
		case "id_cache.memory_size", "im", "idCacheMemory":
			catNode.IDCacheMemory = value
		case "indexing.delete_current", "idc", "indexingDeleteCurrent":
			i, _ := strconv.Atoi(value)
			catNode.IdxDelCur = int32(i)
		case "indexing.delete_time", "idti", "indexingDeleteime":
			catNode.IdxDelTime = value
		case "indexing.delete_total", "idto", "indexingDeleteTotal":
			i, _ := strconv.Atoi(value)
			catNode.IdxDelTotal = int32(i)
		case "indexing.index_current", "iic", "indexingIndexCurrent":
			i, _ := strconv.Atoi(value)
			catNode.IdxIdxCur = int32(i)
		case "indexing.index_time", "iiti", "indexingIndexTime":
			catNode.IdxIdxTime = value
		case "indexing.index_total", "iito", "indexingIndexTotal":
			i, _ := strconv.Atoi(value)
			catNode.IdxIdxTotal = int32(i)
		case "merges.current", "mc", "mergesCurrent":
			i, _ := strconv.Atoi(value)
			catNode.MergCur = int32(i)
		case "merges.current_docs", "mcd", "mergesCurrentDocs":
			i, _ := strconv.Atoi(value)
			catNode.MergCurDocs = int32(i)
		case "merges.current_size", "mcs", "mergesCurrentSize":
			catNode.MergCurSize = value
		case "merges.total", "mt", "mergesTotal":
			i, _ := strconv.Atoi(value)
			catNode.MergTotal = int32(i)
		case "merges.total_docs", "mtd", "mergesTotalDocs":
			i, _ := strconv.Atoi(value)
			catNode.MergTotalDocs = int32(i)
		case "merges.total_size", "mts", "mergesTotalSize":
			catNode.MergTotalSize = value
		case "merges.total_time", "mtt", "mergesTotalTime":
			catNode.MergTotalTime = value
		case "percolate.current", "pc", "percolateCurrent":
			i, _ := strconv.Atoi(value)
			catNode.PercCur = int32(i)
		case "percolate.memory_size", "pm", "percolateMemory":
			catNode.PercMem = value
		case "percolate.queries", "pq", "percolateQueries":
			i, _ := strconv.Atoi(value)
			catNode.PercQueries = int32(i)
		case "percolate.time", "pti", "percolateTime":
			catNode.PercTime = value
		case "percolate.total", "pto", "percolateTotal":
			i, _ := strconv.Atoi(value)
			catNode.PercTotal = int32(i)
		case "refesh.total", "rto", "refreshTotal":
			i, _ := strconv.Atoi(value)
			catNode.RefreshTotal = int32(i)
		case "refresh.time", "rti", "refreshTime":
			catNode.RefreshTime = value
		case "search.fetch_current", "sfc", "searchFetchCurrent":
			i, _ := strconv.Atoi(value)
			catNode.SearchFetchCur = int32(i)
		case "search.fetch_time", "sfti", "searchFetchTime":
			catNode.SearchFetchTime = value
		case "search.fetch_total", "sfto", "searchFetchTotal":
			i, _ := strconv.Atoi(value)
			catNode.SearchFetchTotal = int32(i)
		case "search.open_contexts", "so", "searchOpenContexts":
			i, _ := strconv.Atoi(value)
			catNode.SearchOpenContexts = int32(i)
		case "search.query_current", "sqc", "searchQueryCurrent":
			i, _ := strconv.Atoi(value)
			catNode.SearchQueryCur = int32(i)
		case "search.query_time", "sqti", "searchQueryTime":
			catNode.SearchQueryTime = value
		case "search.query_total", "sqto", "searchQueryTotal":
			i, _ := strconv.Atoi(value)
			catNode.SearchQueryTotal = int32(i)
		case "segments.count", "sc", "segmentsCount":
			i, _ := strconv.Atoi(value)
			catNode.SegCount = int32(i)
		case "segments.memory", "sm", "segmentsMemory":
			catNode.SegMem = value
		case "segments.index_writer_memory", "siwm", "segmentsIndexWriterMemory":
			catNode.SegIdxWriterMem = value
		case "segments.index_writer_max_memory", "siwmx", "segmentsIndexWriterMaxMemory":
			catNode.SegIdxWriterMax = value
		case "segments.version_map_memory", "svmm", "segmentsVersionMapMemory":
			catNode.SegVerMapMem = value
		default:
			return nil, fmt.Errorf("Invalid cat nodes field: %s", field)
		}
	}

	return catNode, nil
}

// GetCatNodeInfo issues an elasticsearch cat nodes request with the specified
// fields and returns a list of CatNodeInfos, one for each node, whose requested
// members are populated with statistics. If fields is nil or empty, the default
// cat output is used.
// NOTE: if you include the name field, make sure it is the last field in the
// list, because name values can contain spaces which screw up the parsing
func (c *Conn) GetCatNodeInfo(fields []string) (catNodes []CatNodeInfo, err error) {

	catNodes = make([]CatNodeInfo, 0)

	// If no fields have been specified, use the "default" arrangement
	if len(fields) < 1 {
		fields = []string{"host", "ip", "heap.percent", "ram.percent", "load",
			"node.role", "master", "name"}
	}

	// Issue a request for stats on the requested fields
	args := map[string]interface{}{
		"v": "",
		"bytes": "b",
		"h":     strings.Join(fields, ","),
	}
	indices, err := c.DoCommand("GET", "/_cat/nodes/", args, nil)
	if err != nil {
		return catNodes, err
	}

	// Create a table of response data
	tab, err := fixedwidth.NewFixedWidthTable(indices)
	for row := 0; row < tab.Height(); row++ {

		data := tab.RowMap(row)

		// Create a CatNodeInfo and append it to the result
		ci, err := newCatNodeInfo(data)
		if ci != nil {
			catNodes = append(catNodes, *ci)
		} else if err != nil {
			return catNodes, err
		}
	}
	return catNodes, nil
}
