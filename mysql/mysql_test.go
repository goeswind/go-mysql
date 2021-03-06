package mysql

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type mysqlTestSuite struct {
}

var _ = check.Suite(&mysqlTestSuite{})

func (s *mysqlTestSuite) SetUpSuite(c *check.C) {

}

func (s *mysqlTestSuite) TearDownSuite(c *check.C) {

}

func (t *mysqlTestSuite) TestMysqlGTIDInterval(c *check.C) {
	i, err := parseInterval("1-2")
	c.Assert(err, check.IsNil)
	c.Assert(i, check.DeepEquals, Interval{1, 3})

	i, err = parseInterval("1")
	c.Assert(err, check.IsNil)
	c.Assert(i, check.DeepEquals, Interval{1, 2})

	i, err = parseInterval("1-1")
	c.Assert(err, check.IsNil)
	c.Assert(i, check.DeepEquals, Interval{1, 2})

	i, err = parseInterval("1-2")
	c.Assert(err, check.IsNil)
}

func (t *mysqlTestSuite) TestMysqlGTIDIntervalSlice(c *check.C) {
	i := IntervalSlice{Interval{1, 2}, Interval{2, 4}, Interval{2, 3}}
	i.Sort()
	c.Assert(i, check.DeepEquals, IntervalSlice{Interval{1, 2}, Interval{2, 3}, Interval{2, 4}})
	n := i.Normalize()
	c.Assert(n, check.DeepEquals, IntervalSlice{Interval{1, 4}})

	i = IntervalSlice{Interval{1, 2}, Interval{3, 5}, Interval{1, 3}}
	i.Sort()
	c.Assert(i, check.DeepEquals, IntervalSlice{Interval{1, 2}, Interval{1, 3}, Interval{3, 5}})
	n = i.Normalize()
	c.Assert(n, check.DeepEquals, IntervalSlice{Interval{1, 5}})

	i = IntervalSlice{Interval{1, 2}, Interval{4, 5}, Interval{1, 3}}
	i.Sort()
	c.Assert(i, check.DeepEquals, IntervalSlice{Interval{1, 2}, Interval{1, 3}, Interval{4, 5}})
	n = i.Normalize()
	c.Assert(n, check.DeepEquals, IntervalSlice{Interval{1, 3}, Interval{4, 5}})

	n1 := IntervalSlice{Interval{1, 3}, Interval{4, 5}}
	n2 := IntervalSlice{Interval{1, 2}}

	c.Assert(n1.Contain(n2), check.Equals, true)
	c.Assert(n2.Contain(n1), check.Equals, false)

	n1 = IntervalSlice{Interval{1, 3}, Interval{4, 5}}
	n2 = IntervalSlice{Interval{1, 6}}

	c.Assert(n1.Contain(n2), check.Equals, false)
	c.Assert(n2.Contain(n1), check.Equals, true)
}

func (t *mysqlTestSuite) TestMysqlGTIDCodec(c *check.C) {
	us, err := ParseUUIDSet("de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2")
	c.Assert(err, check.IsNil)

	c.Assert(us.String(), check.Equals, "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2")

	buf := us.Encode()
	err = us.Decode(buf)
	c.Assert(err, check.IsNil)

	gs, err := ParseMysqlGTIDSet("de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2,de278ad0-2106-11e4-9f8e-6edd0ca20948:1-2")
	c.Assert(err, check.IsNil)

	buf = gs.Encode()
	o, err := DecodeMysqlGTIDSet(buf)
	c.Assert(err, check.IsNil)
	c.Assert(gs, check.DeepEquals, o)
}

func (t *mysqlTestSuite) TestMysqlGTIDContain(c *check.C) {
	g1, err := ParseMysqlGTIDSet("3E11FA47-71CA-11E1-9E33-C80AA9429562:23")
	c.Assert(err, check.IsNil)

	g2, err := ParseMysqlGTIDSet("3E11FA47-71CA-11E1-9E33-C80AA9429562:21-57")
	c.Assert(err, check.IsNil)

	c.Assert(g2.Contain(g1), check.Equals, true)
	c.Assert(g1.Contain(g2), check.Equals, false)
}
