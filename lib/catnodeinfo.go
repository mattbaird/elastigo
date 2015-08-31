package elastigo

import (
	"fmt"
	"strconv"
	"strings"
)

// newCatNodeInfo returns an instance of CatNodeInfo populated with the
// the information in the cat output indexLine which contains the
// specified fields. An err is returned if a field is not known.
func newCatNodeInfo(fields []string, indexLine string) (catNode *CatNodeInfo, err error) {

	split := strings.Fields(indexLine)
	catNode = &CatNodeInfo{}

	// Check the fields length compared to the number of stats
	lf, ls := len(fields), len(split)
	if lf > ls {
		return nil, fmt.Errorf("Number of fields (%d) greater than number of stats (%d)", lf, ls)
	}

	// Populate the apropriate field in CatNodeInfo
	for i, field := range fields {
		switch field {
		case "id", "nodeId":
			catNode.Id = split[i]
		case "pid", "p":
			i, _ := strconv.Atoi(split[i])
			catNode.PID = int32(i)
		case "host", "h":
			catNode.Host = split[i]
		case "ip", "i":
			catNode.IP = split[i]
		case "port", "po":
			i, _ := strconv.Atoi(split[i])
			catNode.Port = int16(i)
		case "version", "v":
			catNode.Version = split[i]
		case "build", "b":
			catNode.Build = split[i]
		case "jdk", "j":
			catNode.JDK = split[i]
		case "disk.avail", "d", "disk", "diskAvail":
			catNode.DiskAvail = split[i]
		case "heap.current", "hc", "heapCurrent":
			catNode.HeapCur = split[i]
		case "heap.percent", "hp", "heapPercent":
			i, _ := strconv.Atoi(split[i])
			catNode.HeapPerc = int16(i)
		case "heap.max", "hm", "heapMax":
			catNode.HeapMax = split[i]
		case "ram.current", "rc", "ramCurrent":
			catNode.RamCur = split[i]
		case "ram.percent", "rp", "ramPercent":
			i, _ := strconv.Atoi(split[i])
			catNode.RamPerc = int16(i)
		case "ram.max", "rm", "ramMax":
			catNode.RamMax = split[i]
		case "file_desc.current", "fdc", "fileDescriptorCurrent":
			catNode.FileDescCur = split[i]
		case "file_desc.percent", "fdp", "fileDescriptorPercent":
			i, _ := strconv.Atoi(split[i])
			catNode.FileDescPerc = int16(i)
		case "file_desc.max", "fdm", "fileDescriptorMax":
			catNode.FileDescMax = split[i]
		case "load", "l":
			catNode.Load = split[i]
		case "uptime", "u":
			catNode.UpTime = split[i]
		case "node.role", "r", "role", "dc", "nodeRole":
			catNode.NodeRole = split[i]
		case "master", "m":
			catNode.Master = split[i]
		case "name", "n":
			catNode.Name = split[i]
		case "completion.size", "cs", "completionSize":
			catNode.CmpltSize = split[i]
		case "fielddata.memory_size", "fm", "fielddataMemory":
			catNode.FieldMem = split[i]
		case "fielddata.evictions", "fe", "fieldataEvictions":
			i, _ := strconv.Atoi(split[i])
			catNode.FieldEvict = int32(i)
		case "filter_cache.memory_size", "fcm", "filterCacheMemory":
			catNode.FiltMem = split[i]
		case "filter_cache.evictions", "fce", "filterCacheEvictions":
			i, _ := strconv.Atoi(split[i])
			catNode.FiltEvict = int32(i)
		case "flush.total", "ft", "flushTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.FlushTotal = int32(i)
		case "flush.total_time", "ftt", "flushTotalTime":
			i, _ := strconv.Atoi(split[i])
			catNode.FlushTotalTime = int32(i)
		case "get.current", "gc", "getCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.GetCur = int32(i)
		case "get.time", "gti", "getTime":
			catNode.GetTime = split[i]
		case "get.total", "gto", "getTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.GetTotal = int32(i)
		case "get.exists_time", "geti", "getExistsTime":
			catNode.GetExistsTime = split[i]
		case "get.exists_total", "geto", "getExistsTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.GetExistsTotal = int32(i)
		case "get.missing_time", "gmti", "getMissingTime":
			catNode.GetMissingTime = split[i]
		case "get.missing_total", "gmto", "getMissingTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.GetMissingTotal = int32(i)
		case "id_cache.memory_size", "im", "idCacheMemory":
			catNode.IDCacheMemory = split[i]
		case "indexing.delete_current", "idc", "indexingDeleteCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.IdxDelCur = int32(i)
		case "indexing.delete_time", "idti", "indexingDeleteime":
			catNode.IdxDelTime = split[i]
		case "indexing.delete_total", "idto", "indexingDeleteTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.IdxDelTotal = int32(i)
		case "indexing.index_current", "iic", "indexingIndexCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.IdxIdxCur = int32(i)
		case "indexing.index_time", "iiti", "indexingIndexTime":
			catNode.IdxIdxTime = split[i]
		case "indexing.index_total", "iito", "indexingIndexTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.IdxIdxTotal = int32(i)
		case "merges.current", "mc", "mergesCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.MergCur = int32(i)
		case "merges.current_docs", "mcd", "mergesCurrentDocs":
			i, _ := strconv.Atoi(split[i])
			catNode.MergCurDocs = int32(i)
		case "merges.current_size", "mcs", "mergesCurrentSize":
			catNode.MergCurSize = split[i]
		case "merges.total", "mt", "mergesTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.MergTotal = int32(i)
		case "merges.total_docs", "mtd", "mergesTotalDocs":
			i, _ := strconv.Atoi(split[i])
			catNode.MergTotalDocs = int32(i)
		case "merges.total_size", "mts", "mergesTotalSize":
			catNode.MergTotalSize = split[i]
		case "merges.total_time", "mtt", "mergesTotalTime":
			catNode.MergTotalTime = split[i]
		case "percolate.current", "pc", "percolateCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.PercCur = int32(i)
		case "percolate.memory_size", "pm", "percolateMemory":
			catNode.PercMem = split[i]
		case "percolate.queries", "pq", "percolateQueries":
			i, _ := strconv.Atoi(split[i])
			catNode.PercQueries = int32(i)
		case "percolate.time", "pti", "percolateTime":
			catNode.PercTime = split[i]
		case "percolate.total", "pto", "percolateTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.PercTotal = int32(i)
		case "refesh.total", "rto", "refreshTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.RefreshTotal = int32(i)
		case "refresh.time", "rti", "refreshTime":
			catNode.RefreshTime = split[i]
		case "search.fetch_current", "sfc", "searchFetchCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.SearchFetchCur = int32(i)
		case "search.fetch_time", "sfti", "searchFetchTime":
			catNode.SearchFetchTime = split[i]
		case "search.fetch_total", "sfto", "searchFetchTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.SearchFetchTotal = int32(i)
		case "search.open_contexts", "so", "searchOpenContexts":
			i, _ := strconv.Atoi(split[i])
			catNode.SearchOpenContexts = int32(i)
		case "search.query_current", "sqc", "searchQueryCurrent":
			i, _ := strconv.Atoi(split[i])
			catNode.SearchQueryCur = int32(i)
		case "search.query_time", "sqti", "searchQueryTime":
			catNode.SearchQueryTime = split[i]
		case "search.query_total", "sqto", "searchQueryTotal":
			i, _ := strconv.Atoi(split[i])
			catNode.SearchQueryTotal = int32(i)
		case "segments.count", "sc", "segmentsCount":
			i, _ := strconv.Atoi(split[i])
			catNode.SegCount = int32(i)
		case "segments.memory", "sm", "segmentsMemory":
			catNode.SegMem = split[i]
		case "segments.index_writer_memory", "siwm", "segmentsIndexWriterMemory":
			catNode.SegIdxWriterMem = split[i]
		case "segments.index_writer_max_memory", "siwmx", "segmentsIndexWriterMaxMemory":
			catNode.SegIdxWriterMax = split[i]
		case "segments.version_map_memory", "svmm", "segmentsVersionMapMemory":
			catNode.SegVerMapMem = split[i]
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
func (c *Conn) GetCatNodeInfo(fields []string) (catNodes []CatNodeInfo, err error) {
	catNodes = make([]CatNodeInfo, 0)

	// Issue a request for stats on the requested fields
	var args map[string]interface{}
	if len(fields) > 0 {
		args = map[string]interface{}{"h": strings.Join(fields, ",")}
	} else {
		fields = []string{"host", "ip", "heap.percent", "ram.percent", "load",
			"node.role", "master", "name"}
	}
	indices, err := c.DoCommand("GET", "/_cat/nodes/", args, nil)
	if err != nil {
		return catNodes, err
	}

	// Create a CatIndexInfo for each line in the response
	indexLines := strings.Split(string(indices[:]), "\n")
	for _, index := range indexLines {

		// Ignore empty output lines
		if len(index) < 1 {
			continue
		}

		// Create a CatNodeInfo and append it to the result
		ci, err := newCatNodeInfo(fields, index)
		if ci != nil {
			catNodes = append(catNodes, *ci)
		} else if err != nil {
			return catNodes, err
		}
	}
	return catNodes, nil
}
