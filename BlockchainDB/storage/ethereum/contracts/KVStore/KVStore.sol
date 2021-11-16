// SPDX-License-Identifier: MIT
pragma solidity ^0.4.24;

contract KVStore {
  event ItemSet(bytes32 key, bytes value);

  mapping (bytes32 => bytes) public items;

  function set(bytes32 key, bytes value) external {
    items[key] = value;
    emit ItemSet(key, value);
  }
}