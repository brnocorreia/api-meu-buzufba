CREATE TABLE IF NOT EXISTS "stops" (
	"id" PRIMARY KEY INT AUTO_INCREMENT,
	"slug" VARCHAR(255) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"latitude" DOUBLE PRECISION NOT NULL,
	"longitude" DOUBLE PRECISION NOT NULL,
	"security_rating" INT NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
	"updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
)

CREATE INDEX "idx_stops_latitude" ON "stops" ("latitude");
CREATE INDEX "idx_stops_longitude" ON "stops" ("longitude");
CREATE INDEX "idx_stops_name" ON "stops" ("name");