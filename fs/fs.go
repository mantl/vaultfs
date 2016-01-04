// Copyright © 2016 Asteris, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/Sirupsen/logrus"
)

// VaultFS is a vault filesystem
type VaultFS struct{}

// New returns a new VaultFS
func New() *VaultFS {
	logrus.Debug("created new FS")
	return &VaultFS{}
}

// Mount the FS at the given mountpoint
func (v *VaultFS) Mount(mountpoint string) (stop func(), errs chan error) {
	conn, err := fuse.Mount(
		mountpoint,
		fuse.FSName("vault"),
		fuse.VolumeName("vault"),
	)

	stop = func() {
		logrus.Info("closing FUSE connection")
		conn.Close()

		logrus.Debug("closed connection, waiting for ready")
		<-conn.Ready
		if conn.MountError != nil {
			errs <- err
		}
		close(errs)
	}
	errs = make(chan error, 1)

	logrus.Debug("created conn")
	if err != nil {
		errs <- err
		close(errs)
		return stop, errs
	}

	logrus.Debug("starting to serve")
	go func() {
		err := fs.Serve(conn, v)
		if err != nil {
			errs <- err
		}
	}()

	return stop, errs
}

// Root returns the struct that does the actual work
func (VaultFS) Root() (fs.Node, error) {
	logrus.Debug("returning root")
	return Root{}, nil
}