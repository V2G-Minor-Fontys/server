-- name: AddController :exec
INSERT INTO controllers (id, cpu_id, firmware_version)
VALUES ($1, $2, $3);

-- name: AddControllerSettings :exec
INSERT INTO controller_settings (id, auto_start, heartbeat_rate)
VALUES ($1, $2, $3);

-- name: PairUserToController :execrows
UPDATE controllers
SET user_id = $2, updated_at = CURRENT_TIMESTAMP
WHERE cpu_id = $1;

-- name: GetPairedControllerByUserId :one
SELECT c.id, c.cpu_id, c.firmware_version, sqlc.embed(cs)
FROM controllers c
JOIN public.controller_settings cs ON c.id = cs.id
WHERE c.user_id = $1;

-- name: GetControllerByCpuId :one
SELECT c.id, c.cpu_id, c.firmware_version, sqlc.embed(cs)
FROM controllers c
         JOIN public.controller_settings cs ON c.id = cs.id
WHERE c.cpu_id = $1;


-- name: AddControllerTelemetry :exec
INSERT INTO controller_telemetry (id, controller_id, battery_voltage, battery_current, battery_power, battery_state, internal_temperature, module_temperature, radiator_temperature, grid_power_r, total_inverter_power, ac_active_power, load_power_r, total_load_power, total_energy_to_grid, daily_energy_to_grid, total_energy_from_grid, daily_energy_from_grid, work_mode, operation_mode, error_message, warning_code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22);

-- name: GetControllerTelemetryByControllerId :many
SELECT * FROM controller_telemetry
WHERE controller_id = $1;

-- name: UpdateControllerSettings :execrows
UPDATE controller_settings
SET auto_start = $2, heartbeat_rate = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;