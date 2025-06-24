package charging_preferences_tests

import (
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/V2G-Minor-Fontys/server/internal/charging_preferences"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestParseDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    pgtype.Date
		wantErr bool
	}{
		{"Happy path", args{"2025-12-20 14:09"}, pgtype.Date{Time: time.Date(2025, 12, 20, 14, 9, 0, 0, time.UTC), Valid: true}, false},
		{"Empty date", args{""}, pgtype.Date{Valid: false}, false},
		{"Invalid format", args{"2025-12-20"}, pgtype.Date{Valid: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charging_preferences.ParseDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	type args struct {
		value int
	}
	tests := []struct {
		name string
		args args
		want pgtype.Int2
	}{
		{"Happy path", args{1}, pgtype.Int2{Int16: int16(1), Valid: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := charging_preferences.ParseInt(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	type args struct {
		value float32
	}
	tests := []struct {
		name    string
		args    args
		want    pgtype.Numeric
		wantErr bool
	}{
		{"Happy path", args{3.14}, pgtype.Numeric{Int: big.NewInt(314), Exp: -2, Valid: true}, false},
		{"Negative number", args{-3.14}, pgtype.Numeric{Int: big.NewInt(-314), Exp: -2, Valid: true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charging_preferences.ParseFloat(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWeekDay(t *testing.T) {
	type args struct {
		day string
	}
	tests := []struct {
		name    string
		args    args
		want    repository.NullWeekday
		wantErr bool
	}{
		{"Happy path", args{"Mon"}, repository.NullWeekday{Weekday: repository.WeekdayMon, Valid: true}, false},
		{"Empty day", args{""}, repository.NullWeekday{Valid: false}, false},
		{"Invalid day", args{"-"}, repository.NullWeekday{Valid: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charging_preferences.ParseWeekDay(tt.args.day)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWeekDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseWeekDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNthOfMonth(t *testing.T) {
	type args struct {
		nth string
	}
	tests := []struct {
		name    string
		args    args
		want    repository.NullOrdinal
		wantErr bool
	}{
		{"Happy path", args{"1st"}, repository.NullOrdinal{Ordinal: repository.Ordinal1st, Valid: true}, false},
		{"Happy path - last", args{"last"}, repository.NullOrdinal{Ordinal: repository.OrdinalLast, Valid: true}, false},
		{"Empty nth", args{""}, repository.NullOrdinal{Valid: false}, false},
		{"Invalid nth", args{"-"}, repository.NullOrdinal{Valid: false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charging_preferences.ParseNthOfMonth(tt.args.nth)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNthOfMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNthOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToChargingPreferenceParams(t *testing.T) {
	type args struct {
		preference *charging_preferences.ChargingPreference
	}
	tests := []struct {
		name string
		args args
		want repository.CreateChargingPreferenceParams
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := charging_preferences.ToChargingPreferenceParams(tt.args.preference); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToChargingPreferenceParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToRegularOccurrenceParams(t *testing.T) {
	type args struct {
		occurrence *charging_preferences.RegularOccurrence
	}
	tests := []struct {
		name    string
		args    args
		want    repository.CreateRegularOccurrenceParams
		wantErr bool
	}{
		{"Happy path - until", args{&charging_preferences.RegularOccurrence{
			TimeOfDay:  time.Date(2025, 12, 20, 12, 30, 0, 0, time.UTC),
			Until:      "2025-12-20 14:09",
			DayOfWeek:  "Mon",
			NthOfMonth: "2nd",
		}}, repository.CreateRegularOccurrenceParams{
			TimeOfDay: pgtype.Time{Microseconds: int64(12)*int64(time.Hour/time.Microsecond) +
				int64(30)*int64(time.Minute/time.Microsecond), Valid: true},
			Until:      pgtype.Date{Time: time.Date(2025, 12, 20, 14, 9, 0, 0, time.UTC), Valid: true},
			DayOfWeek:  repository.NullWeekday{Weekday: repository.WeekdayMon, Valid: true},
			NthOfMonth: repository.NullOrdinal{Ordinal: repository.Ordinal2nd, Valid: true},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charging_preferences.ToRegularOccurrenceParams(tt.args.occurrence)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToRegularOccurrenceParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want.ID = got.ID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRegularOccurrenceParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToOneTimeOccurrenceParams(t *testing.T) {
	type args struct {
		occurrence *charging_preferences.OneTimeOccurrence
	}
	tests := []struct {
		name    string
		args    args
		want    repository.CreateOneTimeOccurrenceParams
		wantErr bool
	}{
		{"Happy path", args{&charging_preferences.OneTimeOccurrence{Start: "2025-12-20 03:30", End: "2025-12-24 13:30"}}, repository.CreateOneTimeOccurrenceParams{
			DateStart: pgtype.Date{Time: time.Date(2025, 12, 20, 3, 30, 0, 0, time.UTC), Valid: true},
			DateEnd:   pgtype.Date{Time: time.Date(2025, 12, 24, 13, 30, 0, 0, time.UTC), Valid: true},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := charging_preferences.ToOneTimeOccurrenceParams(tt.args.occurrence)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToOneTimeOccurrenceParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got.ID = tt.want.ID // Ensure IDs match for comparison
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToOneTimeOccurrenceParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
