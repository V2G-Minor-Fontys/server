-- name: AddController :exec
INSERT INTO controllers (id, cpu_id, firmware_version)
VALUES ($1, $2, $3);

-- name: AddControllerSettings :exec
INSERT INTO controller_settings (id, auto_start, heartbeat_rate)
VALUES ($1, $2, $3);

-- name: PairUserToController :exec
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
INSERT INTO controller_telemetry (id, controller_id, output_power, soc, ev_discharging)
VALUES ($1, $2, $3, $4, $5);

-- name: GetControllerTelemetryByControllerId :many
SELECT * FROM controller_telemetry
WHERE controller_id = $1;

-- name: UpdateControllerSettings :exec
UPDATE controller_settings
SET auto_start = $2, heartbeat_rate = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;