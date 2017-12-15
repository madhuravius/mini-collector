package influxdb

import (
	"github.com/aptible/mini-collector/batch"
	client "github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
)

func buildBatchPoints(database string, entries []batch.Entry) client.BatchPoints {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})

	if err != nil {
		// These are just fatal errors because they're caused
		// by coding errors, not transient problems.
		log.Fatalf("could not build batch points: %+v", err)
	}

	for _, entry := range entries {
		fields := entryToFields(&entry)

		pt, err := client.NewPoint("metrics", entry.Tags, fields, entry.Time)
		if err != nil {
			// Same as above
			log.Fatalf("could not build point: %+v", err)
		}

		bp.AddPoint(pt)
	}

	return bp
}
