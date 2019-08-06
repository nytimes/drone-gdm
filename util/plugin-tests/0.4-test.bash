#!/usr/bin/env bash
#------------------------------------------------------------------------------
# 0.4-test.bash: Demonstrate/test 0.4-style invocation
#------------------------------------------------------------------------------
thispath="${0}"
thisdir="${thispath%/*}"
DRONE_GDM="${DRONE_GDM:-"${thisdir}/../../drone-gdm"}"

${DRONE_GDM} -- "$( cat "${thisdir}/drone/0.4-test.json" )"

# EOF

