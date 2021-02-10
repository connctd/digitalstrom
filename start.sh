#!/bin/bash

CONSOLE_FILE=console/console

if [[ ! -f $CONSOLE_FILE ]];
then
  printf './%s is not existing, will be build now.\n' "$CONSOLE_FILE"
  cd console
  go build
  cd ..
fi

./$CONSOLE_FILE -at a49c2cdd96b62681bdf846b54f8fcc23cda575c59d81e6d63a6e5085347eb8a2
