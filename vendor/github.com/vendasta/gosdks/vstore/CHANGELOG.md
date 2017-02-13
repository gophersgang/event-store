# VStore Client Changelog

## 0.6.1
- Fix a bug where the lookup stub was returning too many results relative to specified pageSize.

## 0.6.0
- Added in memory happy path stub for vstore.

## 0.5.0
* Added EnvFromEnvironmentVariable function to automatically calculate the environment and address that the VStore client should use

## 0.4.0
* Added GetSecondaryIndexName function to vstore client interface

## 0.3.0
* Added CloudSQL secondary indexes to the schema builder

## 0.2.0
* More documentation and examples. Not finished.
* Move changelog to its own file.
* Add lint task for all libs and lint VStore.

## 0.1.0
* Initial release.
* Available operations include Get, GetMulti, Lookup, and Transaction.
* Namespace and kind management is also available.
