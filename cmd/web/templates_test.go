package main

import (
	"gostats/cmd/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {

	var tesVal = map[string]string{
		"2004-mar-1":  "2004-03-01",
		"2004-lis-12": "2004-11-12",
		"1991-gru-31": "1991-12-31",
	}

	for k, v := range tesVal {

		t.Run("humanDate", func(t *testing.T) {

			m, _ := time.Parse("2006-01-02", v)
			hd := humanDate(m)

			assert.Equal(t, k, hd)
		})

	}

}
