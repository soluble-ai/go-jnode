# Copyright 2019 Soluble Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:alpine

# https://github.com/golangci/golangci-lint
RUN wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

RUN apk update && apk add git bash ncurses zip

# This is specific to our buildkite setup where the buildkite-agent
# runs as uid/gid 999 ... so here we'll create that user.  We do this
# so any stuff that gets written out will still be owned by the
# buildkite-agent user on the host.  (NB - alpine has a 'ping' group with
# gid 999 so we delete that first.)
RUN delgroup ping && adduser -D -u 999 -h /workdir buildkite
