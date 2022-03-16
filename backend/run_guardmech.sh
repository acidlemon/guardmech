#!/bin/bash

if [ "$GOOGLE_APPLICATION_CREDENTIALS_BASE64" != "" ]; then
  echo $GOOGLE_APPLICATION_CREDENTIALS_BASE64 | base64 -d > credentials.json
  export -n GOOGLE_APPLICATION_CREDENTIALS_BASE64
  export GOOGLE_APPLICATION_CREDENTIALS=credentials.json
fi

# replace string on dist/index.html
if [ "$GUARDMECH_MOUNT_PATH" != "" ]; then
  perl -pi -e "s@href=\"/guardmech/admin@href=\"$GUARDMECH_MOUNT_PATH/guardmech/admin@g" dist/index.html
  perl -pi -e "s@src=\"/guardmech/admin@src=\"$GUARDMECH_MOUNT_PATH/guardmech/admin@g" dist/index.html
fi

exec ./guardmech
