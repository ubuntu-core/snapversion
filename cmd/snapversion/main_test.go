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
	"testing"

	"github.com/ubuntu-core/snapversion/pkg/server"
	"github.com/zenazn/goji/web"
)

func TestGojiSetUp(t *testing.T) {
	var receivedPattern web.PatternType
	var receivedHandler web.HandlerType
	var serveCalled bool

	fakeGet := func(pattern web.PatternType, handler web.HandlerType) {
		receivedPattern = pattern
		receivedHandler = handler
	}
	fakeServe := func() {
		serveCalled = true
	}

	backGojiGet := gojiGet
	defer func() { gojiGet = backGojiGet }()
	gojiGet = fakeGet
	backGojiServe := gojiServe
	defer func() { gojiServe = backGojiServe }()
	gojiServe = fakeServe
	backSrv := srv
	defer func() { srv = backSrv }()
	srv = &server.Server{}

	main()

	if srv.Source != defaultSource {
		t.Fatalf("Expecting %s server source, got %s", defaultSource, srv.Source)
	}

	if receivedPattern != urlPattern {
		t.Fatalf("Expecting %s pattern for goji.Get, got %s", urlPattern, receivedPattern)
	}

	srv := &server.Server{}
	if h, ok := receivedHandler.(web.HandlerFunc); ok {
		handler := web.HandlerFunc(srv.Get)
		if &h != &handler {
			t.Fatal("Expected handler not received")
		}
	}

	if !serveCalled {
		t.Fatal("goji.Serve not called")
	}
}
