#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status.

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

OBSERVABILITY_PROJECT_NAME="20250507_observability"
OBSERVABILITY_DIR_PATH="$SCRIPT_DIR/../$OBSERVABILITY_PROJECT_NAME"

# Check if the observability project directory exists
if [ ! -d "$OBSERVABILITY_DIR_PATH" ]; then
  echo "Error: Observability project directory '$OBSERVABILITY_PROJECT_NAME' not found at '$OBSERVABILITY_DIR_PATH'."
  echo "Please ensure 'minilab_backend' and '$OBSERVABILITY_PROJECT_NAME' are sibling directories."
  exit 1
fi

echo "Navigating to observability project: $OBSERVABILITY_DIR_PATH"
cd "$OBSERVABILITY_DIR_PATH"

echo ""
echo "Ensuring minilab-backend is rebuilt without cache to include latest dependencies..."
docker compose build --no-cache minilab-backend

echo ""
echo "Building (if necessary) and starting all services (including minilab_backend) with the observability stack..."
echo "This might take a while on the first run as images are downloaded and built."
docker compose up --build -d # --build will efficiently build other services if needed

echo ""
echo "--------------------------------------------------------------------"
echo "Minilab Backend (service name: minilab-backend) should be starting up."
echo "The full observability stack (Grafana, Prometheus, etc.) is also being started."
echo ""
echo "Key services & ports (on your host machine):"
echo "  Minilab Backend API: Expected at http://localhost:8080 (see docker-compose.yaml in $OBSERVABILITY_PROJECT_NAME)"
echo "  Minilab DB (PostgreSQL for minilab-backend): Exposed on host port 5433 (service name: minilab-db)"
echo "  Grafana (Visualization): http://localhost:3001"
echo "  Prometheus (Metrics): http://localhost:9090"
echo "  Loki (Logs via Grafana): Port 3100 (Loki's API port, data source for Grafana)"
echo ""
echo "To view logs for all services (from '$OBSERVABILITY_DIR_PATH'):"
echo "  docker compose logs -f"
echo "To view logs specifically for minilab-backend:"
echo "  docker compose logs -f minilab-backend"
echo ""
echo "To stop all services (from '$OBSERVABILITY_DIR_PATH'):"
echo "  docker compose down"
echo "--------------------------------------------------------------------" 