# Copyright 2024-2025 the original author or authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

##@ Tools

.PHONY: tools
tools: ## Install ci tools

	@$(LOG_TARGET)
	go version
	python --version
	node --version
	npm --version

	@echo "Installing licenses-eyes"
	go install github.com/apache/skywalking-eyes/cmd/license-eye@v0.6.1-0.20250110091440-69f34abb75ec

	@echo "Installing markdownlint-cli"
	npm install markdownlint-cli --global

	@echo "Installing linkinator"
	npm install linkinator --global

	@echo "Installing codespell"
	pip install codespell

	@echo "Installing yamllint"
	pip install yamllint==1.35.1

	@echo "Installing yamlfmt"
	go install github.com/google/yamlfmt/cmd/yamlfmt@latest

	@echo "Installing golangci-lint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.2.0
