// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/golang/dep"
	"github.com/golang/dep/internal/gps"
	"github.com/golang/dep/internal/gps/pkgtree"
)

// gopathAnalyzer deduces configuration from the projects in the GOPATH
type gopathAnalyzer struct {
	loggers *Loggers
	ctx     *dep.Ctx
	pkgT    pkgtree.PackageTree
	cpr     string
	sm      *gps.SourceMgr

	pd projectData
}

func newGopathAnalyzer(loggers *Loggers, ctx *dep.Ctx, pkgT pkgtree.PackageTree, cpr string, sm *gps.SourceMgr) *gopathAnalyzer {
	return &gopathAnalyzer{
		loggers: loggers,
		ctx:     ctx,
		pkgT:    pkgT,
		cpr:     cpr,
		sm:      sm,
	}
}

func (a *gopathAnalyzer) DeriveRootManifestAndLock(path string, n gps.ProjectRoot) (*dep.Manifest, *dep.Lock, error) {
	var err error

	a.pd, err = getProjectData(a.ctx, a.loggers, a.pkgT, a.cpr, a.sm)
	if err != nil {
		return nil, nil, err
	}
	m := &dep.Manifest{
		Dependencies: a.pd.constraints,
	}

	// Make an initial lock from what knowledge we've collected about the
	// versions on disk
	l := &dep.Lock{
		P: make([]gps.LockedProject, 0, len(a.pd.ondisk)),
	}

	for pr, v := range a.pd.ondisk {
		// That we have to chop off these path prefixes is a symptom of
		// a problem in gps itself
		pkgs := make([]string, 0, len(a.pd.dependencies[pr]))
		prslash := string(pr) + "/"
		for _, pkg := range a.pd.dependencies[pr] {
			if pkg == string(pr) {
				pkgs = append(pkgs, ".")
			} else {
				pkgs = append(pkgs, trimPathPrefix(pkg, prslash))
			}
		}

		l.P = append(l.P, gps.NewLockedProject(
			gps.ProjectIdentifier{ProjectRoot: pr}, v, pkgs),
		)
	}

	return m, l, nil
}

func (a *gopathAnalyzer) PostSolveShenanigans(m *dep.Manifest, l *dep.Lock) {
	// Pick notondisk project constraints from solution and add to manifest
	for k, _ := range a.pd.notondisk {
		for _, x := range l.Projects() {
			if k == x.Ident().ProjectRoot {
				m.Dependencies[k] = getProjectPropertiesFromVersion(x.Version())
				break
			}
		}
	}
}
