package pinlib

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/ipfs"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinlib/manifest"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/progress"
	ants "github.com/panjf2000/ants/v2"
)

type ErrMissingChunk struct {
	fileName string
}

func (e *ErrMissingChunk) Error() string {
	return "missing chunk: " + e.fileName
}

type chunkMetadata struct {
	cacheType cache.CacheType
	fileName  string
	file      *os.File
	size      int64
	checksum  []byte
	archive   *os.File
	cid       manifest.IpfsHash
}

func validateChunksNotMissing(pins []manifest.PinDescriptor, existingFiles map[string]struct{}) error {
	for _, pin := range pins {
		// We require all chunks listed in manifest to exist on disk
		if _, ok := existingFiles[pin.FileName]; !ok {
			return &ErrMissingChunk{
				fileName: pin.FileName,
			}
		}
	}

	return nil
}

type ErrMismatchedCid struct {
	fileName    string
	expectedCid string
	actualCid   string
}

func (err *ErrMismatchedCid) Error() string {
	return fmt.Sprintf(
		"wrong CID for file %s: expected %s, but got %s",
		err.fileName,
		err.expectedCid,
		err.actualCid,
	)
}

type pinDescriptorFragment struct {
	ChunkType cache.CacheType
	FileName  string
	FileSize  int64
	Cid       manifest.IpfsHash
}

func prepareChunksByType(chain string, cacheType cache.CacheType, existingManifest *manifest.Manifest) ([]pinDescriptorFragment, error) {
	cachePath := cache.Path{}
	cachePath.New(chain, cacheType)
	ctx, cancel := context.WithCancel(context.Background())
	poolSize := 4

	readGroup := sync.WaitGroup{}
	writeGroup := sync.WaitGroup{}

	// existingFiles := make(map[string]struct{})
	files, err := ioutil.ReadDir(cachePath.String())
	if err != nil {
		cancel()
		return nil, err
	}
	pinToCid := make(map[string]string, len(existingManifest.Pins))
	for _, pin := range existingManifest.Pins {
		cid := pin.BloomHash
		if cacheType == cache.IndexChunk {
			cid = pin.IndexHash
		}

		pinToCid[pin.FileName] = cid
	}

	// metadata := make([]chunkMetadata, 0, len(files))

	progressChannel := make(chan progress.Progress)
	fileChannel := make(chan *os.File)
	fragmentChannel := make(chan *pinDescriptorFragment)

	readPool, err := ants.NewPoolWithFunc(poolSize, func(param interface{}) {
		defer readGroup.Done()
		select {
		case <-ctx.Done():
			return
		default:
			file := param.(fs.FileInfo)
			openedFile, err := os.Open(cachePath.GetFullPath(file.Name()))
			if err != nil {
				progressChannel <- progress.Progress{
					Event:   progress.Error,
					Message: err.Error(),
				}
				cancel()
				return
			}
			fileChannel <- openedFile
		}
	})
	if err != nil {
		cancel()
		return nil, err
	}
	defer readPool.Release()

	writePool, err := ants.NewPoolWithFunc(poolSize, func(param interface{}) {
		file := param.(*os.File)
		// - gzip
		archive := bytes.Buffer{}
		archiveWriter := gzip.NewWriter(&archive)

		// We are not setting archive timestamp to make sure we always get the same
		// archive for a given content
		archiveWriter.Name = file.Name()
		written, err := io.Copy(archiveWriter, file)
		if err != nil {
			progressChannel <- progress.Progress{
				Event:   progress.Error,
				Message: err.Error(),
			}
			cancel()
			return
		}

		// - keep file size

		// - add to ipfs

		ipfsShell := ipfs.Connect(config.GetRootConfig().Settings.IpfsNode)
		cid, err := ipfsShell.Add(bytes.NewReader(archive.Bytes()))
		if err != nil {
			progressChannel <- progress.Progress{
				Event:   progress.Error,
				Message: err.Error(),
			}
			cancel()
			return
		}

		// - check if CID matches manifest
		pinFileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		previousCid, ok := pinToCid[pinFileName]
		if ok && previousCid != cid {
			err = &ErrMismatchedCid{
				fileName:    pinFileName,
				expectedCid: previousCid,
				actualCid:   cid,
			}
			progressChannel <- progress.Progress{
				Event:   progress.Error,
				Message: err.Error(),
			}
			cancel()
			return
		}

		fragmentChannel <- &pinDescriptorFragment{
			FileName: pinFileName,
			FileSize: written,
			Cid:      cid,
		}
	})
	if err != nil {
		cancel()
		return nil, err
	}
	defer writePool.Release()

	writeGroup.Add(1)
	go func() {
		for file := range fileChannel {
			writeGroup.Add(1)
			writePool.Invoke(file)
		}
		writeGroup.Done()
	}()

	for _, file := range files {
		readGroup.Add(1)
		readPool.Invoke(file)
	}

	var result []pinDescriptorFragment
	for fragment := range fragmentChannel {
		result = append(result, *fragment)
	}

	// err = validateChunksNotMissing(existingManifest.Pins, existingFiles)

	// TODO
	return result, err
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
