#!/bin/bash

find . -type f -name "*.sh" -exec stat -c "%A %a %n" {} \;
