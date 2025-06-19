-- name: CreateChargingPreference :exec
INSERT INTO charging_preferences (
  id, user_id, name, priority, enabled, charging_policy_id, keep_battery_at, one_time_occurrence_id, regular_occurrence_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: CreateOneTimeOccurrence :exec
INSERT INTO one_time_occurrences (
  id, date_start, date_end
) VALUES ($1, $2, $3);

-- name: CreateRegularOccurrence :exec
INSERT INTO regular_occurrences (
  id,  time_of_day, repeat, until, day_of_week, nth_of_month, day_of_month
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateChargingPolicies :exec
INSERT INTO charging_policies (
  id, min_charge, max_charge, charge_if_price_below, discharge_if_price_above
) VALUES ($1, $2, $3, $4, $5);

-- name: ListChargingPreferencesForUser :many
SELECT
  p.id AS preference_id,
  p.user_id,
  p.name,
  p.priority,
  p.enabled,
  oo.date_start,
  oo.date_end,
  ro.time_of_day,
  ro.day_of_week,
  ro.nth_of_month,
  ro.day_of_month,
  ro.repeat,
  ro.until,
  cp.min_charge,
  cp.max_charge,
  cp.charge_if_price_below,
  cp.discharge_if_price_above
FROM charging_preferences p
LEFT JOIN regular_occurrences ro ON p.regular_occurrence_id = ro.id
LEFT JOIN one_time_occurrences oo ON p.one_time_occurrence_id = oo.id
LEFT JOIN charging_policies cp ON p.charging_policy_id = cp.id
WHERE p.user_id = $1;

-- name: DeleteChargingPreference :exec
DELETE FROM charging_preferences
WHERE id = $1;

-- name: DeleteOneTimeOccurrence :exec
DELETE FROM one_time_occurrences
WHERE id = $1;

-- name: DeleteRegularOccurrences :exec
DELETE FROM regular_occurrences
WHERE id = $1;

-- name: DeleteChargingPolicies :exec
DELETE FROM charging_policies
WHERE id = $1;
