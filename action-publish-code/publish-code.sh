#!/bin/bash

# Publishes a subdirectory in a given git repository as a standalone repository of it's own.

set -euo pipefail
set -o errexit
set -o nounset

if [ "$#" = 0 ]; then
	echo "usage: $0 <org>/<repo> <dir1>[:<repo1>] [<dirN>[:<repoN>]]..." 1>&2
	exit 1
fi

if [[ "$1" != *"/"* ]]; then
	echo "error: first argument must be in the form <owner>/<repo>" 1>&2
	echo "usage: $0 <org>/<repo> <dir1>[:<repo1>] [<dirN>[:<repoN>]]..." 1>&2
	exit 1
fi

# Parse the source repo from the arguments:
org="${1%/*}"
origin_repo="${1##*/}"
shift

# Create a temporary directory into which to clone the repos:
TMPDIR=$(mktemp -d)
function cleanup() {
	echo "Deleting '${TMPDIR}' ..." 1>&2
	rm -rf "${TMPDIR}"
}
trap cleanup EXIT INT # Clean up on exit

# Clone the origin repo:
echo "Cloning '${org}/${origin_repo}' ..." 1>&2
git clone -b main --single-branch "git@github.com:${org}/${origin_repo}.git" "${TMPDIR}/${origin_repo}"

# Check that the directories we want to publish exist:
echo -e "\n"
cd "${TMPDIR}/${origin_repo}"
for arg in "$@"; do
	dir="${arg%:*}"
	repo="${dir##*/}"
	if [[ "$arg" != *":"* ]]; then
		repo="${arg##*:}"
	fi

	echo "Validating $arg"
	# Just to be extra safe, make sure we're not trying to publish to origin repo itself
	if [ "${repo}" = "${origin_repo}" ]; then
		echo "Error: Cannot publish to the source repo: '${origin_repo}'" 1>&2
		exit 1
	fi

	# Verify the directory exists:
	if [ ! -d "${dir}" ]; then
		echo "Error: Directory '${dir}' does NOT exist, cannot publish. Aborting." 1>&2
		exit 1
	fi
done

for arg in "$@"; do
	dir="${arg%:*}"
	repo="${dir##*/}"
	if [[ "$arg" = *":"* ]]; then
		repo="${arg##*:}"
	fi
	echo "Publishing dir '${dir}' to repo '${repo}' ..."

	set -o xtrace # Print commands as they are executed

	# Remove everything we don't want from the source repo:
 
	# Only keep the tags belonging to the repo we care about
	git tag -d $(git tag -l | grep -v "${repo}/*")
	# TODO: Rewrite using https://github.com/newren/git-filter-repo since filter-branch is no longer
	# recommended by git.
	FILTER_BRANCH_SQUELCH_WARNING=1 git filter-branch \
		--tag-name-filter "grep '^${repo}/' | cut -f 2- -d '/'" \
		--subdirectory-filter "${dir}" --prune-empty -- --all

	git clone "git@github.com:jetify-examples/${repo}.git" "${TMPDIR}/${repo}"
	pushd "${TMPDIR}/${repo}"

	# Here is the trick. Connect your source repository as a remote using a local reference.
	git remote add "${origin_repo}" "../${origin_repo}/"

	# After that simply fetch the remote source, create a branch and merge it with the destination repository in usual way
	git fetch "${origin_repo}"
	git branch "${origin_repo}" "remotes/${origin_repo}/main"

	git merge "${origin_repo}" --allow-unrelated-histories --no-edit --ff-only

	# This is pretty much it, all your code and history were moved from one repository to another.
	# All you need is to clean up a bit and push the changes to the server

	git remote rm "${origin_repo}"
	git branch -d "${origin_repo}"

	git push origin main --follow-tags

	popd # Back to origin repo
	set +o xtrace # Turn off tracing

	# Undo the filtering so we can re-use the source repo for another rewrite.
	git for-each-ref --format="update %(refname:lstrip=2) %(objectname)" refs/original/ | git update-ref --stdin
	git for-each-ref --format="delete %(refname) %(objectname)" refs/original/ | git update-ref --stdin
	git fetch --tags --force  # Restore all tags
	git reset --hard HEAD

	echo "[DONE] Published dir '${dir}' to repo '${repo}' ..."
	echo -e "\n"
done