#!/usr/bin/env bash
#===============================================================================
#
# util/gen-dockerhub-readme.sh: Convert repo README into dockerhub project page.
#
#-------------------------------------------------------------------------------

#--------------------------------------
#          Script Context:
#--------------------------------------
__thispath="${0}"
__thisdir="${__thispath%/*}"
__thisname="${__thispath##*/}"
README_IN="${__thisdir}/../README.md"
README_OUT="${__thisdir}/../README.dockerhub.md"


#--------------------------------------
#              main:
#--------------------------------------
main() {
    cat "${README_IN}" \
        | sed 's^(\./^(https://github.com/nytimes/drone-gdm/blob/main/^g' \
        | sed 's/:[a-z][a-z_]\{1,\}[a-z]://g' \
        | tee "${README_OUT}"
}

main
# EOF

