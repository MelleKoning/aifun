#!/bin/bash

# Usage: ./run.sh [mode]
# Modes: base, ui, research, all

# Define valid modes
VALID_MODES=("base" "ui" "research" "all", "stop")

# Check if a parameter is provided
if [ -z "$1" ]; then
    echo "Error: No mode specified."
    echo "Usage: ./run.sh [$(IFS="|"; echo "${VALID_MODES[*]}")]"
    exit 1
fi

# Check if the provided parameter is valid
if [[ ! " ${VALID_MODES[@]} " =~ " $1 " ]]; then
    echo "Error: Invalid mode '$1'."
    echo "Valid modes are: ${VALID_MODES[*]}"
    exit 1
fi

# Execute based on parameter
case "$1" in
stop)
    # The "*" profile ensures everything in the file is stopped
    docker compose --profile "*" down
    echo "All services stopped."
    ;;
  base)
    docker compose up -d ollama
    ;;
  ui)
    docker compose --profile ui up -d
    ;;
  research)
    docker compose --profile research up
    ;;
  all)
    docker compose --profile ui --profile research up -d
    ;;
esac

echo "Sucessfully executed: $1"
