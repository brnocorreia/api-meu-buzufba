-- Create stops table
CREATE TABLE IF NOT EXISTS "stops" (
	"id" VARCHAR(255) PRIMARY KEY,
	"slug" VARCHAR(255) NOT NULL UNIQUE,
	"name" VARCHAR(255) NOT NULL,
	"latitude" DOUBLE PRECISION NOT NULL,
	"longitude" DOUBLE PRECISION NOT NULL,
	"security_rating" INT NOT NULL,
	"is_active" BOOLEAN DEFAULT true,
	"created_at" TIMESTAMPTZ DEFAULT now(),
	"updated_at" TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX "idx_stops_latitude" ON "stops" ("latitude");
CREATE INDEX "idx_stops_longitude" ON "stops" ("longitude");
CREATE INDEX "idx_stops_name" ON "stops" ("name");
CREATE INDEX "idx_stops_slug" ON "stops" ("slug");