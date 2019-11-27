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
