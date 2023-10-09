# 🧪 Integration Tests Module

This module is dedicated to integration tests for validating the functionality and robustness of our implementation.

## 🚀 Getting Started

Before diving into the tests, ensure to run the script `./scripts/generateJsonHash.go`. This will generate the required `example.json` file with test cases to be consumed by the tests.

## 📜 Tests Available

Currently, we have the following tests:

1. **triehash_test.go**: This is an integration test that loads data from `example.json`. It compares the hash from the JSON with a newly generated hash using this library to ensure consistency and correctness.

Ensure your environment is set up correctly and enjoy testing! ✅