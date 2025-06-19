CREATE TYPE weekday AS ENUM ('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun');
CREATE TYPE ordinal AS ENUM ('1st', '2nd', '3rd', '4th', 'last');


CREATE TABLE IF NOT EXISTS one_time_occurrences (
  id                   UUID PRIMARY KEY,
  date_start           DATETIME NOT NULL,
  date_end             DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS regular_occurrences (
  id              UUID PRIMARY KEY,
  time_of_day     TIME NOT NULL,
  repeat          SMALLINT, -- Number of repetitions
  until           DATE,

  day_of_week     weekday,
  nth_of_month    ordinal,-- CHECK (day_of_week IS NOT NULL),
  
  day_of_month    SMALLINT,

  CHECK (day_of_week IS NOT NULL OR day_of_month IS NOT NULL), -- either `day_of_week (+ nth_of_month)` or `day_of_month``
  CHECK (day_of_week IS NULL OR day_of_month IS NULL), -- either `day_of_week (+ nth_of_month)` or `day_of_month``
  CHECK (repeat IS NULL OR until IS NULL)
);

CREATE TABLE IF NOT EXISTS charging_policies (
  id                        UUID PRIMARY KEY,
  min_charge                SMALLINT CHECK (min_charge >= 0 AND min_charge <= 100),
  max_charge                SMALLINT CHECK (max_charge >= 0 AND max_charge >= min_charge AND max_charge <= 100),
  charge_if_price_below     DECIMAL(5, 3) CHECK (charge_if_price_below >= 0),
  discharge_if_price_above  DECIMAL(5, 3) CHECK (discharge_if_price_above >= 0)
);


CREATE TABLE IF NOT EXISTS charging_preferences (
  id                   UUID PRIMARY KEY,
  user_id              UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name                 TEXT NOT NULL,
  priority             SMALLINT NOT NULL,
  enabled              BOOLEAN NOT NULL DEFAULT TRUE,
  
  charging_policy_id   UUID REFERENCES charging_policies(id),
  keep_battery_at      SMALLINT,

  one_time_occurrence_id  UUID REFERENCES one_time_occurrences(id),
  regular_occurrence_id   UUID REFERENCES regular_occurrences(id),

  CHECK (one_time_occurrence_id IS NOT NULL OR regular_occurrence_id IS NOT NULL),
  CHECK (one_time_occurrence_id IS NULL OR regular_occurrence_id IS NULL),

  CHECK (charging_policy_id IS NOT NULL OR keep_battery_at IS NOT NULL),
  CHECK (charging_policy_id IS NULL OR keep_battery_at IS NULL)
);


CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_charging_preferences ON charging_preferences(created_by, name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_preferences_enabled ON charging_preferences (enabled);

