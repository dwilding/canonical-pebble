// Copyright (c) 2014-2020 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"gopkg.in/check.v1"
)

type responseSuite struct{}

var _ = check.Suite(&responseSuite{})

func (s *responseSuite) TestRespSetsLocationIfAccepted(c *check.C) {
	rec := httptest.NewRecorder()

	rsp := &resp{
		Status: http.StatusAccepted,
		Result: map[string]any{
			"resource": "foo/bar",
		},
	}

	rsp.ServeHTTP(rec, nil)
	hdr := rec.Result().Header
	c.Check(hdr.Get("Location"), check.Equals, "foo/bar")
}

func (s *responseSuite) TestRespSetsLocationIfCreated(c *check.C) {
	rec := httptest.NewRecorder()

	rsp := &resp{
		Status: http.StatusCreated,
		Result: map[string]any{
			"resource": "foo/bar",
		},
	}

	rsp.ServeHTTP(rec, nil)
	hdr := rec.Result().Header
	c.Check(hdr.Get("Location"), check.Equals, "foo/bar")
}

func (s *responseSuite) TestRespDoesNotSetLocationIfOther(c *check.C) {
	rec := httptest.NewRecorder()

	rsp := &resp{
		Status: http.StatusTeapot, // I'm a teapot
		Result: map[string]any{
			"resource": "foo/bar",
		},
	}

	rsp.ServeHTTP(rec, nil)
	hdr := rec.Result().Header
	c.Check(hdr.Get("Location"), check.Equals, "")
}

func (s *responseSuite) TestFileResponseSetsContentDisposition(c *check.C) {
	const filename = "icon.png"

	path := filepath.Join(c.MkDir(), filename)
	err := os.WriteFile(path, nil, os.ModePerm)
	c.Check(err, check.IsNil)

	rec := httptest.NewRecorder()
	rsp := fileResponse(path)
	req, err := http.NewRequest("GET", "", nil)
	c.Check(err, check.IsNil)

	rsp.ServeHTTP(rec, req)

	hdr := rec.Result().Header
	c.Check(hdr.Get("Content-Disposition"), check.Equals,
		fmt.Sprintf("attachment; filename=%s", filename))
}

// This diverges from snapd. For historical reasons snapd must send a null result
// in this case, but there are no old clients to be worried about here.
func (s *responseSuite) TestRespJSONWithNullResult(c *check.C) {
	rj := &respJSON{Result: nil}
	data, err := json.Marshal(rj)
	c.Assert(err, check.IsNil)
	c.Check(string(data), check.Equals, `{"type":"","status-code":0}`)
}

func (s *responseSuite) TestErrorResponderPrintfsWithArgs(c *check.C) {
	teapot := makeErrorResponder(http.StatusTeapot)

	rec := httptest.NewRecorder()
	rsp := teapot("system memory below %d%%.", 1)
	req, err := http.NewRequest("GET", "", nil)
	c.Assert(err, check.IsNil)
	rsp.ServeHTTP(rec, req)

	var v struct{ Result errorResult }
	c.Assert(json.NewDecoder(rec.Body).Decode(&v), check.IsNil)

	c.Check(v.Result.Message, check.Equals, "system memory below 1%.")
}

func (s *responseSuite) TestErrorResponderDoesNotPrintfAlways(c *check.C) {
	teapot := makeErrorResponder(http.StatusTeapot)

	rec := httptest.NewRecorder()
	rsp := teapot("system memory below 1%.")
	req, err := http.NewRequest("GET", "", nil)
	c.Assert(err, check.IsNil)
	rsp.ServeHTTP(rec, req)

	var v struct{ Result errorResult }
	c.Assert(json.NewDecoder(rec.Body).Decode(&v), check.IsNil)

	c.Check(v.Result.Message, check.Equals, "system memory below 1%.")
}
