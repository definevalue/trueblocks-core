package pinlib

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinlib/manifest"
)

// func makeFileMap(files []fs.FileInfo) map[string]struct{} {
// 	fmap := make(map[string]struct{})

// 	for _, file := range files {
// 		fmap[file.Name()] = struct{}{}
// 	}

// 	return fmap
// }

// func filterNewChunks(files []fs.FileInfo, pins []manifest.PinDescriptor) []fs.FileInfo {
// 	pmap := make(map[string]struct{})
// 	for _, pin := range pins {
// 		pmap[pin.FileName] = struct{}{}
// 	}

// 	newChunks := make([]fs.FileInfo, 0, len(files))
// 	for _, file := range files {
// 		if _, ok := pmap[file.Name()]; ok {
// 			continue
// 		}
// 		newChunks = append(newChunks, file)
// 	}

// 	return newChunks
// }

func makeChecksum(file *os.File) ([]byte, error) {
	checksum := md5.New()
	if _, err := io.Copy(checksum, file); err != nil {
		return nil, err
	}

	return checksum.Sum(nil), nil
}

type ErrMissingChunk struct {
	fileName string
}

func (e *ErrMissingChunk) Error() string {
	return "missing chunk: " + e.fileName
}

type validatedChunk struct {
	cacheType cache.CacheType
	fileName  string
	file      *os.File
	size      int64
	checksum  []byte
	archive   *os.File
	cid       manifest.IpfsHash
}

type ErrMismatchedChecksum struct {
	wrongChecksum    []byte
	expectedChecksum []byte
	fileName         string
}

func (err *ErrMismatchedChecksum) Error() string {
	return fmt.Sprintf("mismatched checksum for file %s: expected %s, got %s", err.fileName, err.expectedChecksum, err.wrongChecksum)
}

func validateExistingChunks(chunks []validatedChunk, pins []manifest.PinDescriptor, existingFiles map[string]struct{}) error {
	pinFileNameToHash := make(map[string]string)
	for _, pin := range pins {
		// We require all chunks listed in manifest to exist on disk
		if _, ok := existingFiles[pin.FileName]; !ok {
			return &ErrMissingChunk{
				fileName: pin.FileName,
			}
		}

		pinHash := pin.IndexChecksum
		// We expect all chunks to be of the same cacheType
		if chunks[0].cacheType == cache.BloomChunk {
			pinHash = pin.IndexChecksum
		}

		pinFileNameToHash[pin.FileName] = pinHash
	}

	for _, chunk := range chunks {
		pinHash, ok := pinFileNameToHash[chunk.fileName]
		// The chunk is new, so its hash is missing in manifest
		if !ok {
			continue
		}

		if pinHash != string(chunk.checksum) {
			return &ErrMismatchedChecksum{
				fileName:         chunk.fileName,
				expectedChecksum: []byte(pinHash),
				wrongChecksum:    chunk.checksum,
			}
		}
	}

	return nil
}

func prepareChunksByType(cacheType cache.CacheType, existingManifest *manifest.Manifest) ([]validatedChunk, error) {
	cachePath := cache.Path{}
	cachePath.New(cacheType)

	existingFiles := make(map[string]struct{})
	files, err := ioutil.ReadDir(cachePath.String())
	if err != nil {
		return nil, err
	}

	validatedChunks := make([]validatedChunk, 0, len(files))
	for _, file := range files {
		openedFile, err := os.Open(cachePath.GetFullPath(file.Name()))
		if err != nil {
			return nil, err
		}
		existingFiles[file.Name()] = struct{}{}

		checksum, err := makeChecksum(openedFile)
		if err != nil {
			return nil, err
		}

		// TODO: replace with archive size
		// stat, err := openedFile.Stat()
		// if err != nil {
		// 	return nil, err
		// }
		validatedChunks = append(validatedChunks, validatedChunk{
			fileName: openedFile.Name(),
			file:     openedFile,
			// size:     stat.Size(),
			checksum: checksum,
		})
	}

	err = validateExistingChunks(validatedChunks, existingManifest.Pins, existingFiles)

	// TODO
	return nil, err
}

// func PublishNewManifest(existingManifest *manifest.Manifest) (*manifest.Manifest, error) {
// 	// 0. check: all chunks listed in manifest are on the disk
// 	// 1. Read chunk list from disk and leave ones not listed in manifest
// 	// 2. compute their md5
// 	// 3. zip them and save zip archive file size
// 	// 4. pin them
// 	// 5. write new manifest
// 	// 6. pin manifest
// 	// bloomsPath := &cache.Path{}
// 	// bloomsPath.New(cache.BloomChunk)
// 	// indexPath := &cache.Path{}
// 	// indexPath.New(cache.IndexChunk)
// 	indexPath := config.ReadTrueBlocks().Settings.IndexPath

// 	files, err := ioutil.ReadDir(indexPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	missingChunk, missingChunkName := areAnyChunksMissing(files, existingManifest.Pins)
// 	if missingChunk {
// 		return nil, &ErrMissingChunk{
// 			fileName: missingChunkName,
// 		}
// 	}

// 	// Not really needed, just open all files, check checksum, then filer out the already present ones
// 	// or do not filter out and try to pin all, then validate CIDs (redundant?)
// 	// newChunks := filterNewChunks(files, existingManifest.Pins)
// 	// if len(newChunks) == 0 {
// 	// 	return nil, errors.New("no new chunks")
// 	// }

// 	// for _, newChunk := range newChunks {
// 	// 	os.Open()
// 	// }

// 	checksums := makeChecksums(newChunks)
// }
