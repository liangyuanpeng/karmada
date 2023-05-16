#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# This script is working for github action.
# Today, it is only work for rerun workflow when the github action CI is failed.

# Usage:
# hack/retest.sh
# Environments:
#   ISSUE_COMMENT: Only support /retest to rerun workflow,Just confirm that your issue comment is include "/retest"
#   PR_NUM: Your PR number
#   GH_TOKEN: Github Token for run workflow
#   REPO: Github Repo, karmada-io/karmada by default

function rerun_workflow(){

  PR_NUM=${PR_NUM:-0}

  if [ $PR_NUM -le 0 ];then
    echo "Invalid pr num:"$PR_NUM
    exit 0
  fi

  REPO=${REPO:-"liangyuanpeng/karmada"}

  EVENT=`gh api repos/$REPO/pulls/$PR_NUM `

  EVENT_HEAD=$(echo $EVENT | jq .head )

  SHA=$(echo $EVENT_HEAD | jq .sha | sed 's/\"//g')
  BRANCH=$(echo $EVENT_HEAD | jq .ref | sed 's/\"//g')
  ACTOR=$(echo $EVENT_HEAD | jq .user.login | sed 's/\"//g')

  PR_TITLE=$(echo $EVENT | jq .title | sed 's/\"//g' )
  PR_URL=$(echo $EVENT | jq .html_url | sed 's/\"//g' )

  datas=$(gh api "repos/$REPO/actions/runs?actor="$ACTOR"&branch="$BRANCH"&status=failure" | jq ".workflow_runs[] | select(.head_sha==\"$SHA\") | [{id,name}] ")

  echo $datas | jq -r '.[] | "\(.id)\t\(.name)"' | while read -r id name; do
      echo -e "Reruning workflow...\nPR:"$PR_TITLE "\nURL:"$PR_URL "\nID:$id \nWorkflowRunName:$name \n==============================="
      gh run rerun -R $REPO $id
  done
}

ISSUE_COMMENT=${ISSUE_COMMENT:-""}

if [ -z "$ISSUE_COMMENT" ];then
    echo "Have not issue comment,exit..."
    exit -1
fi

echo -e $ISSUE_COMMENT |sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//'| while read line
do
  echo "line:"$line "line don"
  if [[ "$line" == "/retest-failed" ]]; then
    echo "Matching /retest-failed and rerun workflow..."
    rerun_workflow
    break
  fi
done
