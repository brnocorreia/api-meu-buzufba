CREATE TABLE IF NOT EXISTS routes (
    route_id SERIAL PRIMARY KEY,
    route_name VARCHAR(30) NOT NULL,
    trip_length INTEGER NOT NULL,
    departure_location VARCHAR(100) NOT NULL,
    arrival_location VARCHAR(100) NOT NULL,
    status_cd VARCHAR(10) NOT NULL DEFAULT 'A',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS departures (
    route_id INTEGER REFERENCES routes(id),
    departure_time TIME,
    obs VARCHAR(255),
    status_cd VARCHAR(10) NOT NULL DEFAULT 'A',
    PRIMARY KEY (route_id, departure_time)
);

CREATE TABLE IF NOT EXISTS locations (
    location_id SERIAL PRIMARY KEY,
    location_name VARCHAR(255) NOT NULL,
    location_lat DECIMAL(10, 8),
    location_lng DECIMAL(11, 8),
    PRIMARY KEY (route_id, location)
);

CREATE TABLE IF NOT EXISTS route_locations (
    route_id INTEGER REFERENCES routes(id),
    location_id INTEGER REFERENCES locations(id),
    obs VARCHAR(255),
    PRIMARY KEY (route_id, location_id)
);

CREATE TABLE stops (
    stop_id SERIAL PRIMARY KEY,
    stop_name VARCHAR(100) NOT NULL,
    stop_lat DECIMAL(10, 8),
    stop_lng DECIMAL(11, 8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS route_stops (
    route_id INTEGER REFERENCES routes(id),
    stop_id INTEGER REFERENCES stops(id),
    stop_type VARCHAR(10) NOT NULL,
    stop_order INTEGER NOT NULL,
    PRIMARY KEY (route_id, stop_type, stop_order)
);