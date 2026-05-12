CREATE TABLE route_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    city TEXT NOT NULL,
    stop_order INTEGER NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT route_stops_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT route_stops_city_not_empty CHECK (length(trim(city)) > 0),
    CONSTRAINT route_stops_stop_order_positive CHECK (stop_order > 0),
    CONSTRAINT route_stops_latitude_valid CHECK (latitude >= -90 AND latitude <= 90),
    CONSTRAINT route_stops_longitude_valid CHECK (longitude >= -180 AND longitude <= 180)
);

CREATE UNIQUE INDEX route_stops_route_order_unique_idx ON route_stops (route_id, stop_order);
CREATE INDEX route_stops_route_id_idx ON route_stops (route_id);
CREATE INDEX route_stops_route_city_idx ON route_stops (route_id, lower(city));
