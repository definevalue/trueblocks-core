// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package manifest

import (
	"encoding/json"
	"fmt"
	"io"
)

var supportedManifestVersion = uint(2)

// ReadJSONManifest reads manifest in JSON format
func ReadJSONManifest(reader io.Reader) (*Manifest, error) {
	decoder := json.NewDecoder(reader)
	manifest := &Manifest{}

	err := decoder.Decode(manifest)
	if err == nil {
		version, err := manifest.Version.Int64()
		if err != nil || version != int64(supportedManifestVersion) {
			return nil, fmt.Errorf("manifest version %d is not supported", version)
		}
	}

	return manifest, err
}
