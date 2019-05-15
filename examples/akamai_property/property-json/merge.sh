#! /bin/bash
akamai property-manager merge -p $1
jq -n --arg result "$1" '{"result":"'$1'"}'
