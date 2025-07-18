#!/bin/bash

# f:name=test-params f:verb=test
# f:params=secretRef:my-secret:SECRET_VAR|prompt:"Enter name":NAME_VAR|text:default-value:DEFAULT_VAR

echo "Secret: ${SECRET_VAR:0:8}..."
echo "Name: $NAME_VAR"
echo "Default: $DEFAULT_VAR"