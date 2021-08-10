#!/bin/sh

if [ "$GOOGLE_APPLICATION_CREDENTIALS_BASE64" != "" ]; then
  echo $GOOGLE_APPLICATION_CREDENTIALS_BASE64 | base64 -d > credentials.json
  export GOOGLE_APPLICATION_CREDENTIALS=credentials.json
fi

exec ./guardmech