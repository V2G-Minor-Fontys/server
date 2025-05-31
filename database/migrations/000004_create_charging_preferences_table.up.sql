CREATE TYPE action_type AS ENUM ('remain', 'min', 'max', 'automatic');
CREATE TYPE weekday AS ENUM ('Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun');
CREATE TYPE ordinal AS ENUM ('1st', '2nd', '3rd', '4th', 'last');

CREATE TABLE IF NOT EXISTS charging_preferences (
  id              UUID PRIMARY KEY,
  name            TEXT NOT NULL,
  action_id       UUID NOT NULL REFERENCES actions(id),
  created_by      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  occurrence_id   UUID REFERENCES occurrences(id),
  date            DATE CHECK (occurrence_id IS NULL),
  priority        SMALLINT NOT NULL,
  enabled         BOOLEAN NOT NULL DEFAULT TRUE,

  CHECK (NOT (occurrence_id IS NOT NULL AND date IS NOT NULL)) -- either `occurrence_id` OR `date`
);

CREATE TABLE IF NOT EXISTS occurrences (
  id              UUID PRIMARY KEY,
  created_by      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  time            TIME NOT NULL,
  repeat          INT, -- Number of repetitions
  until           DATE,
  day_of_week     weekday,
  nth_of_month    ordinal CHECK (day_of_week IS NOT NULL),
  day_of_month    SMALLINT CHECK (day_of_week IS NULL AND day_of_month BETWEEN 1 AND 31),

  CHECK (NOT (day_of_week IS NOT NULL AND day_of_month IS NOT NULL)), -- either `day_of_week (+ nth_of_month)` OR `day_of_month`
  CHECK (repeat IS NULL OR until IS NULL)
);

CREATE TABLE IF NOT EXISTS actions (
  id                        UUID PRIMARY KEY,
  created_by                UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type                      action_type NOT NULL,

  battery_charge            SMALLINT CHECK (type != 'automatic' AND battery_charge BETWEEN 0 AND 100),

  charge_if_price_below     DECIMAL(5, 3) CHECK (charge_if_price_below >= 0),
  discharge_if_price_above  DECIMAL(5, 3) CHECK (discharge_if_price_above >= 0),

  CHECK ((type != 'automatic' AND battery_charge IS NOT NULL))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_charging_preferences ON charging_preferences(created_by, name);
CREATE UNIQUE INDEX IF NOT EXISTS idx_preferences_enabled ON charging_preferences (enabled);

