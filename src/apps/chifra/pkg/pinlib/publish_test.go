package pinlib

// import (
// 	"bytes"
// 	"crypto/md5"
// 	"io"
// 	"os"
// 	"testing"

// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinlib/manifest"
// )

// func Test_validateExistingChunks_Checksum(t *testing.T) {
// 	h := md5.New()
// 	io.WriteString(h, "some string")
// 	checksum := h.Sum(nil)

// 	chunks := []chunkMetadata{
// 		{
// 			fileName:  "test-file",
// 			file:      &os.File{},
// 			cacheType: cache.BloomChunk,
// 			checksum:  checksum,
// 		},
// 	}
// 	pins := []manifest.PinDescriptor{
// 		{
// 			FileName:      "test-file",
// 			IndexChecksum: string(checksum),
// 		},
// 	}
// 	files := map[string]struct{}{
// 		"test-file": {},
// 	}

// 	err := validateExistingChunks(chunks, pins, files)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func Test_validateExistingChunks_ChecksumInvalid(t *testing.T) {
// 	h := md5.New()
// 	io.WriteString(h, "some string")
// 	checksum := h.Sum(nil)

// 	chunks := []chunkMetadata{
// 		{
// 			fileName:  "test-file",
// 			file:      &os.File{},
// 			cacheType: cache.BloomChunk,
// 			checksum:  checksum,
// 		},
// 	}
// 	pins := []manifest.PinDescriptor{
// 		{
// 			FileName:      "test-file",
// 			IndexChecksum: "5e5037601c72f18ab11ee51ca59c9fae",
// 		},
// 	}
// 	files := map[string]struct{}{
// 		"test-file": {},
// 	}

// 	err := validateExistingChunks(chunks, pins, files)
// 	if err == nil {
// 		t.Error("expected error")
// 	}

// 	errMismatch, ok := err.(*ErrMismatchedChecksum)
// 	if !ok {
// 		t.Error("expected ErrMismatchedChecksum")
// 	}
// 	if errMismatch.fileName != chunks[0].fileName ||
// 		!bytes.Equal(errMismatch.expectedChecksum, []byte(pins[0].IndexChecksum)) ||
// 		!bytes.Equal(errMismatch.wrongChecksum, chunks[0].checksum) {
// 		t.Errorf(
// 			"wrong error details: %s, %s, %s",
// 			errMismatch.fileName,
// 			errMismatch.expectedChecksum,
// 			errMismatch.wrongChecksum,
// 		)
// 	}
// }

// func Test_validateExistingChunks_FileMismatch(t *testing.T) {
// 	h := md5.New()
// 	io.WriteString(h, "some string")
// 	checksum := h.Sum(nil)

// 	chunks := []chunkMetadata{
// 		{
// 			fileName:  "test-file",
// 			file:      &os.File{},
// 			cacheType: cache.BloomChunk,
// 			checksum:  checksum,
// 		},
// 	}
// 	pins := []manifest.PinDescriptor{
// 		{
// 			FileName:      "test-file",
// 			IndexChecksum: string(checksum),
// 		},
// 	}
// 	files := map[string]struct{}{
// 		"different-test-file": {},
// 	}

// 	err := validateExistingChunks(chunks, pins, files)
// 	if err == nil {
// 		t.Error("expected error")
// 	}
// }
