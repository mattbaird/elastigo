package fixedwidth

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFixedWidth(t *testing.T) {

	Convey("Analyze basic header line", t, func() {
		// 0----*----1----*----2----*----3----*----4----*----5----*----6
		fields := parseHeaderLine("field1    field2  field3")
		So(len(fields), ShouldEqual, 3)
		So(fields[0].Header, ShouldEqual, "field1")
		So(fields[0].Width, ShouldEqual, 10)
		So(fields[1].Header, ShouldEqual, "field2")
		So(fields[1].Width, ShouldEqual, 8)
		So(fields[2].Header, ShouldEqual, "field3")
		So(fields[2].Width, ShouldEqual, -1)

		tab, err := parseDataLines([]string{
			"data1     data2   data three",
			"data four data_5  data six ok",
			"data7     dataochodata niner",
		}, fields)
		So(err, ShouldBeNil)
		So(tab, ShouldNotBeNil)
		So(len(tab), ShouldEqual, 3)
		So(len(tab[0].Data), ShouldEqual, 3)
		So(tab.Width(), ShouldEqual, 3)
		So(tab.Height(), ShouldEqual, 3)
		So(tab.Item(1, 1), ShouldEqual, "data_5")
		So(tab.Item(2, 1), ShouldEqual, "dataocho")

		//		fmt.Println("Columns")
		//		for _, c := range tab {
		//			fmt.Println(c.String())
		//		}

		//		fmt.Println("Table")
		//		fmt.Println(tab.String())

	})

	Convey("Test public API", t, func() {

		tabstr := "field1    f2 field3  field4\n" +
			"data1     d2 data3   data four\n" +
			"data fivesd6 datasevsdata 8888 888 88\n" +
			"data9     d10data_eledata 12"

		tab, err := NewFixedWidthTable(bytes.NewBufferString(tabstr).Bytes())

		So(err, ShouldBeNil)
		So(tab, ShouldNotBeNil)
		So(tab.Width(), ShouldEqual, 4)
		So(tab.Height(), ShouldEqual, 3)
		So(tab.Item(-1, -2), ShouldEqual, "")
		So(tab.Item(0, 0), ShouldEqual, "data1")
		So(tab.Item(2, 1), ShouldEqual, "d10")
		So(tab.Item(1, 3), ShouldEqual, "data 8888 888 88")
		So(tab.Item(2, 3), ShouldEqual, "data 12")
		So(tab.Item(3, 3), ShouldEqual, "")

		m := tab.RowMap(1)
		So(len(m), ShouldEqual, 4)
		So(m["field1"], ShouldEqual, "data fives")
		So(m["f2"], ShouldEqual, "d6")

		m = tab.RowMap(3)
		So(len(m), ShouldEqual, 0)
		
		So(tab.String(), ShouldEqual, tabstr)
	})

}
