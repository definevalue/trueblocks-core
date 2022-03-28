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
#include "options.h"

string_q STR_INDEX_STATS =
    ("start|end|nBlks|nAddrs|nApps|nBits|nA/nB|nP/nB|nP/nA|nA/Bt|nBlooms|recWid|bloomSz|chunkSz|comp");
static char delim = ',';

//--------------------------------------------------------------
class CThing {
  public:
    blkrange_t range;
    CBloomArray blooms;
    CIndexHeader header;
    string_q getHeaders(void) const {
        string_q fieldList = STR_INDEX_STATS;
        CStringArray fields;
        explode(fields, fieldList, '|');
        ostringstream os;
        for (auto field : fields) {
            os << field << delim;
        }
        return os.str();
    }
    uint64_t nBlocks(void) const {
        return (range.second - range.first) + 1;
    }
    uint64_t nBlooms(void) const {
        return blooms.size();
    }
    uint64_t nBloomInserts(void) const {
        size_t cnt = 0;
        for (auto bloom : blooms) {
            cnt += bloom.nInserted;
        }
        return cnt;
    }
    uint64_t totalBits(void) const {
        size_t ret = 0;
        for (auto bloom : blooms) {
            ret += 1;
        }
        return ret;
    }
};

//--------------------------------------------------------------
static bool bloomVisitFunc(const string_q& path, void* data) {
    if (endsWith(path, "/")) {
        return forEveryFileInFolder(path + "*", bloomVisitFunc, data);

    } else {
        if (!endsWith(path, ".bloom"))
            return true;

        COptions* opts = (COptions*)data;

        CThing thing;
        thing.range.second = NOPOS;
        thing.range.first = path_2_Bn(path, thing.range.second);
        // cout << path << "\t" << thing.range << "\t" << opts->last << endl;

        if (thing.range.first == 0) {
            if (opts->last != 0)
                return true;
        } else {
            if (opts->last >= thing.range.first)
                return true;
        }

        if (isTestMode()) {
            if (thing.range.first > 2000000 && thing.range.first < 3000000) {
                // too slow, so skip for testing
                return true;
            }

            if (thing.range.first > 4000000) {
                // enough already
                return false;
            }
        }

        string_q chunkPath = substitute(substitute(path, "blooms", "finalized"), ".bloom", ".bin");

        readBloomFromBinary(path, thing.blooms, false /* readBits */);
        readIndexHeader(chunkPath, thing.header);

        size_t recordSize = (sizeof(uint32_t) + getBloomWidthInBytes());

        ostringstream os;
        os << thing.range.first << delim;
        os << thing.range.second << delim;
        os << thing.nBlocks() << delim;
        os << thing.header.nAddrs << delim;
        os << thing.header.nRows << delim;
        os << thing.totalBits() << delim;
        os << double_2_Str(thing.nBlocks() ? float(thing.header.nAddrs) / float(thing.nBlocks()) : 0., 3) << delim;
        os << double_2_Str(thing.nBlocks() ? float(thing.header.nRows) / float(thing.nBlocks()) : 0., 3) << delim;
        os << double_2_Str(thing.header.nAddrs ? float(thing.header.nRows) / float(thing.header.nAddrs) : 0., 3)
           << delim;
        os << double_2_Str(thing.totalBits() ? float(thing.nBloomInserts()) / float(thing.totalBits()) : 0., 3)
           << delim;
        os << thing.nBlooms() << delim;
        os << recordSize << delim;
        os << fileSize(path) << delim;
        os << fileSize(chunkPath) << delim;
        os << double_2_Str(fileSize(path) ? float(fileSize(chunkPath)) / float(fileSize(path)) : 0., 3);
        os << endl;
        appendToAsciiFile(cacheFolder_tmp + "chunk_stats.csv", os.str());
        cout << os.str();
    }

    return true;
}

//----------------------------------------------------------------
bool COptions::handle_stats() {
    CThing thing;
    cout << thing.getHeaders() << endl;

    readCache();
    return forEveryFileInFolder(indexFolder_blooms, bloomVisitFunc, this);
}

//----------------------------------------------------------------
bool COptions::readCache() {
    CStringArray lines;
    asciiFileToLines(cacheFolder_tmp + "chunk_stats.csv", lines);
    for (auto line : lines) {
        blknum_t start = str_2_Uint(line);
        if (!isTestMode() || start <= 4000000) {
            cout << line << endl;
        }
    }

    last = 0;
    if (lines.size() > 0)
        last = str_2_Uint(lines[lines.size() - 1]);

    return true;
}

// string_q checkSize = sizeof(uint32_t) + (thing.nBlooms() * recordSize) == fileSize(path) ? greenCheck : redX;
// os << checkSize << delim;
// os << ((thing.header.nAddrs == thing.nBloomInserts()) ? greenCheck : redX);
