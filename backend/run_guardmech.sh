#!/bin/sh

if [ "$GOOGLE_APPLICATION_CREDENTIALS_TEXT" != "" ]; then
  echo $GOOGLE_APPLICATION_CREDENTIALS_TEXT > credentials.json
  export GOOGLE_APPLICATION_CREDENTIALS=credentials.json
fi

exec ./guardmech