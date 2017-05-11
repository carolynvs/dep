// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dep

import (
	"os"
	"path/filepath"

	"github.com/golang/dep/internal/gps"
)

type Analyzer struct{}

func (a Analyzer) HasConfig(path string) bool {
	mf := filepath.Join(path, ManifestName)
	fileOK, err := IsRegular(mf)
	return err == nil && fileOK
}

func (a Analyzer) DeriveManifestAndLock(path string, n gps.ProjectRoot) (gps.Manifest, gps.Lock, error) {
	if !a.HasConfig(path) {
		return nil, nil, nil
	}

	f, err := os.Open(filepath.Join(path, ManifestName))
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	m, err := readManifest(f)
	if err != nil {
		return nil, nil, err
	}
	// TODO: No need to return lock til we decide about preferred versions, see
	// https://github.com/sdboyer/gps/wiki/gps-for-Implementors#preferred-versions.
	return m, nil, nil
}

func (a Analyzer) Info() (string, int) {
	return "dep", 1
}
