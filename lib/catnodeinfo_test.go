package elastigo

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCatNode(t *testing.T) {

	c := NewTestConn()

	Convey("Basic cat nodes", t, func() {

		fields := []string{"fm", "fe", "fcm", "fce", "ft", "ftt", "im", "rp", "n"}
		catNodes, err := c.GetCatNodeInfo(fields)

		So(err, ShouldBeNil)
		So(catNodes, ShouldNotBeNil)
		So(len(catNodes), ShouldBeGreaterThan, 0)


		for _, catNode := range catNodes {
			So(catNode.FieldMem, ShouldNotBeEmpty)
			So(catNode.FiltMem, ShouldNotBeEmpty)
			So(catNode.IDCacheMemory, ShouldNotBeEmpty)
			So(catNode.RamPerc, ShouldNotBeEmpty)
			So(catNode.Name, ShouldNotBeEmpty)
		}
	})

	Convey("Cat nodes with default arguments", t, func() {

		fields := []string{}
		catNodes, err := c.GetCatNodeInfo(fields)

		So(err, ShouldBeNil)
		So(catNodes, ShouldNotBeNil)
		So(len(catNodes), ShouldBeGreaterThan, 0)

		for _, catNode := range catNodes {
			So(catNode.Host, ShouldNotBeEmpty)
			So(catNode.IP, ShouldNotBeEmpty)
			So(catNode.NodeRole, ShouldNotBeEmpty)
			So(catNode.Name, ShouldNotBeEmpty)
		}
	})

	Convey("Cat nodes with all output fields", t, func() {

		fields := []string{
			"id", "pid", "h", "i", "po", "v", "b", "j", "d", "hc", "hp", "hm",
			"rc", "rp", "rm", "fdc", "fdp", "fdm", "l", "u", "r", "m", "n",
			"cs", "fm", "fe", "fcm", "fce", "ft", "ftt", "gc", "gti", "gto",
			"geti", "geto", "gmti", "gmto", "im", "idc", "idti", "idto", "iic",
			"iiti", "iito", "mc", "mcd", "mcs", "mt", "mtd", "mts", "mtt",
			"pc", "pm", "pq", "pti", "pto", "rto", "rti", "sfc", "sfti", "sfto",
			"so", "sqc", "sqti", "sqto", "sc", "sm", "siwm", "siwmx", "svmm",
		}
		catNodes, err := c.GetCatNodeInfo(fields)

		So(err, ShouldBeNil)
		So(catNodes, ShouldNotBeNil)
		So(len(catNodes), ShouldBeGreaterThan, 0)

		for _, catNode := range catNodes {
			So(catNode.Host, ShouldNotBeEmpty)
			So(catNode.IP, ShouldNotBeEmpty)
			So(catNode.NodeRole, ShouldNotBeEmpty)
			So(catNode.Name, ShouldNotBeEmpty)
			So(catNode.MergTotalSize, ShouldNotBeEmpty)
			So(catNode.GetMissingTime, ShouldNotBeEmpty)
			So(catNode.SegIdxWriterMem, ShouldNotBeEmpty)
		}
	})

	Convey("Invalid field error behavior", t, func() {

		fields := []string{"fm", "bogus"}
		catNodes, err := c.GetCatNodeInfo(fields)
		
		So(err, ShouldNotBeNil)

		for _, catNode := range catNodes {
			So(catNode.FieldMem, ShouldNotBeEmpty)
		}
	})
}
