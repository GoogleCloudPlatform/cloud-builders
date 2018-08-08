# Copyright 2016 Google, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

prepare_workspace() {
    # We record the original working directory because some tools need to
    # remember, and sometimes this function changes directory.
    export WORKSPACE="$PWD"

    # If $GOPATH is set, do nothing. The user has already set up the workspace.
    if [[ "$GOPATH" ]]; then
        # We allow relative paths here because of how the Cloud Build service
        # works. But, the Go tool does not allow relative paths. To accomodate
        # the Go tool, we prepend the CWD to $GOPATH if $GOPATH is not
        # absolute.
        if [[ "${GOPATH:0:1}" != '/' ]]; then
            export GOPATH="$PWD/$GOPATH"
        fi
        return 0
    fi


    # If there is a src directory, accept that as a GOPATH indicator.
    if [[ -d src ]]; then
        export GOPATH="$PWD"
        return 0
    fi

    # Otherwise, we need to set up $GOPATH and map the source in somehow.


    # If $PROJECT_ROOT is set, map the current directory into $GOPATH/src/$PROJECT_ROOT.
    if [[ -z "$PROJECT_ROOT" ]]; then
        # Nothing in the environment or directory structure gives a strong enough clue, let's inspect source.
        PROJECT_ROOT="$(/builder/go_workspace)"
        if [[ $? != 0 ]]; then
            echo '
Unable to determine the Go workspace structure.

To determine the workspace structure, this tool checks the following, in order:
1) Is $GOPATH set? Use that.
2) Is there a ./src directory? Set GOPATH=$PWD.
3) Is $PROJECT_ROOT set? Make a temporary workspace in GOPATH=./gopath, and link
   the contents of the current directory into ./gopath/src/$PROJECT_ROOT/*.
4) Does a .go file in the current directory have a comment like
   // import "$PROJECT_ROOT"? Use the $PROJECT_ROOT found in the import
   comment instead of a provided $PROJECT_ROOT environment variable.

You are seeing this message because none of the checks succeeded.
'
            return 1
        fi
    fi

    shadow_workspace="./gopath/src/$PROJECT_ROOT"
    link_dir="$(dirname "$shadow_workspace")"

    echo "Creating shadow workspace and symlinking source into \"$shadow_workspace\"."

    mkdir -p "$link_dir" || return

    # We provide the -T option so that another workspace symlink is not placed inside
    # our workspace. If it fails, that's ok because it might be because a previous step
    # already created the symlink.
    ln  -s $PWD "$shadow_workspace" -T 2> /dev/null || stat "$shadow_workspace" 2> /dev/null || return

    export GOPATH="$PWD/gopath"

    # The tools that use this workspace preparation should be run from inside
    # the shadow workspace. If they are run from their original locations,
    # relative targets (eg. ".") will not work properly. Any files read or
    # written from inside the shadow workspace will have the same affect as if
    # they were read or written from the original directory, except that the new
    # working directory will be properly couched in GOPATH.
    cd "$shadow_workspace"
}
