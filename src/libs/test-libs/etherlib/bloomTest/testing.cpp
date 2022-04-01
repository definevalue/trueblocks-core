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

bool addAddrToBloom2(CBloomFilter& blooms, const address_t& addr);
bloom_t addr_2_Bloom2(const address_t& addr, CStringArray& parts, CUintArray& values, CUintArray& litBits);
#define BLOOM_WIDTH_IN_BYTES (1048576 / 8)
#define BLOOM_WIDTH_IN_BITS (BLOOM_WIDTH_IN_BYTES * 8)
#define MAX_ADDRS_IN_BLOOM 50000
#define NIBBLE_WID 8
#define K 5

class CBloomStats {
  public:
    uint64_t nBlooms;
    uint64_t nInserted;
    uint64_t nBitsLit;
    uint64_t nBitsNotLit;
    uint64_t sz;
    CUintArray bitsLit;
    CBloomStats(const CBloomFilter& bloom);
    void Report(void) const;
};

//------------------------------------------------------------------------
int main(int argc, const char* argv[]) {
    loadEnvironmentPaths("", "");
    etherlib_init(quickQuitHandler);

    CStringArray addrsInSet{"0x0371a82e4a9d0a4312f3ee2ac9c6958512891372", "0x3d493c51a916f86d6d1c04824b3a7431e61a3ca3",
                            "0xe1c15164dcfe79431f8421b5a311a829cf0907f3", "0x1296d3bb6dc0efbae431a12939fc15b2823db49b",
                            "0xd09022c48298f268c2c431dadb9ca4c2534d9c1c"};

    CStringArray addrsNotInSet{
        "0x1296d3bb6dc0efbae431a12939fc15b2823db79b", "0xd09022c48298f268c2c431dadb9ca4c2534d9c1e",
        "0x0341a82e4a9d0a4312f3ee2ac9c6958512891342", "0x3d493c51a916f86d6d1c04824b3a4431e61a3ca3",
        "0xe1c15164dcfe49431f8421b5a311a829cf0904f3"};

    CStringArray allAddrs;
    for (auto addr : addrsInSet) {
        allAddrs.push_back(addr);
    }
    for (auto addr : addrsNotInSet) {
        allAddrs.push_back(addr);
    }

    // First, we create a bloom filter and add each address in the addrsInSet set to it.
    CBloomFilter bloom;
    for (auto addr : addrsInSet) {
        addAddrToBloom2(bloom, addr);
    }

    // Next we show a few statistics about the bloom filter
    CBloomStats stats(bloom);
    stats.Report();

    // Next we show each address being processed by the bloom filter algorithm and accumulate the list
    // of all the bits that we expect to be lit in the bloom filter.
    cout << endl;
    for (auto addr : allAddrs) {
        if (addr == "0x1296d3bb6dc0efbae431a12939fc15b2823db79b")
            cout << endl;

        CStringArray parts;
        CUintArray values;
        CUintArray bits;
        bloom_t addrBloom = addr_2_Bloom2(addr, parts, values, bits);
        const char* STR_ADDR_AS_BLOOM =
            "[{ADDR}]\n"
            "\tParts:  [{PARTS}]\n"
            "\tValues: [{VALUES}]\n"
            "\tBits:   [{BITS}]";

        ASSERT(parts.size() == values.size() && parts.size() == bits.size() && parts.size() == 5);
        ostringstream pStream, vStream, bStream;
        for (size_t i = 0; i < parts.size(); i++) {
            pStream << parts[i] << ",";
            vStream << values[i] << ",";
            bStream << bits[i] << ",";
        }
        string_q line = STR_ADDR_AS_BLOOM;
        replace(line, "[{ADDR}]", addr);
        replace(line, "[{PARTS}]", substitute(trim(pStream.str(), ','), ",", ", "));
        replace(line, "[{VALUES}]", substitute(trim(vStream.str(), ','), ",", ", "));
        replace(line, "[{BITS}]", substitute(trim(bStream.str(), ','), ",", ", "));
        cout << line << endl;
    }

    // Next, we test to see if the address we've put into the bloom filter reporting true when queried
    cout << endl;
    for (auto addr : addrsInSet) {
        cout << addr << ": " << isInBloomFilter(bloom, addr) << endl;
    }

    // Next, we test that those that were not put in the filter report false (even though they could report true)
    for (auto addr : addrsNotInSet) {
        cout << addr << ": " << isInBloomFilter(bloom, addr) << endl;
    }

    // const char* STR_GOARRAY =
    //     "{\n"
    //     "    Addr:     common.HexToAddress(\"[{ADDR}]\"),\n"
    //     "    Values:   [5]uint32{[{VALUES}]},\n"
    //     "    Bits:     [5]uint32{[{BITS}]},\n"
    //     "    Member:   [{MEMBER}],\n"
    //     "    FalsePos: [{FALSEPOS}],\n"
    //     "},";

    // CBloomFilter blooms, addrBlooms;
    // for (size_t i = 0; i < addrs.size(); i++) {
    //     address_t addr = addrs[i];
    //     CUintArray values;
    //     CUintArray bits;
    //     bloom_t b = addr_2_Bloom2(addr, values, bits);
    //     if (!(i % 2)) {
    //         addAddrToBloom(blooms, addr);
    //     }
    //     addrBlooms.push_back(b);
    // }
    // cout << endl;

    // for (size_t i = 0; i < addrs.size(); i++) {
    //     bool hit = isInBloomFilter(blooms, addrBlooms[i]);
    //     LOG_INFO(addrs[i], "-", hit, ((i % 2) && hit ? "-fp" : ""));
    // }

    // for (auto bb : blooms) {
    //     for (size_t i = 0; i < BLOOM_WIDTH_IN_BYTES; i++) {
    //         std::bitset<8> x(bb.bits[i]);
    //         ostringstream os;
    //         os << x << endl;
    //         if (contains(os.str(), "1")) {
    //             cout << i << ": " << x << endl;
    //         }
    //     }
    // }

    etherlib_cleanup();
    return 0;
}

//---------------------------------------------------------------------------
bloom_t addr_2_Bloom2(const address_t& addr, CStringArray& parts, CUintArray& values, CUintArray& litBits) {
    bloom_t ret;
    for (size_t k = 0; k < K; k++) {
        string_q four_byte = extract(addr, 2 + (k * NIBBLE_WID), NIBBLE_WID);
        uint64_t value = str_2_Uint("0x" + four_byte);
        uint64_t bit = (value % BLOOM_WIDTH_IN_BITS);
        // cout << bit << " ";
        ret.lightBit(bit);
        parts.push_back(four_byte);
        values.push_back(value);
        litBits.push_back(bit);
    }
    return ret;
}

/*
func(bloom* BloomFilter)
getStats()(nBlooms uint64, nInserted uint64, nBitsLit uint64, nBitsNotLit uint64, sz uint64, bitsLit[] uint64) {
    bitsLit = []uint64{}
    sz += 4
    for _, bf := range bloom.Blooms {
            nInserted += uint64(bf.NInserted)
            sz += 4 + uint64(len(bf.Bytes))
            cnt := uint64(0)
            for _, b := range bf.Bytes {
                if b != 0 {
                        nBitsLit++ bitsLit = append(bitsLit, cnt)
                    }
                else {
                    nBitsNotLit++
                    // fmt.Printf("%d", b)
                }
                cnt++
            }
    }
    return
}
*/

CBloomStats::CBloomStats(const CBloomFilter& bloom) {
    nBlooms = nInserted = nBitsLit = nBitsNotLit = sz = 0;
    bitsLit.clear();
    nBlooms = bloom.size();
    for (const auto& b : bloom) {
        nInserted += b.nInserted;
        for (size_t i = 0; i < BLOOM_WIDTH_IN_BITS; i++) {
            if (isBitLit(i, b.bits)) {
                nBitsLit++;
                bitsLit.push_back(i);
            } else {
                nBitsNotLit++;
            }
        }
    }
}

void CBloomStats::Report(void) const {
    cout << "nBlooms:     " << nBlooms << " nInserted:   " << nInserted << " nBitsLit:    " << nBitsLit
         << " nBitsNotLit: " << nBitsNotLit << endl;
    cout << "bitsLit:" << endl;
    for (auto b : bitsLit) {
        cout << b << ",";
    }
    cout << endl;
}

static CUintArray unused;
static const bloom_t zeroBloom = addr_2_Bloom("0x0", unused);
//----------------------------------------------------------------------
bool addAddrToBloom2(CBloomFilter& blooms, const address_t& addr) {
    cout << endl << "--------------------------------" << endl;
    if (blooms.size() == 0) {
        blooms.push_back(zeroBloom);  // so we have something to add to
        cout << "Adds zero bloom" << endl;
    }

    CStringArray parts;
    CUintArray values;
    CUintArray litBits;
    bloom_t addrBloom = addr_2_Bloom2(addr, parts, values, litBits);
    for (auto bit : litBits) {
        blooms[blooms.size() - 1].lightBit(bit);
    }
    blooms[blooms.size() - 1].nInserted++;

    if (blooms[blooms.size() - 1].nInserted > MAX_ADDRS_IN_BLOOM)
        blooms.push_back(zeroBloom);

    return true;
}
