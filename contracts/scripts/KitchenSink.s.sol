pragma solidity ^0.8.33;

import {KitchenSink} from "../src/KitchenSink.sol";
import {Script} from "@forge-std/Script.sol";

contract KitchenSinkScript is Script {
    KitchenSink public bcn;

    function setUp() public {}

    function run() public {
        vm.startBroadcast();
        bcn = new KitchenSink();
        vm.stopBroadcast();
    }
}
