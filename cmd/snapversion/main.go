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

package main

import (
	"github.com/ubuntu-core/snapversion/pkg/server"
	"github.com/zenazn/goji"
)

const (
	urlPattern    = "/:name/:channel/:arch"
	defaultSource = "https://search.apps.ubuntu.com"
)

var (
	// dependency aliasing
	gojiGet   = goji.Get
	gojiServe = goji.Serve
	srv       = &server.Server{}
)

func main() {
	srv.Source = defaultSource
	gojiGet(urlPattern, srv.Get)
	gojiServe()
}
