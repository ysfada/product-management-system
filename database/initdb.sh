# #!/bin/bash
set -e

echo "Running initdb.sh..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
-- ADD YOUR INITIAL SQL HERE
EOSQL

echo "Finished running initdb.sh..."
