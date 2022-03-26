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

uint32_t bytesPerLine = (2048 / 64);
//----------------------------------------------------------------
static bool bloomVisitFunc(const string_q& path, void* data) {
    if (endsWith(path, "/")) {
        return forEveryFileInFolder(path + "*", bloomVisitFunc, data);

    } else {
        blknum_t endBlock = NOPOS;
        blknum_t startBlock = path_2_Bn(path, endBlock);
        blknum_t last = *(blknum_t*)data;
        if (last > startBlock)
            return true;

        CBloomArray blooms;
        readBloomFromBinary(path, blooms);

        ostringstream os;
        cout << "range: {" << startBlock << " " << endBlock << "}" << endl;
        cout << "nBlooms: " << blooms.size() << endl;
        cout << "byteWidth: " << getBloomWidthInBytes() << endl;
        for (auto bloom : blooms) {
            cout << "nInserted: " << bloom.nInserted << endl;
            for (size_t i = 0; i < getBloomWidthInBytes(); i++) {
                if (!(i % bytesPerLine)) {
                    if (i != 0)
                        cout << endl;
                    cout << padNum7T(uint64_t(i)) << ": ";
                }
                uint8_t ch = bloom.bits[i];
                cout << bitset<8>(ch) << ' ';
            }
            cout << endl;
            if (isTestMode()) {
                return false;
            }
        }
    }

    return !shouldQuit();
}

//----------------------------------------------------------------
bool COptions::handle_extract_blooms(void) {
    // bytesPerLine = (2048 / 16);

    blknum_t last = 0;
    if (isTestMode()) {
        bloomVisitFunc(indexFolder_blooms + "000000000-000000000.bloom", &last);
        bloomVisitFunc(indexFolder_blooms + "000000001-000590501.bloom", &last);
        return true;
    }

    return forEveryFileInFolder(indexFolder_blooms, bloomVisitFunc, &last);
}
