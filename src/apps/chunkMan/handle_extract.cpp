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

//----------------------------------------------------------------
static bool chunkVisitFunc(const string_q& path, void* data) {
    // LOG_WARN(path);
    if (endsWith(path, "/")) {
        return forEveryFileInFolder(path + "*", chunkVisitFunc, data);

    } else {
        if (!endsWith(path, ".bin"))
            return true;

        blknum_t endBlock = NOPOS;
        blknum_t startBlock = path_2_Bn(path, endBlock);

        COptions* opts = (COptions*)data;
        blknum_t startTest = opts->blocks.start == NOPOS ? 0 : opts->blocks.start;
        blknum_t endTest = opts->blocks.stop == 0 ? NOPOS : opts->blocks.stop;
        if (!inRange(startBlock, startTest, endTest)) {
            LOG_PROG("Skipped: " + path + "\r");
            return true;
        }
        if (!inRange(endBlock, startTest, endTest)) {
            LOG_PROG("Skipped: " + path + "\r");
            return true;
        }

        CIndexArchive index(READING_ARCHIVE);
        if (index.ReadIndexFromBinary(path)) {
            string_q msg = "start: {0} end: {1} fileSize: {2} bloomSize: {3} nAddrs: {4} nRows: {5}";
            replace(msg, "{0}", "{" + padNum9T(startBlock) + "}");
            replace(msg, "{1}", "{" + padNum9T(endBlock) + "}");
            replace(msg, "{2}", "{" + padNum9T(fileSize(path)) + "}");
            replace(
                msg, "{3}",
                "{" + padNum9T(fileSize(substitute(substitute(path, "finalized", "blooms"), ".bin", ".bloom"))) + "}");
            replace(msg, "{4}", "{" + padNum9T(uint64_t(index.header->nAddrs)) + "}");
            replace(msg, "{5}", "{" + padNum9T(uint64_t(index.header->nRows)) + "}");
            if (verbose) {
                string_q m = msg;
                replaceAny(m, "{}", "");
                replaceAll(m, "  ", " ");
                cout << "# " << m << endl;
                cout << "address\tstart\tcount" << endl;
            }
            replaceAll(msg, "{", cGreen);
            replaceAll(msg, "}", cOff);
            cout << msg << endl;

            if (verbose > 0) {
                for (uint32_t a = 0; a < index.nAddrs; a++) {
                    CIndexedAddress* aRec = &index.addresses[a];
                    cout << bytes_2_Addr(aRec->bytes) << "\t" << aRec->offset << "\t" << aRec->cnt << endl;
                    if (verbose > 4) {
                        for (uint32_t b = aRec->offset; b < (aRec->offset + aRec->cnt); b++) {
                            CIndexedAppearance* bRec = &index.appearances[b];
                            cout << "\t" << bRec->blk << "\t" << bRec->txid << endl;
                        }
                    }
                }
            }
        }
    }

    return true;
}

//----------------------------------------------------------------
bool COptions::handle_extract() {
    if (extract == "stats") {
        return handle_stats();
    } else {
        return forEveryFileInFolder(indexFolder_finalized, chunkVisitFunc, this);
    }
    LOG_PROG("Finished");
    return true;
}
