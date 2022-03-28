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

//-------------------------------------------------------------------------------------------------
static char delim = '|';

//-------------------------------------------------------------------------------------------------
string_q STR_INDEX_STATS =
    "+start|end|nAddrs|nApps|nBlocks|nBlooms|addrsPerBlock|appsPerBlock|appsPerAddr|recWid|bloomSz|chunkSz|ratio+";

//-------------------------------------------------------------------------------------------------
string_q getHeaders(void) {
    string_q fieldList = STR_INDEX_STATS;
    replaceAll(fieldList, "|", "\",\"");
    replaceAll(fieldList, "+", "\"");
    return fieldList;
}

//-------------------------------------------------------------------------------------------------
class CThing {
  public:
    string_q indexPath;
    string_q bloomPath;
    blkrange_t range;
    CBloomArray blooms;
    CIndexHeader header;
    CThing(void) {
    }
    CThing(const string_q& fn) {
        bloomPath = fn;
        indexPath = substitute(substitute(bloomPath, "blooms", "finalized"), ".bloom", ".bin");
        range.second = NOPOS;
        range.first = path_2_Bn(bloomPath, range.second);
    }
    uint64_t nBlocks(void) const {
        return (range.second - range.first) + 1;
    }
    uint64_t nBlooms(void) const {
        return blooms.size();
    }
    string_q addressesPerBlock(void) const {
        return double_2_Str((nBlocks() ? float(header.nAddrs) / float(nBlocks()) : 0.), 3);
    }
    string_q appearancesPerBlock(void) const {
        return double_2_Str(nBlocks() ? float(header.nRows) / float(nBlocks()) : 0., 3);
    }
    string_q appearancesPerAddress(void) const {
        return (double_2_Str(header.nAddrs ? float(header.nRows) / float(header.nAddrs) : 0., 3));
    }
    size_t bloomSize(void) const {
        return (fileSize(bloomPath));
    }
    size_t indexSize(void) const {
        return (fileSize(indexPath));
    }
    string_q bloomRatio(void) const {
        return (double_2_Str(bloomSize() ? float(indexSize()) / float(bloomSize()) : 0., 3));
    }
    size_t recordSize(void) const {
        return (sizeof(uint32_t) + getBloomWidthInBytes());
    }
    bool filter(blknum_t last, bool& filtered) const {
        if (isTestMode()) {
            if (range.first > 2000000 && range.first < 3000000) {
                // too slow, so skip for testing
                filtered = true;
                return true;

            } else if (range.first > 4000000) {
                // enough already
                filtered = false;
                return true;
            }
        } else {
            cerr << range << "\r";
            cerr.flush();
        }

        if (range.first == 0) {
            if (last != 0) {
                filtered = true;
                return true;
            }
        } else {
            if (last >= range.first) {
                filtered = true;
                return true;
            }
        }

        return false;
    }
};

//--------------------------------------------------------------
bool bloomVisitFunc(const string_q& pp, void* data) {
    if (endsWith(pp, "/")) {
        return forEveryFileInFolder(pp + "*", bloomVisitFunc, data);

    } else {
        if (!endsWith(pp, ".bloom"))
            return true;

        bool ret;
        CThing thing(pp);
        if (thing.filter(((COptions*)data)->last, ret))
            return ret;
        readBloomFromBinary(thing.bloomPath, thing.blooms, false /* readBits */);
        readIndexHeader(thing.indexPath, thing.header);

        ostringstream os;
        os << thing.range.first << delim;
        os << thing.range.second << delim;
        os << thing.header.nAddrs << delim;
        os << thing.header.nRows << delim;
        os << thing.nBlocks() << delim;
        os << thing.nBlooms() << delim;
        os << thing.addressesPerBlock() << delim;
        os << thing.appearancesPerBlock() << delim;
        os << thing.appearancesPerAddress() << delim;
        os << thing.recordSize() << delim;
        os << thing.bloomSize() << delim;
        os << thing.indexSize() << delim;
        os << thing.bloomRatio();
        string_q str = os.str();
        replaceAll(str, "|", "\",\"");
        str = "\"" + str + "\"";
        cout << str << endl;
        appendToAsciiFile(cacheFolder_tmp + "chunk_stats.csv", str + "\n");
    }

    return true;
}

//----------------------------------------------------------------
bool COptions::handle_stats() {
    cout << getHeaders() << endl;
    return forEveryFileInFolder(indexFolder_blooms, bloomVisitFunc, this);
}
