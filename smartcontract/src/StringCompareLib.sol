// SPDX-License-Identifier: M
pragma solidity ^0.8.0;

library StringCompareLib {
    function safeCompare(string memory s1, string memory s2) internal pure returns (bool) {
        bytes memory b1 = bytes(s1);
        bytes memory b2 = bytes(s2);
        if (b1.length != b2.length) {
            return false;
        }
        return keccak256(b1) == keccak256(b2);
    }
}
