#!/bin/sh

dir=/usr/local/app/taf/app_log/%s/%s

if [ -d $dir ]; then
  cd $dir
fi

if [ ! -f /bin/bash ]; then
  sh
else
  bash
fi
