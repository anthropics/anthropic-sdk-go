# Changelog

## 0.2.0-alpha.2 (2024-10-17)

Full Changelog: [v0.2.0-alpha.1...v0.2.0-alpha.2](https://github.com/anthropics/anthropic-sdk-go/compare/v0.2.0-alpha.1...v0.2.0-alpha.2)

### Features

* move pagination package from internal to packages ([#33](https://github.com/anthropics/anthropic-sdk-go/issues/33)) ([ee3edb1](https://github.com/anthropics/anthropic-sdk-go/commit/ee3edb16dcd406435ade212cb7553f75b161e297))


### Bug Fixes

* **beta:** merge betas param with the default value ([#32](https://github.com/anthropics/anthropic-sdk-go/issues/32)) ([9191ae0](https://github.com/anthropics/anthropic-sdk-go/commit/9191ae0b8a47c3c6ea9dfbb5073f1e66f5b4e1d8))


### Chores

* fix GetNextPage docstring ([#29](https://github.com/anthropics/anthropic-sdk-go/issues/29)) ([acf8009](https://github.com/anthropics/anthropic-sdk-go/commit/acf8009c886ec27cc07665b0377a2a3b3493c336))
* **internal:** update spec URL ([#31](https://github.com/anthropics/anthropic-sdk-go/issues/31)) ([240f1c3](https://github.com/anthropics/anthropic-sdk-go/commit/240f1c3d7e4dc145988d2f8d11e45ccfd255861e))

## 0.2.0-alpha.1 (2024-10-08)

Full Changelog: [v0.1.0-alpha.2...v0.2.0-alpha.1](https://github.com/anthropics/anthropic-sdk-go/compare/v0.1.0-alpha.2...v0.2.0-alpha.1)

### Features

* **api:** add message batches api ([#28](https://github.com/anthropics/anthropic-sdk-go/issues/28)) ([169eb3c](https://github.com/anthropics/anthropic-sdk-go/commit/169eb3c83d39126b4f9ec3a8d7f70c06466d9ef6))


### Bug Fixes

* **beta:** pass beta header by default ([#27](https://github.com/anthropics/anthropic-sdk-go/issues/27)) ([c79ba68](https://github.com/anthropics/anthropic-sdk-go/commit/c79ba6826c452ca1eeefd34db1638722fa942082))


### Refactors

* **types:** improve metadata type names ([#26](https://github.com/anthropics/anthropic-sdk-go/issues/26)) ([95f0266](https://github.com/anthropics/anthropic-sdk-go/commit/95f0266f62ba90590db68f1f98e41d80ea8f5388))
* **types:** improve tool type names ([#23](https://github.com/anthropics/anthropic-sdk-go/issues/23)) ([79e4d75](https://github.com/anthropics/anthropic-sdk-go/commit/79e4d75d26bbf2339841d27696477817c01a55fc))

## 0.1.0-alpha.2 (2024-10-04)

Full Changelog: [v0.1.0-alpha.1...v0.1.0-alpha.2](https://github.com/anthropics/anthropic-sdk-go/compare/v0.1.0-alpha.1...v0.1.0-alpha.2)

### Features

* **api:** support disabling parallel tool use ([#22](https://github.com/anthropics/anthropic-sdk-go/issues/22)) ([1d8c00b](https://github.com/anthropics/anthropic-sdk-go/commit/1d8c00b317536d77a26f74d0008e1a4760b17d2e))
* **client:** send retry count header ([#19](https://github.com/anthropics/anthropic-sdk-go/issues/19)) ([d1c8ea1](https://github.com/anthropics/anthropic-sdk-go/commit/d1c8ea1f84d05002705ac7aa4d47a5ba13c388e9))
* improve error message ([#15](https://github.com/anthropics/anthropic-sdk-go/issues/15)) ([98d1ffd](https://github.com/anthropics/anthropic-sdk-go/commit/98d1ffd29f97e85ea543f36ce104c341e729a7d2))


### Bug Fixes

* **requestconfig:** copy over more fields when cloning ([#17](https://github.com/anthropics/anthropic-sdk-go/issues/17)) ([d5e7415](https://github.com/anthropics/anthropic-sdk-go/commit/d5e741578ac0ff88db3b04564810321b18f4dd40))


### Chores

* **ci:** add CODEOWNERS file ([#12](https://github.com/anthropics/anthropic-sdk-go/issues/12)) ([71c33b8](https://github.com/anthropics/anthropic-sdk-go/commit/71c33b841dece97e77f04ea4feae3d586b59b0d6))


### Documentation

* improve and reference contributing documentation ([#21](https://github.com/anthropics/anthropic-sdk-go/issues/21)) ([7288df1](https://github.com/anthropics/anthropic-sdk-go/commit/7288df1e1e62401487bee0685f77119bae5287ee))
* update CONTRIBUTING.md ([#18](https://github.com/anthropics/anthropic-sdk-go/issues/18)) ([dcfcbf8](https://github.com/anthropics/anthropic-sdk-go/commit/dcfcbf8d07e3d7a7d6b6398d60724f38eca050a4))

## 0.1.0-alpha.1 (2024-08-14)

Full Changelog: [v0.0.1-alpha.0...v0.1.0-alpha.1](https://github.com/anthropics/anthropic-sdk-go/compare/v0.0.1-alpha.0...v0.1.0-alpha.1)

### Features

* **api:** add prompt caching beta ([#11](https://github.com/anthropics/anthropic-sdk-go/issues/11)) ([78f8c72](https://github.com/anthropics/anthropic-sdk-go/commit/78f8c7266dd98ef5f76d258f485ee284b7a0e590))
* publish ([5ff0ff8](https://github.com/anthropics/anthropic-sdk-go/commit/5ff0ff8cc5706c39a6dde75ae69d11c892ef8bb3))
* simplify system prompt ([#3](https://github.com/anthropics/anthropic-sdk-go/issues/3)) ([cd3fcef](https://github.com/anthropics/anthropic-sdk-go/commit/cd3fcefad20baef3c28375adf16ab266f97e7d94))
* simplify system prompt ([#4](https://github.com/anthropics/anthropic-sdk-go/issues/4)) ([85e1b34](https://github.com/anthropics/anthropic-sdk-go/commit/85e1b349619e7dd26c06ed0d9f566ddbbe80db2a))


### Bug Fixes

* deserialization of struct unions that implement json.Unmarshaler ([#6](https://github.com/anthropics/anthropic-sdk-go/issues/6)) ([a883a3a](https://github.com/anthropics/anthropic-sdk-go/commit/a883a3a8232dfca1ce8a139047a0356a3fd6015f))
* handle nil pagination responses when HTTP status is 200 ([#2](https://github.com/anthropics/anthropic-sdk-go/issues/2)) ([2bb2325](https://github.com/anthropics/anthropic-sdk-go/commit/2bb232557a9f75d58b7e7145c69771c927574dd3))
* message accumulation with union content block ([09457cb](https://github.com/anthropics/anthropic-sdk-go/commit/09457cb2ef8019cc23bcdefa0d3102e642d64b3d))


### Chores

* add back custom code ([106c404](https://github.com/anthropics/anthropic-sdk-go/commit/106c40466382daaa403e7f472647248e14d939d7))
* bump Go to v1.21 ([#7](https://github.com/anthropics/anthropic-sdk-go/issues/7)) ([928ed50](https://github.com/anthropics/anthropic-sdk-go/commit/928ed50c83154eb4f56575cf9f405a132000888e))
* **ci:** bump prism mock server version ([#5](https://github.com/anthropics/anthropic-sdk-go/issues/5)) ([0b326c6](https://github.com/anthropics/anthropic-sdk-go/commit/0b326c6b18effa222b8b03a17c1e562d0aedce1d))
* **examples:** minor formatting changes ([#8](https://github.com/anthropics/anthropic-sdk-go/issues/8)) ([4195c55](https://github.com/anthropics/anthropic-sdk-go/commit/4195c5541a1a517a3890bfe43eb84e3ddc496bfe))


### Documentation

* add examples to README ([df47298](https://github.com/anthropics/anthropic-sdk-go/commit/df4729897b782faeaa6a0795359ecf20b4a833ca))
