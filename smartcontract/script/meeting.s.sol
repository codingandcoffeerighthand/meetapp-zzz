// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/meeting.sol";

contract DeployMeeting is Script {
    function run() public {
        // Bắt đầu broadcast giao dịch từ tài khoản deployer
        vm.startBroadcast();
        new DAppMeeting();
        vm.stopBroadcast();
    }
}
