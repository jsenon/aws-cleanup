package configs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConst(t *testing.T) {
	Convey("Test Constant", t, func() {
		serv := Service
		desc := Description
		So(serv, ShouldEqual, "aws-cleanup")
		So(desc, ShouldEqual, "Cleaning AWS")
		So(Version, ShouldHaveSameTypeAs, "asdf")
		So(GitCommit, ShouldHaveSameTypeAs, "asdf")
		So(BuildDate, ShouldHaveSameTypeAs, "asdf")
		So(BuildDate, ShouldNotHaveSameTypeAs, 123)
		So(GitCommit, ShouldNotHaveSameTypeAs, 123)
		So(Version, ShouldNotHaveSameTypeAs, 123)
	})
}
