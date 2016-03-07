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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zenazn/goji/web"
)

const (
	archHeader = "X-Ubuntu-Architecture"
)

// Server is the type that encapsulates the service functionality
type Server struct {
	Source string
}

// Get is the default Goji based http handler
func (s *Server) Get(c web.C, w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	url := s.Source + "/api/v1/package/" + c.URLParams["name"] + "/" + c.URLParams["channel"]

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
	}

	req.Header.Add(archHeader, c.URLParams["arch"])

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error requesting: ", err)
	}

	defer resp.Body.Close()

	version := getVersion(resp.Body)

	_, err = w.Write([]byte(version))
	if err != nil {
		fmt.Println("Error writing response: ", err)
	}
}

func getVersion(body io.ReadCloser) string {
	type snapInfo struct {
		Version string
	}

	var snap snapInfo
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&snap)
	if err != nil {
		fmt.Println("Error decoding json: ", err)
	}

	return snap.Version
}
