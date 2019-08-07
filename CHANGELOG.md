Releases
========

2.0.5 - 2019/08/06
------------------
 - BUGFIX: add json/yaml annotations for config structure to ensure that configs loaded via `configfile` retain all of their properties
 - Fix ANSI term escape sequences for color after `ERROR:` messages + during command invocation

2.0.2 - 2019/07/29
------------------
 - Add support for [type-providers](https://cloud.google.com/deployment-manager/docs/configuration/type-providers/creating-type-provider)
 - Bump gcloud SDK version to `255.0.0`

2.0.1b - 2018/09/06 + 2.0.1 - 2018/10/03
----------------------------------------
 - Bump to cloud sdk `215.0.0`

2.0.0 - 2018/07/18
------------------
 - Updated support for external configuration files and more advanced templating
 - Fixed bug introduced in `1.2.Xa` series which _did not pass `deletePolicy` spec_

1.2.1a - 2018/05/09
-------------------
 - Use google-sdk base image

1.2.0a - 2018/05/09
-------------------
 - Add external configuration ability

1.1.0b - 2018/04/26
-------------------
 - Add drone 0.8 compatibility

1.0.10 - 2017/11/14
-------------------
 - Preserve additional error information for failed gcloud invocations

1.0.8b - 2017/08/28
-------------------
 - Use exact-match filter operator instead of pattern-match operator

1.0.7b - 2017/08/14
-------------------
 - Use `-q` option to prevent prompting on delete

1.0.6b - 2017/08/14
-------------------
Initial public BETA

(Previous stable is [andrewcanaday/drone-gdm:0.3.4b](https://hub.docker.com/r/andrewcanaday/drone-gdm/tags/))

