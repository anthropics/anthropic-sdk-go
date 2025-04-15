# Changelog

## 0.2.0-beta.4 (2025-04-15)

Full Changelog: [v0.2.0-beta.3...v0.2.0-beta.4](https://github.com/anthropics/anthropic-sdk-go/compare/v0.2.0-beta.3...v0.2.0-beta.4)

### Features

* **api:** extract ContentBlockDelta events into their own schemas ([#165](https://github.com/anthropics/anthropic-sdk-go/issues/165)) ([6d75486](https://github.com/anthropics/anthropic-sdk-go/commit/6d75486e9f524f5511f787181106a679e3414498))
* **api:** manual updates ([a92a382](https://github.com/anthropics/anthropic-sdk-go/commit/a92a382976d595dd32208109b480bf26dbbdc00f))
* **api:** manual updates ([59bd507](https://github.com/anthropics/anthropic-sdk-go/commit/59bd5071282403373ddca9333fafc9efc90a16d6))
* **client:** support param struct overrides ([#167](https://github.com/anthropics/anthropic-sdk-go/issues/167)) ([e0d5eb0](https://github.com/anthropics/anthropic-sdk-go/commit/e0d5eb098c6441e99d53c6d997c7bcca460a238b))
* **client:** support unions in query and forms ([#171](https://github.com/anthropics/anthropic-sdk-go/issues/171)) ([6bf1ce3](https://github.com/anthropics/anthropic-sdk-go/commit/6bf1ce36f0155dba20afd4b63bf96c4527e2baa5))


### Bug Fixes

* **client:** deduplicate stop reason type ([#155](https://github.com/anthropics/anthropic-sdk-go/issues/155)) ([0f985ad](https://github.com/anthropics/anthropic-sdk-go/commit/0f985ad54ef47849d7d478c84d34c7350a4349b5))
* **client:** return error on bad custom url instead of panic ([#169](https://github.com/anthropics/anthropic-sdk-go/issues/169)) ([b086b55](https://github.com/anthropics/anthropic-sdk-go/commit/b086b55f4886474282d4e2ea9ee3495cbf25ec6b))
* **client:** support multipart encoding array formats ([#170](https://github.com/anthropics/anthropic-sdk-go/issues/170)) ([611a25a](https://github.com/anthropics/anthropic-sdk-go/commit/611a25a427fc5303bb311fa4a2fec836d55b0933))
* **client:** unmarshal stream events into fresh memory ([#168](https://github.com/anthropics/anthropic-sdk-go/issues/168)) ([9cc1257](https://github.com/anthropics/anthropic-sdk-go/commit/9cc1257a67340e446ac415ec9ddddded24bb1f9a))


### Chores

* **docs:** doc improvements ([#173](https://github.com/anthropics/anthropic-sdk-go/issues/173)) ([aebe8f6](https://github.com/anthropics/anthropic-sdk-go/commit/aebe8f68afa3de4460cda6e4032c7859e13cda81))
* **docs:** update file uploads in README ([#166](https://github.com/anthropics/anthropic-sdk-go/issues/166)) ([a4a36bf](https://github.com/anthropics/anthropic-sdk-go/commit/a4a36bfbefa5a166774c23d8c5428fb55c1b4abe))
* **internal:** remove CI condition ([#160](https://github.com/anthropics/anthropic-sdk-go/issues/160)) ([adfa1e2](https://github.com/anthropics/anthropic-sdk-go/commit/adfa1e2e349842aa88262af70b209d1a59dbb419))
* **internal:** update config ([#157](https://github.com/anthropics/anthropic-sdk-go/issues/157)) ([46f0194](https://github.com/anthropics/anthropic-sdk-go/commit/46f019497bd9533390c4b9f0ebee6863263ce009))

## 0.2.0-beta.3 (2025-03-27)

Full Changelog: [v0.2.0-beta.2...v0.2.0-beta.3](https://github.com/anthropics/anthropic-sdk-go/compare/v0.2.0-beta.2...v0.2.0-beta.3)

### Chores

* add hash of OpenAPI spec/config inputs to .stats.yml ([#154](https://github.com/anthropics/anthropic-sdk-go/issues/154)) ([76b91b5](https://github.com/anthropics/anthropic-sdk-go/commit/76b91b56fbf42fe8982e7b861885db179b1bdcc5))
* fix typos ([#152](https://github.com/anthropics/anthropic-sdk-go/issues/152)) ([1cf6a6a](https://github.com/anthropics/anthropic-sdk-go/commit/1cf6a6ae25231b88d2eedbe0758f1281cbe439d8))

## 0.2.0-beta.2 (2025-03-25)

Full Changelog: [v0.2.0-beta.1...v0.2.0-beta.2](https://github.com/anthropics/anthropic-sdk-go/compare/v0.2.0-beta.1...v0.2.0-beta.2)

### Bug Fixes

* **client:** use raw json for tool input ([1013c2b](https://github.com/anthropics/anthropic-sdk-go/commit/1013c2bdb87a27d2420dbe0dcadc57d1fe3589f2))


### Chores

* add request options to client tests ([#150](https://github.com/anthropics/anthropic-sdk-go/issues/150)) ([7c70ae1](https://github.com/anthropics/anthropic-sdk-go/commit/7c70ae134a345aff775694abcad255c76e7dfcba))

## 0.2.0-beta.1 (2025-03-25)

Full Changelog: [v0.2.0-alpha.13...v0.2.0-beta.1](https://github.com/anthropics/anthropic-sdk-go/compare/v0.2.0-alpha.13...v0.2.0-beta.1)

### ⚠ BREAKING CHANGES

* **api:** migrate to v2

### Features

* add SKIP_BREW env var to ./scripts/bootstrap ([#137](https://github.com/anthropics/anthropic-sdk-go/issues/137)) ([4057111](https://github.com/anthropics/anthropic-sdk-go/commit/40571110129d5c66f171ead36f5d725663262bc4))
* **api:** migrate to v2 ([fcd95eb](https://github.com/anthropics/anthropic-sdk-go/commit/fcd95eb8f45d0ffedcd1e47cd0879d7e66783540))
* **client:** accept RFC6838 JSON content types ([#139](https://github.com/anthropics/anthropic-sdk-go/issues/139)) ([78d17cd](https://github.com/anthropics/anthropic-sdk-go/commit/78d17cd4122893ba62b1e14714a1da004c128344))
* **client:** allow custom baseurls without trailing slash ([#135](https://github.com/anthropics/anthropic-sdk-go/issues/135)) ([9b30fce](https://github.com/anthropics/anthropic-sdk-go/commit/9b30fce0a71a35910315e02cd3a2f2afc1fd7962))
* **client:** improve default client options support ([07f82a6](https://github.com/anthropics/anthropic-sdk-go/commit/07f82a6f9e07bf9aadf4ca150287887cb9e75bc4))
* **client:** improve default client options support ([#142](https://github.com/anthropics/anthropic-sdk-go/issues/142)) ([f261355](https://github.com/anthropics/anthropic-sdk-go/commit/f261355e497748bcb112eecb67a95d7c7c5075c0))
* **client:** support v2 ([#147](https://github.com/anthropics/anthropic-sdk-go/issues/147)) ([6b3af98](https://github.com/anthropics/anthropic-sdk-go/commit/6b3af98e02a9b6126bd715d43f83b8adf8b861e8))


### Chores

* **docs:** clarify breaking changes ([#146](https://github.com/anthropics/anthropic-sdk-go/issues/146)) ([a2586b4](https://github.com/anthropics/anthropic-sdk-go/commit/a2586b4beb2b9a0ad252e90223fbb471e6c25bc1))
* **internal:** codegen metadata ([ce0eca2](https://github.com/anthropics/anthropic-sdk-go/commit/ce0eca25c6a83fca9ececccb41faf04e74566e2d))
* **internal:** remove extra empty newlines ([#143](https://github.com/anthropics/anthropic-sdk-go/issues/143)) ([2ed1584](https://github.com/anthropics/anthropic-sdk-go/commit/2ed1584c7d80fddf2ef5143eabbd33b8f1a4603d))


### Refactors

* tidy up dependencies ([#140](https://github.com/anthropics/anthropic-sdk-go/issues/140)) ([289cc1b](https://github.com/anthropics/anthropic-sdk-go/commit/289cc1b007094421305dfc4ef01ae68bb2d50ee5))
