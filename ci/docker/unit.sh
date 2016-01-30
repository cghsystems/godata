#!/bin/bash

#
# This script allows Concourse and the GOPATH to play nice together. It will
# move the code cloned by concourse into the GOPATH setup by the container running
# the job 
#

set -eux

target_dir=$GOPATH/src/github.com/cghsystems/godata

move_project_to_gopath() {
  mkdir -p ${target_dir}
  cp -r ./* ${target_dir}
}

install_dependencies() {
  godep restore ./...
}

execute_tests() {
  ginkgo -r --randomizeAllSpecs --failOnPending --randomizeSuites --race
}

#
# The Docker container that run the scripts will default $TMPDIR to /tmp which
# is where Concourse checks out code. Ginkgo's gexec methods will create
# artifacts and remove the contents of /tmp after every test run hence we need
# to redirect $TMPDIR
#
setup_tmp_dir() {
  local tmp_dir=/ccp/tmp
  mkdir -p ${tmp_dir}
  export TMPDIR=${tmp_dir}
}
 
main() {
  setup_tmp_dir
  move_project_to_gopath
  cd ${target_dir}
  install_dependencies
  execute_tests
}

main
