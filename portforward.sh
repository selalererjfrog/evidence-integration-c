#!/bin/bash

echo "Starting port-forwarding for ArgoCD, Quotopia UI, and QuoteOfDay..."

# Function to run a port-forward command and restart it if it fails
run_port_forward() {
  local service_name=$1
  local local_port=$2
  local service_port=$3
  local namespace=$4

  while true; do
    echo "Starting port-forward for $service_name on port $local_port..."
    kubectl -n "$namespace" port-forward "$service_name" "$local_port":"$service_port"

    # The command above will exit if the connection is broken (e.g., pod replaced).
    echo "Port-forward for $service_name exited. Restarting in 2 seconds..."
    sleep 2
  done
}

# Run each port-forward command in the background
run_port_forward svc/argocd-server 8080 443 argocd &
PID1="$!"
run_port_forward service/quotopia-ui-service 9000 80 quotopia &
PID2="$!"
run_port_forward service/quoteofday-service 8001 8001 quotopia &
PID3="$!"

killcommands() {
    echo "Killing the port-forward commands..."
    kill "$PID1"
    kill "$PID2"
    kill "$PID3"
    echo "Done"
}

trap killcommands SIGINT SIGTERM

echo "All port-forwarding processes are running in the background. Press Ctrl+C to stop them."

wait


