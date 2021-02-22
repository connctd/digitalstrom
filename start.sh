#!/bin/bash

CONSOLE_FILE=console/console

if [[ ! -f $CONSOLE_FILE ]];
then
  printf './%s is not existing, will be build now.\n' "$CONSOLE_FILE"
  cd console
  go build
  cd ..
fi

./$CONSOLE_FILE -at <your application>
