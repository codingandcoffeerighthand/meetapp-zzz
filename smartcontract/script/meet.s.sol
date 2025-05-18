// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.29;

import "forge-std/Script.sol";
import "../src/meet.sol";
import "../src/StringCompareLib.sol";

contract DeployMeeting is Script {
    function run() public {
        // Bắt đầu broadcast giao dịch từ tài khoản deployer
        vm.startBroadcast();
        new Meet();
        vm.stopBroadcast();
    }
}
