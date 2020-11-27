package utils

import uuid "github.com/satori/go.uuid"

func GenUUID() string {
	v4, err := uuid.NewV4()
	if err != nil {
		// something has to be very wrong, if the uuid generation fails...
		panic(err)
	}

	return v4.String()
}
