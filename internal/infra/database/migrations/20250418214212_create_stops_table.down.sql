-- Drop indexes
DROP INDEX IF EXISTS "idx_stops_latitude";
DROP INDEX IF EXISTS "idx_stops_longitude";
DROP INDEX IF EXISTS "idx_stops_name";

-- Drop table
DROP TABLE IF EXISTS "stops";