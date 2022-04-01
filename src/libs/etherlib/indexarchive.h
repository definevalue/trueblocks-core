#pragma once
/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
#include "etherlib.h"
#include "appearance.h"
#include "bloom.h"
#include "indexheader.h"
#include "indexedaddress.h"
#include "indexedappearance.h"

namespace qblocks {

typedef struct CReverseAppMapEntry {
  public:
    uint32_t n;
    uint32_t blk;
    uint32_t tx;
} CReverseAppMapEntry;

typedef enum {
    IP_NONE = 0,
    IP_HEADER = (1 << 1),
    IP_ADDRS = (1 << 2),
    IP_APPS = (1 << 3),
    IP_ALL = (IP_HEADER | IP_ADDRS | IP_APPS),
} indexparts_t;

//---------------------------------------------------------------------------
class CIndexArchive : public CArchive {
  public:
    CIndexHeader header;
    CIndexedAddress* addresses;
    CIndexedAppearance* appearances;
    CBlockRangeArray reverseAddrRanges;
    CReverseAppMapEntry* reverseAppMap{nullptr};

    explicit CIndexArchive(bool mode);
    ~CIndexArchive(void);
    bool ReadIndexFromBinary(const string_q& fn, indexparts_t parts);
    bool LoadReverseMaps(const blkrange_t& range);

  private:
    char* rawData;
    CIndexArchive(void) : CArchive(READING_ARCHIVE) {
    }
    void clean(void);
};

//-----------------------------------------------------------------------
#define MAGIC_NUMBER ((uint32_t)str_2_Uint("0xdeadbeef"))
extern hash_t versionHash;
//--------------------------------------------------------------
typedef bool (*INDEXCHUNKFUNC)(CIndexArchive& chunk, void* data);
typedef bool (*INDEXBLOOMFUNC)(CBloomFilter& blooms, void* data);
typedef bool (*ADDRESSFUNC)(const address_t& addr, void* data);
class CChunkVisitor {
  public:
    INDEXCHUNKFUNC indexFunc = nullptr;
    ADDRESSFUNC addrFunc = nullptr;
    void* callData = nullptr;
    blkrange_t range = make_pair(0, NOPOS);
};
extern bool readIndexHeader(const string_q& inFn, CIndexHeader& header);
}  // namespace qblocks

#if 0
// extern bool forEveryIndexChunk(INDEXCHUNKFUNC func, void* data);
// extern bool forEveryAddressInIndex(ADDRESSFUNC func, const blkrange_t& range, void* data);
// extern bool forEverySmartContractInIndex(ADDRESSFUNC func, void* data);
#endif
