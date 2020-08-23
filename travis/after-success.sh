#!/bin/bash
#===============================================================================
#
# drone-gdm/travis/after-success.sh:
#   Executed when main travis script succeeds.
#   Used to build and (optionally) push the docker image.
#
#-------------------------------------------------------------------------------
DRONE_GDM_IMAGE_NAME="${DRONE_GDM_IMAGE_NAME:-"${TRAVIS_REPO_SLUG}"}"

#------------------------
# PR Handling:
#------------------------
# Don't bother building the docker image for pull requests.
# TODO: why would we not do this again?
if [ "${TRAVIS_PULL_REQUEST}" != "false" ]; then
    exit 0
fi

#------------------------
# Tagged Builds:
#------------------------
if [ -n "$TRAVIS_TAG" ]; then
  major="${TRAVIS_TAG%%.*}"
  last_char="${TRAVIS_TAG:((${#TRAVIS_TAG}-1)):1}"
  if [ "${last_char}" == "a" ]; then
    img_lbl="alpha"
  elif [ "${last_char}" == "b" ]; then
    img_lbl="beta"
  else
    img_lbl="stable"
  fi

  docker build \
      -t "${DRONE_GDM_IMAGE_NAME}:v${major}-${img_lbl}" \
      -t "${DRONE_GDM_IMAGE_NAME}:$TRAVIS_TAG" . ;
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker push "${DRONE_GDM_IMAGE_NAME}:$TRAVIS_TAG"
  docker push "${DRONE_GDM_IMAGE_NAME}:v${major}-${img_lbl}"

#------------------------
# Mainline builds:
#------------------------
elif [ "$TRAVIS_BRANCH" == "main" ]; then
  docker build -t "${DRONE_GDM_IMAGE_NAME}:latest" .
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker push "${DRONE_GDM_IMAGE_NAME}:latest"

#------------------------
# Development builds:
#------------------------
else
  img_name="${DRONE_GDM_IMAGE_NAME}:test-image"
  printf "%s\n" "Building docker image %s...\n" \
      "${img_name}"
  docker build -t "${img_name}" .
  echo "Docker push skipped for branch \"${TRAVIS_BRANCH}\""
fi

# EOF

