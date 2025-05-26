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

    battery_voltage FLOAT NOT NULL,
    battery_current FLOAT NOT NULL,
    battery_power FLOAT NOT NULL,
    battery_state SMALLINT NOT NULL,

    internal_temperature FLOAT NOT NULL,
    module_temperature FLOAT NOT NULL,
    radiator_temperature FLOAT NOT NULL,

    grid_power_r INT NOT NULL,
    total_inverter_power INT NOT NULL,
    ac_active_power INT NOT NULL,
    load_power_r INT NOT NULL,
    total_load_power INT NOT NULL,

    total_energy_to_grid FLOAT NOT NULL,
    daily_energy_to_grid FLOAT NOT NULL,
    total_energy_from_grid FLOAT NOT NULL,
    daily_energy_from_grid FLOAT NOT NULL,

    work_mode SMALLINT NOT NULL,
    operation_mode SMALLINT NOT NULL,

    error_message BIGINT NOT NULL,
    warning_code SMALLINT NOT NULL
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_controllers_user_id ON controllers(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_controllers_cpu_id ON controllers(cpu_id);
CREATE INDEX IF NOT EXISTS idx_controller_telemetry_controller_id ON controller_telemetry(controller_id);
