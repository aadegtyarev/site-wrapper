#!/usr/bin/env bash
set -euo pipefail

docker volume create site-wrapper-android-sdk >/dev/null
docker volume create site-wrapper-gradle-cache >/dev/null

docker run --rm \
  --volume site-wrapper-android-sdk:/sdk \
  --volume site-wrapper-gradle-cache:/gradle-cache \
  --env ANDROID_HOME=/sdk \
  gradle:9.8.2-jdk17 bash -lc '
    set -e
    if [ ! -x /sdk/cmdline-tools/latest/bin/sdkmanager ]; then
      curl -fsSL --retry 3 https://dl.google.com/android/repository/commandlinetools-linux-15859902_latest.zip -o /tmp/cmdline.zip
      echo "4e4c464f145a7512b57d088ac6c278c03c9eea610886b35a5e0804e74eedf583  /tmp/cmdline.zip" | sha256sum -c -
      mkdir -p /sdk/cmdline-tools/latest /tmp/site-wrapper-tools
      unzip -q /tmp/cmdline.zip -d /tmp/site-wrapper-tools
      mv /tmp/site-wrapper-tools/cmdline-tools/* /sdk/cmdline-tools/latest/
    fi
    yes | /sdk/cmdline-tools/latest/bin/sdkmanager --licenses >/dev/null
    /sdk/cmdline-tools/latest/bin/sdkmanager "platforms;android-37.2" "build-tools;36.0.0" "platform-tools"
    chmod 0777 /gradle-cache
  '

docker run --rm \
  --user "$(id -u):$(id -g)" \
  --volume site-wrapper-android-sdk:/sdk \
  --volume site-wrapper-gradle-cache:/gradle-cache \
  --volume "$(pwd)/android:/src" \
  --env ANDROID_HOME=/sdk \
  --env GRADLE_USER_HOME=/gradle-cache \
  --workdir /src \
  gradle:9.8.2-jdk17 gradle :app:testDebugUnitTest :app:lintDebug :app:assembleDebug --no-daemon --console=plain
