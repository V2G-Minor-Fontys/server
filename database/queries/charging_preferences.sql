-- name: CreateChargingPreference :exec
INSERT INTO charging_preferences (
  id, created_by, name, occurrence_id, date, action_id, priority
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);

-- name: CreateOccurrence :exec
INSERT INTO occurrences (
  id, created_by, time, repeat, until, day_of_week, nth_of_month, day_of_month
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: CreateAction :exec
INSERT INTO actions (
  id, created_by, type, battery_charge,
  charge_if_price_below, discharge_if_price_above
) VALUES (
  $1, $2, $3, $4, $5, $6
);


-- name: GetChargingPreferencesByUserId :one
SELECT
  p.id AS preference_id,
  p.name,
  o.time,
  o.day_of_week,
  o.nth_of_month,
  o.day_of_month,
  o.repeat,
  o.until,
  a.type AS action_type,
  a.battery_charge,
  a.charge_if_price_below,
  a.discharge_if_price_above,
  p.priority,
  p.enabled
FROM charging_preferences p
LEFT JOIN occurrences o ON p.occurrence_id = o.id
JOIN actions a ON p.action_id = a.id
WHERE p.id = $1;

-- name: ListChargingPreferencesForUser :many
SELECT
  p.id AS preference_id,
  p.name,
  o.time,
  o.day_of_week,
  o.nth_of_month,
  o.day_of_month,
  o.repeat,
  o.until,
  a.type AS action_type,
  a.battery_charge,
  a.charge_if_price_below,
  a.discharge_if_price_above,
  p.priority,
  p.enabled
FROM charging_preferences p
LEFT JOIN occurrences o ON p.occurrence_id = o.id
JOIN actions a ON p.action_id = a.id
WHERE p.created_by = $1;


-- name: ListActionsForUser :many
SELECT * FROM actions
WHERE created_by = $1;

-- name: ListOccurrencesForUser :many
SELECT * FROM occurrences
WHERE created_by = $1;


-- name: DeleteChargingPreference :exec
DELETE FROM charging_preferences
WHERE id = $1;

-- name: DeleteOccurrence :exec
DELETE FROM occurrences
WHERE id = $1;

-- name: DeleteAction :exec
DELETE FROM actions
WHERE id = $1;
