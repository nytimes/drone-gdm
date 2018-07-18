#!/bin/bash
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
      -t "nytimes/drone-gdm:v${major}-${img_lbl}" \
      -t "nytimes/drone-gdm:$TRAVIS_TAG" . ;
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker push "nytimes/drone-gdm:$TRAVIS_TAG"
  docker push "nytimes/drone-gdm:v${major}-${img_lbl}"

elif [ "$TRAVIS_BRANCH" == "master" ]; then
  docker build -t "nytimes/drone-gdm:develop" .
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker push "nytimes/drone-gdm:develop"
fi

# EOF

