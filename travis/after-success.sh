#!/bin/bash
ORG_NAME="${ORG_NAME:-"nytimes"}"

if [ "${TRAVIS_PULL_REQUEST}" != "false" ]; then
    exit 0
fi

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
      -t "${ORG_NAME}/drone-gdm:v${major}-${img_lbl}" \
      -t "${ORG_NAME}/drone-gdm:$TRAVIS_TAG" . ;
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker push "${ORG_NAME}/drone-gdm:$TRAVIS_TAG"
  docker push "${ORG_NAME}/drone-gdm:v${major}-${img_lbl}"

elif [ "$TRAVIS_BRANCH" == "main" ]; then
  docker build -t "${ORG_NAME}/drone-gdm:develop" .
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker push "${ORG_NAME}/drone-gdm:develop"
fi

# EOF

