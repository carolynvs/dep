// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"

	"github.com/golang/dep"
	"github.com/golang/dep/internal/gps"
)

// compositeAnalyzer overlays configuration from multiple analyzers
type compositeAnalyzer struct {
	// Analyzers is the set of analyzers to apply, last one wins any conflicts
	Analyzers []rootProjectAnalyzer
}

func (a compositeAnalyzer) DeriveRootManifestAndLock(path string, n gps.ProjectRoot) (*dep.Manifest, *dep.Lock, error) {
	return nil, nil, errors.New("Not implemented")
}

func (a compositeAnalyzer) PostSolveShenanigans(*dep.Manifest, *dep.Lock) {
	panic("not implemented")
}
