package pinlib

import (
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinlib/manifest"
)

// ReadLocalManifest loads the manifest saved in ConfigPath
func ReadLocalManifest() (*manifest.Manifest, error) {
	file, err := os.Open(config.GetPathToConfig(false /* withChain */) + "manifest/manifest.json")
	if err != nil {
		return nil, err
	}

	return manifest.ReadJSONManifest(file)
}
