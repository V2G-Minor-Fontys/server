package charging_preferences

import (
	"errors"
	"fmt"
	"time"

	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ChargingPreference struct {
	Id                uuid.UUID          `json:"id"`
	UserId            uuid.UUID          `json:"user_id"`
	Name              string             `json:"name"`
	Priority          int16              `json:"priority"`
	Enabled           bool               `json:"enabled"`
	RegularOccurrence *RegularOccurrence `json:"regular_occurrence"`
	OneTimeOccurrence *OneTimeOccurrence `json:"one_time_occurrence"`
	ChargingPolicy    *ChargingPolicy    `json:"charging_policy"`
	KeepChargeAt      *int               `json:"keep_charge_at"`
}

type RegularOccurrence struct {
	TimeOfDay  time.Time `json:"time_of_day"`
	Repeat     int       `json:"repeat"`
	Until      string    `json:"until"`
	DayOfWeek  string    `json:"day_of_week"`
	NthOfMonth string    `json:"nth_of_month"`
	DayOfMonth int       `json:"day_of_month"`
}

type OneTimeOccurrence struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ChargingPolicy struct {
	MinCharge             int     `json:"min_charge"`
	MaxCharge             int     `json:"max_charge"`
	ChargeIfPriceBelow    float32 `json:"charge_if_price_below"`
	DischargeIfPriceAbove float32 `json:"discharge_if_price_above"`
}

func ParseDate(date string) (pgtype.Date, error) {
	if date == "" {
		return pgtype.Date{
			Valid: false,
		}, nil
	}

	dateVal, err := time.Parse("2006-01-02 15:04", date)
	if err != nil {
		return pgtype.Date{
			Valid: false,
		}, errors.New("Invalid format")
	}

	return pgtype.Date{
		Time:  dateVal,
		Valid: true,
	}, nil
}

func ParseInt(value int) pgtype.Int2 {
	if value == 0 {
		return pgtype.Int2{
			Valid: false,
		}
	}

	return pgtype.Int2{
		Int16: int16(value),
		Valid: true,
	}
}

func ParseFloat(value float32) (pgtype.Numeric, error) {
	var n pgtype.Numeric

	err := n.Scan(fmt.Sprintf("%.2f", value))
	if err != nil {
		return pgtype.Numeric{}, errors.New("Invalid format")
	}

	return n, nil
}

func ParseWeekDay(day string) (repository.NullWeekday, error) {
	if day == "" {
		return repository.NullWeekday{}, nil
	}

	var weekday repository.Weekday
	err := weekday.Scan(day)

	if err != nil || (weekday != repository.WeekdayMon && weekday != repository.WeekdayTue && weekday != repository.WeekdayWed && weekday != repository.WeekdayThu && weekday != repository.WeekdayFri && weekday != repository.WeekdaySat && weekday != repository.WeekdaySun) {
		return repository.NullWeekday{}, errors.New("Invalid format")
	}

	return repository.NullWeekday{
		Weekday: weekday,
		Valid:   true,
	}, nil
}

func ParseNthOfMonth(nth string) (repository.NullOrdinal, error) {
	if nth == "" {
		return repository.NullOrdinal{}, nil
	}

	var ordinal repository.Ordinal
	err := ordinal.Scan(nth)

	if err != nil || (ordinal != repository.Ordinal1st &&
		ordinal != repository.Ordinal2nd &&
		ordinal != repository.Ordinal3rd &&
		ordinal != repository.Ordinal4th &&
		ordinal != repository.OrdinalLast) {
		return repository.NullOrdinal{}, errors.New("Invalid format")
	}

	return repository.NullOrdinal{
		Ordinal: ordinal,
		Valid:   true,
	}, nil
}

func ToChargingPreferenceParams(preference *ChargingPreference) repository.CreateChargingPreferenceParams {
	return repository.CreateChargingPreferenceParams{
		ID:       preference.Id,
		UserID:   preference.UserId,
		Name:     preference.Name,
		Priority: preference.Priority,
		Enabled:  preference.Enabled,
	}
}

func ToRegularOccurrenceParams(occurrence *RegularOccurrence) (repository.CreateRegularOccurrenceParams, error) {
	params := repository.CreateRegularOccurrenceParams{}

	if occurrence.TimeOfDay.IsZero() {
		return params, errors.New("Time of day must be specified")
	}

	parsedUntil, err := ParseDate(occurrence.Until)
	if err != nil {
		return params, err
	}
	parsedRepeat := ParseInt(occurrence.Repeat)

	if parsedUntil.Valid && parsedRepeat.Valid {
		return params, errors.New("Either until or repeat must be specified, but not both")
	}

	parsedDayOfMonth := ParseInt(occurrence.DayOfMonth)
	parsedWeekday, err := ParseWeekDay(occurrence.DayOfWeek)
	if err != nil {
		return params, err
	}

	if parsedDayOfMonth.Valid && parsedWeekday.Valid {
		return params, errors.New("Either day_of_month or day_of_week must be specified, but not both")
	} else if !parsedDayOfMonth.Valid && !parsedWeekday.Valid {
		return params, errors.New("Either day_of_month or day_of_week must be specified")
	}

	parsedNthOfMonth, err := ParseNthOfMonth(occurrence.NthOfMonth)
	if err != nil {
		return params, err
	}
	if !parsedWeekday.Valid && parsedNthOfMonth.Valid {
		return params, errors.New("day_of_week must be specified if nth_of_month is specified")
	}

	return repository.CreateRegularOccurrenceParams{
		ID: uuid.New(),
		TimeOfDay: pgtype.Time{
			Microseconds: int64(occurrence.TimeOfDay.Hour())*int64(time.Hour/time.Microsecond) +
				int64(occurrence.TimeOfDay.Minute())*int64(time.Minute/time.Microsecond),
			Valid: true,
		},
		Repeat:     parsedRepeat,
		Until:      parsedUntil,
		DayOfWeek:  parsedWeekday,
		NthOfMonth: parsedNthOfMonth,
		DayOfMonth: parsedDayOfMonth,
	}, nil
}

func ToOneTimeOccurrenceParams(occurrence *OneTimeOccurrence) (repository.CreateOneTimeOccurrenceParams, error) {
	params := repository.CreateOneTimeOccurrenceParams{}
	params.ID = uuid.New()

	start, err := ParseDate(occurrence.Start)
	if err != nil {
		return params, errors.New("Invalid start date")
	}
	params.DateStart = start

	end, err := ParseDate(occurrence.End)
	if err != nil {
		return params, errors.New("Invalid end date")
	}
	params.DateEnd = end

	return params, nil

}
