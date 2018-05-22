#!/usr/bin/env bash
#------------------------------------------------------------------------------
# 0.4-test.bash: Demonstrate/test 0.4-style invocation
#------------------------------------------------------------------------------
thispath="${0}"
thisdir="${thispath%/*}"

${thisdir}/../drone-gdm -- "$( cat "${thisdir}/drone/0.4-test.json" )"

# EOF

