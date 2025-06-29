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

# All make targets should be implemented in tools/make/*.mk
# ====================================================================================================
# Supported Targets: (Run `make help` to see more information)
# ====================================================================================================

# This file is a wrapper around `make` so that we can force on the
# --warn-undefined-variables flag.  Sure, you can set
# `MAKEFLAGS += --warn-undefined-variables` from inside of a Makefile,
# but then it won't turn on until the second phase (recipe execution),
# and won't actually be on during the initial phase (parsing).
# See: https://www.gnu.org/software/make/manual/make.html#Reading-Makefiles

# Have everything-else ("%") depend on _run (which uses
# $(MAKECMDGOALS) to decide what to run), rather than having
# everything else run $(MAKE) directly, since that'd end up running
# multiple sub-Makes if you give multiple targets on the CLI.

_run:
	@$(MAKE) --warn-undefined-variables \
		-f tools/make/common.mk \
		-f tools/make/golang.mk \
		-f tools/make/linter.mk \
		-f tools/make/tools.mk \
		$(MAKECMDGOALS)

.PHONY: _run

$(if $(MAKECMDGOALS),$(MAKECMDGOALS): %: _run)
