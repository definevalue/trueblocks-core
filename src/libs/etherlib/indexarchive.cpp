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
#include "indexarchive.h"

namespace qblocks {
//----------------------------------------------------------------
CIndexArchive::CIndexArchive(bool mode) : CArchive(mode) {
    addresses = nullptr;
    appearances = nullptr;
    rawData = nullptr;
    reverseAppMap = nullptr;
}

//----------------------------------------------------------------
CIndexArchive::~CIndexArchive(void) {
    clean();
}

//----------------------------------------------------------------
void CIndexArchive::clean(void) {
    if (rawData) {
        delete[] rawData;
        rawData = nullptr;
        header = CIndexHeader();
        addresses = nullptr;
        appearances = nullptr;
    }
    if (reverseAppMap) {
        delete[] reverseAppMap;
        reverseAppMap = nullptr;
    }
    reverseAddrRanges.clear();
    Release();
}

//----------------------------------------------------------------
bool CIndexArchive::ReadIndexFromBinary(const string_q& path, indexparts_t parts) {
    if (!m_isReading)
        return false;
    if (!fileExists(path))
        return false;
    if (!Lock(path, modeReadOnly, LOCK_NOWAIT))
        return false;
    if (parts == IP_NONE)
        return false;

    size_t readSize = sizeof(CIndexHeader);
    size_t memorySize = readSize;
    if (parts != IP_HEADER) {
        readSize = fileSize(path);
        memorySize = readSize + (2 * asciiAppearanceSize);  // a bit of extra room for some reason...
    }

    rawData = reinterpret_cast<char*>(malloc(memorySize));
    if (!rawData) {
        LOG_ERR("Could not allocate memory for data.");
        Release();
        return false;
    }
    bzero(rawData, memorySize);

    size_t nRead = Read(rawData, readSize, sizeof(char));
    if (nRead != readSize) {
        LOG_ERR("Could not read entire file.");
        clean();
        return false;
    }

    header = *(reinterpret_cast<CIndexHeader*>(rawData));
    ASSERT(header.magic == MAGIC_NUMBER);
    ASSERT(bytes_2_Hash(header.hash) == versionHash);

    size_t startOfAddrTable = sizeof(CIndexHeader);
    size_t addrTableSize = sizeof(CIndexedAddress) * header.nAddrs;
    size_t startOfAppsTable = startOfAddrTable + addrTableSize;

    if (parts == IP_HEADER) {
        addresses = nullptr;
        appearances = nullptr;
    } else {
        addresses = (CIndexedAddress*)(rawData + startOfAddrTable);       // NOLINT
        appearances = (CIndexedAppearance*)(rawData + startOfAppsTable);  // NOLINT
    }
    Release();
    return true;
}

//--------------------------------------------------------------
bool readIndexHeader(const string_q& path, CIndexHeader& header) {
    header.nApps = header.nAddrs = (uint32_t)-1;
    if (contains(path, "blooms")) {
        return false;
    }

    if (endsWith(path, ".txt")) {
        header.nApps = (uint32_t)fileSize(path) / (uint32_t)asciiAppearanceSize;
        CStringArray lines;
        asciiFileToLines(path, lines);
        CAddressBoolMap addrMap;
        for (auto line : lines)
            addrMap[nextTokenClear(line, '\t')] = true;
        header.nAddrs = (uint32_t)addrMap.size();
        return true;
    }

    CArchive archive(READING_ARCHIVE);
    if (!archive.Lock(path, modeReadOnly, LOCK_NOWAIT))
        return false;

    bzero(&header, sizeof(header));
    // size_t nRead =
    archive.Read(&header, sizeof(header), 1);
    // if (false) { //nRead != sizeof(header)) {
    //    cerr << "Could not read file: " << path << endl;
    //    return;
    //}
    ASSERT(header.magic == MAGIC_NUMBER);
    // ASSERT(bytes_2_Hash(h->hash) == versionHash);
    archive.Release();
    return true;
}

//-----------------------------------------------------------------------
int sortRecords(const void* i1, const void* i2) {
    int32_t* p1 = (int32_t*)i1;
    int32_t* p2 = (int32_t*)i2;
    if (p1[1] == p2[1]) {
        if (p1[2] == p2[2]) {
            return (p1[0] - p2[0]);
        }
        return p1[2] - p2[2];
    }
    return p1[1] - p2[1];
}

//-----------------------------------------------------------------------
bool CIndexArchive::LoadReverseMaps(const blkrange_t& range) {
    if (reverseAppMap) {
        delete[] reverseAppMap;
        reverseAddrRanges.clear();
        reverseAppMap = nullptr;
    }

    uint32_t nAppsHere = header.nApps;

    string_q mapFile = substitute(getFilename(), indexFolder_finalized, indexFolder_map);
    if (fileExists(mapFile)) {
        CArchive archive(READING_ARCHIVE);
        if (!archive.Lock(mapFile, modeReadOnly, LOCK_NOWAIT)) {
            LOG_ERR("Could not open file ", mapFile);
            return false;
        }
        size_t nRecords = fileSize(mapFile) / sizeof(CReverseAppMapEntry);
        ASSERT(nRecords == nAppsHere);
        // Cleaned up on destruction of the chunk
        reverseAppMap = new CReverseAppMapEntry[nRecords];
        if (!reverseAppMap) {
            LOG_ERR("Could not allocate memory for CReverseAppMapEntry");
            return false;
        }
        archive.Read((char*)reverseAppMap, sizeof(char), nRecords * sizeof(CReverseAppMapEntry));
        archive.Release();
        blknum_t cur = 0;
        for (uint32_t i = 0; i < header.nAddrs; i++) {
            blkrange_t r;
            r.first = cur + addresses[i].offset;
            r.second = r.first + addresses[i].cnt - 1;
            reverseAddrRanges.push_back(r);
        }
        return true;
    }

    // Cleaned up on destruction of the chunk
    reverseAppMap = new CReverseAppMapEntry[nAppsHere];
    if (!reverseAppMap) {
        LOG_ERR("Could not allocate memory for CReverseAppMapEntry");
        return false;
    }
    for (uint32_t i = 0; i < nAppsHere; i++) {
        reverseAppMap[i].n = i;
        reverseAppMap[i].blk = appearances[i].blk;
        reverseAppMap[i].tx = appearances[i].txid;
    }

    blknum_t cur = 0;
    for (uint32_t i = 0; i < header.nAddrs; i++) {
        blkrange_t r;
        r.first = cur + addresses[i].offset;
        r.second = r.first + addresses[i].cnt - 1;
        reverseAddrRanges.push_back(r);
    }

    qsort(reverseAppMap, nAppsHere, sizeof(CReverseAppMapEntry), sortRecords);

    CArchive archive(WRITING_ARCHIVE);
    if (!archive.Lock(mapFile, modeWriteCreate, LOCK_WAIT)) {
        LOG_ERR("Could not open file ", mapFile);
        return false;
    }
    archive.Write(reverseAppMap, sizeof(char), nAppsHere * sizeof(CReverseAppMapEntry));
    archive.Release();

    LOG_PROG("Processed: " + getFilename());
    return true;
}

}  // namespace qblocks
