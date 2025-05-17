package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func GuidToPgUUID(id uuid.UUID) pgtype.UUID {
	var uuidArr [16]byte
	copy(uuidArr[:], id[:])
	return pgtype.UUID{
		Bytes: uuidArr,
		Valid: true,
	}
}

func ParsePgUUIDToGuid(pgUUID pgtype.UUID) (uuid.UUID, error) {
	if !pgUUID.Valid {
		return uuid.Nil, fmt.Errorf("invalid UUID (null)")
	}
	return uuid.FromBytes(pgUUID.Bytes[:])
}
