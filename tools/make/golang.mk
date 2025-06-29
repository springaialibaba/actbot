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

.PHONY: fmt
fmt: ## Format the Go code
	go fmt ./...

.PHONY: vet
vet: ## Check the Go code for potential issues
	go vet ./...

.PHONY: test
test: ## Run the Go tests
	go test ./...

.PHONY: golangci-lint
golangci-lint: ## Run the Go linter
	@$(LOG_TARGET)
	./bin/golangci-lint run --timeout 10m

.PHONY: golangci-lint-fix
golangci-lint-fix: ## Run the Go linter with auto-fix
	@$(LOG_TARGET)
	golangci-lint run --fix --timeout 10m
