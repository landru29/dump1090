package main

import (
	"log/slog"

	"github.com/landru29/dump1090/internal/serialize"
	"github.com/landru29/dump1090/internal/serialize/basestation"
	"github.com/landru29/dump1090/internal/serialize/json"
	"github.com/landru29/dump1090/internal/serialize/nmea"
	"github.com/landru29/dump1090/internal/serialize/none"
	"github.com/landru29/dump1090/internal/serialize/text"
)

func provideSerializers(
	log *slog.Logger,
	nmeaVessel nmea.VesselType,
	nmeaMid uint16,
) (map[string]serialize.Serializer, []serialize.Serializer) {
	serializers := map[string]serialize.Serializer{}

	availableSerializers := []serialize.Serializer{
		none.Serializer{},
		json.Serializer{},
		text.Serializer{},
		basestation.Serializer{},
		nmea.New(nmeaVessel, nmeaMid),
	}

	for _, serializer := range availableSerializers {
		log.Info("loading serializer", "name", serializer.String())

		serializers[serializer.String()] = serializer
	}

	return serializers, availableSerializers
}
