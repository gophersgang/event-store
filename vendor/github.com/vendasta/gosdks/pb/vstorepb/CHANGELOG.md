## 0.2.0
- Added the ability to define a backup configuration for a kind to `vstorepb.CreateKindRequest`, and read that configuration from `vstorepb.GetKindResponse`.
- Added backup configuration to the proto for a `vstorepb.Schema`.

## 0.1.3
- The folder name really mattered to protoc's importing mechanisms, avoiding breaking changes

## 0.1.2
- Do away with versioned package name to avoid breaking changes

## 0.1.1
- Fix package name to be relative to this repository to allow for easy code generation

## 0.1.0
- Minimal set of protos required to communicate with the vstore service

