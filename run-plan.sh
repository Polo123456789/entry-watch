#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

ITERATION_COUNT=0
MAX_ITERATIONS=30

PLAN="$1"
FINISH_FILE="plan-finished"

if [ -z "$PLAN" ]; then
    echo "Usage: $0 <plan-file>"
    exit 1
fi

if [ ! -f "$PLAN" ]; then
    echo "Plan file '$PLAN' does not exist."
    exit 1
fi

rm -f $FINISH_FILE

while [ ! -f $FINISH_FILE ]; do
    opencode -m github-copilot/gpt-5-mini run Use the plan in $PLAN to implement auth. Choose whatever you think is the most important task, and implement it, and update its status so future sessions dont reimplement it. Do regular commits once you have verified that the code is correct.
    ITERATION_COUNT=$((ITERATION_COUNT + 1))
    if [ $ITERATION_COUNT -ge $MAX_ITERATIONS ]; then
        echo "Reached maximum iteration count of $MAX_ITERATIONS. Exiting."
        exit 1
    fi
done

echo "Plan finished successfully."
echo "Iterated $ITERATION_COUNT times."
