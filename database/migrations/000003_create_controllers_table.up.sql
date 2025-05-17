CREATE TABLE IF NOT EXISTS controllers (
    id UUID PRIMARY KEY,
    cpu_id TEXT UNIQUE NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    firmware_version VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS controller_settings (
    id UUID PRIMARY KEY REFERENCES controllers(id) ON DELETE CASCADE,
    auto_start BOOLEAN NOT NULL DEFAULT TRUE,
    heartbeat_rate SMALLINT NOT NULL DEFAULT 5,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS controller_telemetry (
    id UUID PRIMARY KEY,
    controller_id UUID REFERENCES controllers(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    output_power INTEGER NOT NULL,
    soc SMALLINT NOT NULL,
    ev_discharging BOOLEAN NOT NULL
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_controllers_user_id ON controllers(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_controllers_cpu_id ON controllers(cpu_id);
CREATE INDEX IF NOT EXISTS idx_controller_telemetry_controller_id ON controller_telemetry(controller_id);
