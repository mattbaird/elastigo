package elastigo

type CatIndexInfo struct {
	Health   string
	Status   string
	Name     string
	Shards   int
	Replicas int
	Docs     CatIndexDocs
	Store    CatIndexStore
}

type CatIndexDocs struct {
	Count   int64
	Deleted int64
}

type CatIndexStore struct {
	Size    int64
	PriSize int64
}

type CatShardInfo struct {
	IndexName string
	Shard     int
	Primary   string
	State     string
	Docs      int64
	Store     int64
	NodeIP    string
	NodeName  string
}

type CatNodeInfo struct {
	Id                 string
	PID                int32
	Host               string
	IP                 string
	Port               int16
	Version            string
	Build              string
	JDK                string
	DiskAvail          string
	HeapCur            string
	HeapPerc           int16
	HeapMax            string
	RamCur             string
	RamPerc            int16
	RamMax             string
	FileDescCur        string
	FileDescPerc       int16
	FileDescMax        string
	Load               string
	UpTime             string
	NodeRole           string
	Master             string
	Name               string
	CmpltSize          string
	FieldMem           string
	FieldEvict         int32
	FiltMem            string
	FiltEvict          int32
	FlushTotal         int32
	FlushTotalTime     string
	GetCur             int32
	GetTime            string
	GetTotal           int32
	GetExistsTime      string
	GetExistsTotal     int32
	GetMissingTime     string
	GetMissingTotal    int32
	IDCacheMemory      string
	IdxDelCur          int32
	IdxDelTime         string
	IdxDelTotal        int32
	IdxIdxCur          int32
	IdxIdxTime         string
	IdxIdxTotal        int32
	MergCur            int32
	MergCurDocs        int32
	MergCurSize        string
	MergTotal          int32
	MergTotalDocs      int32
	MergTotalSize      string
	MergTotalTime      string
	PercCur            int32
	PercMem            string
	PercQueries        int32
	PercTime           string
	PercTotal          int32
	RefreshTotal       int32
	RefreshTime        string
	SearchFetchCur     int32
	SearchFetchTime    string
	SearchFetchTotal   int32
	SearchOpenContexts int32
	SearchQueryCur     int32
	SearchQueryTime    string
	SearchQueryTotal   int32
	SegCount           int32
	SegMem             string
	SegIdxWriterMem    string
	SegIdxWriterMax    string
	SegVerMapMem       string
}
