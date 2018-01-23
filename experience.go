package main

import (
	"math"
	"net/http"
)

const (
	ROWS = 2000
)

// https://github.com/Jleagle/steamranks.com/blob/master/src/Application/Views/ExperienceView.php
func experienceHandler(w http.ResponseWriter, r *http.Request) {

	template := experienceTemplate{}

	template.Rows = append(template.Rows, experienceRow{
		Level: 0,
		Start: 0,
	})

	var xp float64
	xp = 10

	for i := 1; i <= ROWS+1; i++ {

		diff := (math.Ceil((float64(i) + 1) / 10)) * 100

		template.Rows = append(template.Rows, experienceRow{
			Level: i,
			Start: int(xp),
		})

		xp = xp + diff
	}

	for i := 1; i <= ROWS; i++ {

		nextRow := template.Rows[i+1]
		thisRow := template.Rows[i]

		template.Rows[i].Diff = nextRow.Start - thisRow.Start
		template.Rows[i].End = nextRow.Start - 1
	}

	template.Rows = template.Rows[0:ROWS]

	returnTemplate(w, "experience", template)
}

type experienceTemplate struct {
	Rows []experienceRow
}

type experienceRow struct {
	Level int
	Start int
	End   int
	Diff  int
	Count int
}
