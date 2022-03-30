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
#include "testing.h"

bloom_t addr_2_Bloom2(const address_t& addrIn, CUintArray& litBits);

//------------------------------------------------------------------------
int main(int argc, const char* argv[]) {
    loadEnvironmentPaths("", "");
    etherlib_init(quickQuitHandler);

    cout << "addresses\tparts\tbitsLit" << endl;

    CStringArray addrs{"0x0371a82e4a9d0a4312f3ee2ac9c6958512891372", "0x1296d3bb6dc0efbae431a12939fc15b2823db79b",
                       "0x3d493c51a916f86d6d1c04824b3a7431e61a3ca3", "0xd09022c48298f268c2c431dadb9ca4c2534d9c1c",
                       "0xe1c15164dcfe79431f8421b5a311a829cf0907f3", "0x0341a82e4a9d0a4312f3ee2ac9c6958512891342",
                       "0x1296d3bb6dc0efbae431a12939fc15b2823db49b", "0x3d493c51a916f86d6d1c04824b3a4431e61a3ca3",
                       "0xd09022c48298f268c2c431dadb9ca4c2534d9c1c", "0xe1c15164dcfe49431f8421b5a311a829cf0904f3"};

    CBloomArray blooms, addrBlooms;
    for (size_t i = 0; i < addrs.size(); i++) {
        address_t addr = addrs[i];
        CUintArray bits;
        bloom_t b = addr_2_Bloom2(addr, bits);
        if (!(i % 2)) {
            addToBloomFilter(blooms, addr);
        }
        addrBlooms.push_back(b);
    }
    cout << endl;

    for (size_t i = 0; i < addrs.size(); i++) {
        bool hit = isInBloomFilter(blooms, addrBlooms[i]);
        LOG_INFO(addrs[i], "-", hit, ((i % 2) && hit ? "-fp" : ""));
    }

    etherlib_cleanup();
    return 0;
}

#define BLOOM_WIDTH_IN_BYTES (1048576 / 8)
#define BLOOM_WIDTH_IN_BITS (BLOOM_WIDTH_IN_BYTES * 8)
#define MAX_ADDRS_IN_BLOOM 50000
#define NIBBLE_WID 8
#define K 5
//---------------------------------------------------------------------------
bloom_t addr_2_Bloom2(const address_t& addr, CUintArray& litBits) {
    bloom_t ret;
    cout << addr << " ";
    for (size_t k = 0; k < K; k++) {
        string_q four_byte = extract(addr, 2 + (k * NIBBLE_WID), NIBBLE_WID);
        uint64_t bit64 = str_2_Uint("0x" + four_byte);
        uint64_t bit = (bit64 % BLOOM_WIDTH_IN_BITS);
        ret.lightBit(bit);
        litBits.push_back(bit);
        cout << four_byte << "-" << uint_2_Str(bit64) << "-" << uint_2_Str(bit) << " ";
    }
    cout << endl;
    return ret;
}