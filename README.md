<!--
  ~ Copyright 2024-2025 the original author or authors.
  ~
  ~ Licensed under the Apache License, Version 2.0 (the "License");
  ~ you may not use this file except in compliance with the License.
  ~ You may obtain a copy of the License at
  ~
  ~     https://www.apache.org/licenses/LICENSE-2.0
  ~
  ~ Unless required by applicable law or agreed to in writing, software
  ~ distributed under the License is distributed on an "AS IS" BASIS,
  ~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  ~ See the License for the specific language governing permissions and
  ~ limitations under the License.
-->

## Robot working on GitHub Actions

:rocket: Goals of the first phase

* [X] `/retest` in PR

* [X] `/[un] assign` in Issue

* [X] `/sync` in Issue

* [X] `/[un] area` in Issue

* [X] `/[un] kind` in Issue

:memo: Goals of the second phase

* [ ] `/lgtm` in PR 

* [ ] `/[un] cc` in PR

### Quick Start

You can use it in GitHub workflow:

```yaml
name: Issue and PullRequest Command

on:
  issue_comment:
    types:
      - created

jobs:
  actbot:
    runs-on: ubuntu-22.04
    permissions:
      pull-requests: write
      contents: read
      issues: write
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - uses: ./
        name: Actbot Action
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```
