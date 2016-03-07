// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zenazn/goji/web"
)

const (
	responseTpl = `{"whitelist_country_codes": [], "website": null, "last_updated": "2016-03-07T11:59:11.700043Z", "package_name": "ubuntu-core", "sequence": 42, "channels": ["edge"], "screenshot_url": null, "video_urls": [], "anon_download_url": "https://public.apps.ubuntu.com/anon/download/canonical/ubuntu-core.canonical/ubuntu-core.canonical_16.04.0-21+cowboy1_amd64.snap", "framework": [], "terms_of_service": "", "keywords": [], "id": 4142, "allow_unauthenticated": true, "support_url": "mailto:snappy-canonical-storeaccount@canonical.com", "icon_url": "https://myapps.developer.ubuntu.com/site_media/appmedia/2015/12/logo-ubuntu_cof-orange-hex_2.png", "binary_filesize": 81821696, "download_url": "https://public.apps.ubuntu.com/download/canonical/ubuntu-core.canonical/ubuntu-core.canonical_16.04.0-21+cowboy1_amd64.snap", "content": "os", "developer_name": "Shared snappy store account", "version": "%s", "_links": {"curies": [{"href": "https://wiki.ubuntu.com/AppStore/Interfaces/ClickPackageIndex#reltype_{rel}", "name": "clickindex", "templated": true}], "self": {"href": "https://search.apps.ubuntu.com/api/v1/package/ubuntu-core.canonical/stable"}}, "company_name": "Canonical", "channel": "edge", "department": ["accessories"], "screenshot_urls": [], "revision": 42, "status": "Published", "description": "The ubuntu-core OS snap\nno description", "click_framework": [], "price": 0.0, "origin": "canonical", "blacklist_country_codes": [], "date_published": "2015-12-17T15:51:40.296964Z", "alias": "ubuntu-core", "prices": {}, "icon_urls": {"256": "https://myapps.developer.ubuntu.com/site_media/appmedia/2015/12/logo-ubuntu_cof-orange-hex_2.png"}, "download_sha512": "5b27b46e37ad0aec23c4b1ceee2d0b5be7426996216ca0e632950039be4012f432e4587949a4b643018e479ad5dc278099031995be4f260d45d05b684b65a449", "publisher": "Canonical", "name": "ubuntu-core.canonical", "license": "Other Open Source", "changelog": "", "title": "ubuntu-core", "click_version": "0.1", "ratings_average": 0.0, "architecture": ["amd64"], "release": ["rolling-core"], "is_published": true}`
	testArch    = "myArch"
	testName    = "myName"
	testChannel = "myChannel"
	testVersion = "1.0"
)

type mockHandler struct {
	receivedName, receivedArch, receivedChannel string
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// we expect a URL with a Path of the form /api/v1/package/:name/:channel
	urlParts := strings.Split(r.URL.Path, "/")

	totalParts := len(urlParts)

	m.receivedName = urlParts[totalParts-2]
	m.receivedChannel = urlParts[totalParts-1]
	header := r.Header[archHeader]
	if len(header) > 0 {
		m.receivedArch = header[0]
	}

	responseText := fmt.Sprintf(responseTpl, testVersion)
	w.Write([]byte(responseText))
}

func TestOutgoingRequest(t *testing.T) {
	m := &mockHandler{}
	w := httptest.NewRecorder()

	setUpTest(t, m, w)

	if m.receivedName != testName {
		t.Fatalf("Expecting name %s, got %s", testName, m.receivedName)
	}

	if m.receivedChannel != testChannel {
		t.Fatalf("Expecting channel %s, got %s", testChannel, m.receivedChannel)
	}

	if m.receivedArch != testArch {
		t.Fatalf("Expecting arch %s, got %s", testArch, m.receivedArch)
	}
}

func TestResponse(t *testing.T) {
	m := &mockHandler{}
	w := httptest.NewRecorder()

	setUpTest(t, m, w)

	receivedVersion := w.Body.String()
	if receivedVersion != testVersion {
		t.Fatalf("Expecting version %s, got %s", testVersion, receivedVersion)
	}
}

func setUpTest(t *testing.T, m *mockHandler, w *httptest.ResponseRecorder) {
	sourceServer := httptest.NewServer(m)

	subject := &Server{
		Source: sourceServer.URL,
	}

	context := web.C{
		URLParams: map[string]string{
			"name":    testName,
			"channel": testChannel,
			"arch":    testArch},
	}

	req, err := http.NewRequest("GET", "http://localhost/myname/mychannel/myarch", nil)
	if err != nil {
		t.Fatal(err)
	}

	subject.Get(context, w, req)
}
