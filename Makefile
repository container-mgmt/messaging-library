#
# Copyright (c) 2018 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# This Makefile is just a wrapper calling the 'build.py' script, for those used
# to just run 'make'.

.PHONY: build
build:
	./build.py build

.PHONY: binaries
binaries:
	./build.py binaries

.PHONY: lint
lint:
	./build.py lint

.PHONY: test
test:
	./build.py test

.PHONY: fmt
fmt:
	./build.py fmt

.PHONY: clean
clean:
	rm -rf .gopath vendor
