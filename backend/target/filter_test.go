package filter

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	CONFIG_FILE = "target.conf.sample"
)

var (
	manager = New()
)

func init() {
	manager.AddFilterFunc(FilterTestFuncEquals, FilterTestEquals)
	manager.AddFilterFunc(FilterTestFuncTestInt, FilterTestInt)
}

func TestTarget(t *testing.T) {

	Convey("get target name's demensions", t, func() {
		Convey("when xxxx not in map", func() {
			Convey("will be not filtered", func() {
				filtered, e := manager.Compare(
					&FilterParams{Names: []string{"xxxx"}},
					nil, nil)
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)
			})
		})
		Convey("when FilterTestFuncEquals in map", func() {
			Convey("will get FilterTestFuncEquals' dememsions", func() {
				filtered, e := manager.Compare(
					&FilterParams{
						Names: []string{FilterTestFuncEquals}},
					nil, nil)
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)
			})
		})
	})

	Convey("test sequence filter", t, func() {
		Convey("when filter test_key does not equal", func() {
			Convey("will pass the filter", func() {
				filtered, e := manager.Compare(
					&FilterParams{
						Names: []string{FilterTestFuncEquals}},
					nil,
					FilterValues{})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)

				filtered, e = manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncEquals}},
					InputValues{testKeyEquals: ""},
					FilterValues{})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)

				filtered, e = manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncEquals}},
					InputValues{testKeyEquals: "1", testKeyInt: 0},
					FilterValues{testKeyEquals: "2"})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)

				filtered, e = manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncEquals}},
					InputValues{testKeyEquals: 1, testKeyInt: 2},
					FilterValues{testKeyEquals: 2})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)

				filtered, e = manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncEquals}},
					InputValues{testKeyEquals: 1, testKeyInt: 2},
					FilterValues{testKeyEquals: 2})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeFalse)
			})
		})
		Convey("when filter values' key exist", func() {
			Convey("will be filtered", func() {
				filtered, e := manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncEquals}},
					InputValues{},
					FilterValues{testKeyEquals: "1"})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeTrue)

				filtered, e = manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncEquals}},
					InputValues{testKeyEquals: "1"},
					FilterValues{testKeyEquals: "1"})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeTrue)
			})
		})
		Convey("when filter has two more names", func() {
			Convey("will be filtered", func() {
				filtered, e := manager.Compare(
					&FilterParams{Names: []string{
						FilterTestFuncEquals,
						FilterTestFuncTestInt}},
					InputValues{testKeyEquals: "2", testKeyInt: "1"},
					FilterValues{testKeyEquals: "1"})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeTrue)

				filtered, e = manager.Compare(
					&FilterParams{Names: []string{
						FilterTestFuncTestInt, FilterTestFuncEquals}},
					InputValues{testKeyEquals: "1", testKeyInt: 2},
					FilterValues{testKeyEquals: "1"})
				So(e, ShouldBeNil)
				So(filtered, ShouldBeTrue)
			})
		})
	})

	Convey("test sequence filter", t, func() {
		Convey("will be filtered", func() {
			filtered, e := manager.Compare(
				&FilterParams{
					Type: CompareTypeConsistent,
					Names: []string{
						FilterTestFuncEquals,
						FilterTestFuncTestInt}},
				InputValues{},
				FilterValues{testKeyEquals: "1"})
			So(e, ShouldBeNil)
			So(filtered, ShouldBeTrue)

			filtered, e = manager.Compare(
				&FilterParams{
					Type: CompareTypeConsistent,
					Names: []string{
						FilterTestFuncEquals,
						FilterTestFuncTestInt}},
				InputValues{testKeyEquals: "1"},
				FilterValues{testKeyEquals: "1"})
			So(e, ShouldBeNil)
			So(filtered, ShouldBeTrue)

			filtered, e = manager.Compare(
				&FilterParams{
					Type: CompareTypeConsistent,
					Names: []string{
						FilterTestFuncEquals,
						FilterTestFuncTestInt}},
				InputValues{testKeyEquals: 2, testKeyInt: 2},
				FilterValues{testKeyEquals: 2})
			So(e, ShouldBeNil)
			So(filtered, ShouldBeTrue)

			filtered, e = manager.Compare(
				&FilterParams{
					Type: CompareTypeConsistent,
					Names: []string{
						FilterTestFuncEquals,
						FilterTestFuncTestInt}},
				InputValues{testKeyEquals: 1, testKeyInt: "2"},
				FilterValues{testKeyEquals: 2})
			So(e, ShouldBeNil)
			So(filtered, ShouldBeTrue)
		})
	})

	Convey("test filter timeout", t, func() {
		Convey("set invalid timeout", func() {
			Convey("will return error", func() {

				e := manager.SetFilterTimeout(-1)
				So(e.Error(), ShouldEqual, ErrTimeoutMustAboveZero.New().Error())
			})
		})
		manager.AddFilterFunc(FilterTestFuncTimeout, FilterTestTimeout)
		manager.SetFilterTimeout(1)
		Convey("input timeout function and exec sequence function", func() {
			Convey("will timeout", func() {

				filtered, e := manager.Compare(
					&FilterParams{Names: []string{FilterTestFuncTimeout}},
					InputValues{testKeyEquals: 1, testKeyInt: "2"},
					FilterValues{testKeyEquals: 2})
				So(e.Error(), ShouldEqual, ErrFailedExecTimeout.New().Error())
				So(filtered, ShouldBeTrue)
			})
		})
		Convey("input timeout function and exec consistent function", func() {
			Convey("will timeout", func() {

				filtered, e := manager.Compare(
					&FilterParams{
						Type: CompareTypeConsistent,
						Names: []string{
							FilterTestFuncTimeout,
							FilterTestFuncEquals,
							FilterTestFuncTestInt}},
					InputValues{testKeyEquals: 1, testKeyInt: 2},
					FilterValues{testKeyEquals: 2})
				So(e.Error(), ShouldEqual, ErrFailedExecTimeout.New().Error())
				So(filtered, ShouldBeTrue)
			})
		})
	})
	return
}

// FilterTestFuncs
const (
	FilterTestFuncEquals  = "defaultTestEquals"
	FilterTestFuncTestInt = "defaultTestInt"
	FilterTestFuncTimeout = "defaultTestTimeout"

	testKeyEquals = "test_key_equals"
	testKeyInt    = "test_key_int"
)

func FilterTestEquals(req InputValues, cvs FilterValues) (bool, error) {
	if len(cvs) == 0 {
		return false, nil
	}

	if len(req) == 0 {
		return true, nil
	}

	return req[testKeyEquals] == cvs[testKeyEquals], nil
}

func FilterTestTimeout(_ InputValues, _ FilterValues) (bool, error) {
	time.Sleep(time.Second * 3)
	return false, nil
}

func FilterTestInt(req InputValues, _ FilterValues) (bool, error) {
	vs := req[testKeyInt]
	if vs == nil {
		return true, nil
	}

	switch reflect.TypeOf(vs).Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return false, nil
	}

	return true, nil
}
