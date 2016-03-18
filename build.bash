#!/usr/bin/env bash
#
# Caddy build script. Automates proper versioning.
#
# Usage:
#
#     $ ./build.bash [output_filename]
#
# Outputs compiled program in current directory.
# Default file name is 'ecaddy'.

set -euo pipefail

: ${output_filename:="${1:-}"}
: ${output_filename:="ecaddy"}

pkg=main
ldflags=()

# Timestamp of build
name="${pkg}.buildDate"
value=$(date --utc +"%F %H:%M:%SZ")
ldflags+=("-X" "\"${name}=${value}\"")

# Current tag, if HEAD is on a tag
name="${pkg}.gitTag"
set +e
value="$(git describe --exact-match HEAD 2>/dev/null)"
set -e
ldflags+=("-X" "\"${name}=${value}\"")

# Nearest tag on branch
name="${pkg}.gitNearestTag"
value="$(git describe --abbrev=0 --tags HEAD)"
ldflags+=("-X" "\"${name}=${value}\"")

# Commit SHA
name="${pkg}.gitCommit"
value="$(git rev-parse --short HEAD)"
ldflags+=("-X" "\"${name}=${value}\"")

# Summary of uncommitted changes
name="${pkg}.gitShortStat"
value="$(git diff-index --shortstat HEAD)"
ldflags+=("-X" "\"${name}=${value}\"")

# List of modified files
name="${pkg}.gitFilesModified"
value="$(git diff-index --name-only HEAD | tr "\n" "," | sed -e 's:,$::')"
ldflags+=("-X" "\"${name}=${value}\"")

set -x
go build \
  -ldflags "${ldflags[*]}" \
  -o "${output_filename}"
