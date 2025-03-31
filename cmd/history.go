package cmd

import (
	"fmt"
	"time"

	"https://github.com/hagatasdelus/steph/db"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

const (
	maxStepLevel = 4
	weekdays     = 7
)

var stepColors = []tcell.Color{
	tcell.GetColor("#EFF2F5"),
	tcell.GetColor("#82E596"),
	tcell.GetColor("#26A148"),
	tcell.GetColor("#107F32"),
	tcell.GetColor("#014C19"),
}

var weekdayNames = [weekdays]string{
	"",
	"Mon",
	"",
	"Wed",
	"",
	"Fri",
	"",
}

var monthNames = []string{
	"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Display step history in a mapping format",
	Long:  `Display the history of your steps in a mapping format similar to GitHub contribution graph.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayStepHistory()
	},
}

func stepsToLevel(steps uint) int {
	if steps == 0 {
		return 0
	} else if steps < 3000 {
		return 1
	} else if steps < 6000 {
		return 2
	} else if steps < 9000 {
		return 3
	} else {
		return 4
	}
}

func formatDateKey(t time.Time) string {
	return t.Format("2006-01-02")
}

type dateBlockInfo struct {
	date    time.Time
	x, y    int
	width   int
	height  int
}

func displayStepHistory() {
	histories, err := db.GetAllStepHistories()
	if err != nil {
		fmt.Println("Failed to retrieve step histories:", err)
		return
	}

	if len(histories) == 0 {
		fmt.Println("No step history recorded yet.")
		return
	}

	stepsMap := make(map[string]int)
	var oldestDate, newestDate time.Time
	firstRun := true

	for _, h := range histories {
		date := time.Date(h.Date.Year(), h.Date.Month(), h.Date.Day(), 0, 0, 0, 0, time.Local)
		
		dateKey := formatDateKey(date)
		stepsMap[dateKey] = stepsToLevel(h.Steps)
		
		if firstRun {
			oldestDate = date
			newestDate = date
			firstRun = false
		} else {
			if date.Before(oldestDate) {
				oldestDate = date
			}
			if date.After(newestDate) {
				newestDate = date
			}
		}
	}

	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error creating new screen: %v\n", err)
		return
	}
	if err := s.Init(); err != nil {
		fmt.Printf("Error initializing screen: %v\n", err)
		return
	}
	defer s.Fini()

	s.EnableMouse()

	s.Clear()

	var startDate time.Time
	if !firstRun {
		startDate = time.Date(oldestDate.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	} else {
		currentYear := time.Now().Year()
		startDate = time.Date(currentYear, time.January, 1, 0, 0, 0, 0, time.Local)
	}

	minYear := oldestDate.Year()
	maxYear := newestDate.Year()

	var dateBlocks []dateBlockInfo
	
	currentStartDate := startDate
	selectedDateInfo := ""
	
	drawContributionHistory(s, currentStartDate, stepsMap, &dateBlocks, selectedDateInfo)
	s.Show()
	
	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
					close(quit)
					return
				} else if ev.Rune() == 'n' {
					nextYear := currentStartDate.Year() + 1
					if nextYear <= maxYear {
						currentStartDate = time.Date(nextYear, time.January, 1, 0, 0, 0, 0, time.Local)
						dateBlocks = nil
						drawContributionHistory(s, currentStartDate, stepsMap, &dateBlocks, selectedDateInfo)
						s.Show()
					}
				} else if ev.Rune() == 'p' {
					prevYear := currentStartDate.Year() - 1
					if prevYear >= minYear {
						currentStartDate = time.Date(prevYear, time.January, 1, 0, 0, 0, 0, time.Local)
						dateBlocks = nil
						drawContributionHistory(s, currentStartDate, stepsMap, &dateBlocks, selectedDateInfo)
						s.Show()
					}
				}
			case *tcell.EventMouse:
				if ev.Buttons() == tcell.ButtonPrimary {
					x, y := ev.Position()
					
					clickedDate := findDateFromPosition(x, y, dateBlocks)
					if !clickedDate.IsZero() {
						steps, err := db.GetStepsForDate(clickedDate)
						if err != nil {
							selectedDateInfo = fmt.Sprintf("Error retrieving data for %s: %v", clickedDate.Format("2006-01-02"), err)
						} else if steps > 0 {
							selectedDateInfo = fmt.Sprintf("Date: %s  Steps: %d", clickedDate.Format("2006-01-02"), steps)
						} else {
							selectedDateInfo = fmt.Sprintf("No steps recorded for %s", clickedDate.Format("2006-01-02"))
						}
						
						drawContributionHistory(s, currentStartDate, stepsMap, &dateBlocks, selectedDateInfo)
						s.Show()
					}
				}
			case *tcell.EventResize:
				s.Sync()
				dateBlocks = nil
				drawContributionHistory(s, currentStartDate, stepsMap, &dateBlocks, selectedDateInfo)
				s.Show()
			}
		}
	}()

	<-quit
}

func findDateFromPosition(x, y int, blocks []dateBlockInfo) time.Time {
	for _, block := range blocks {
		if x >= block.x && x < block.x+block.width &&
		   y >= block.y && y < block.y+block.height {
			return block.date
		}
	}
	return time.Time{}
}

func drawContributionHistory(s tcell.Screen, startDate time.Time, contributions map[string]int, 
                           dateBlocks *[]dateBlockInfo, selectedDateInfo string) {
	s.Clear()

	baseX, baseY := 5, 5
	cellWidth, cellHeight := 2, 1

	*dateBlocks = []dateBlockInfo{}

	for i := 0; i < weekdays; i++ {
		if weekdayNames[i] != "" {
			drawText(s, baseX-3, baseY+i*cellHeight, tcell.StyleDefault, weekdayNames[i])
		}
	}

	drawText(s, baseX, baseY-2, tcell.StyleDefault.Bold(true), fmt.Sprintf("Step History for %d", startDate.Year()))

	currentDate := startDate
	year := startDate.Year()
	endDate := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)

	x, y := 0, int(currentDate.Weekday())

	monthPositions := make(map[time.Month]int)

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		if currentDate.Day() == 1 {
			monthPositions[currentDate.Month()] = x
		}

		dateKey := formatDateKey(currentDate)
		
		level, exists := contributions[dateKey]
		if !exists {
			level = 0
		}

		drawCell(s, baseX+x*cellWidth, baseY+y*cellHeight, level)
		
		*dateBlocks = append(*dateBlocks, dateBlockInfo{
			date:   currentDate,
			x:      baseX + x*cellWidth,
			y:      baseY + y*cellHeight,
			width:  cellWidth,
			height: cellHeight,
		})

		currentDate = currentDate.AddDate(0, 0, 1)

		y = int(currentDate.Weekday())
		if y == 0 {
			x++
		}
	}

	for month, xPos := range monthPositions {
		monthName := monthNames[month-1]
		drawText(s, baseX+xPos*cellWidth, baseY-1, tcell.StyleDefault, monthName)
	}

	legendY := baseY + 8
	drawLegend(s, baseX, legendY, cellWidth)
	
	if selectedDateInfo != "" {
		drawText(s, baseX, legendY + 2, tcell.StyleDefault, selectedDateInfo)
	}
}

func drawLegend(s tcell.Screen, x, y, cellWidth int) {
	drawText(s, x, y, tcell.StyleDefault, "Less")
	
	for i := 0; i <= maxStepLevel; i++ {
		color := stepColors[i]
		style := tcell.StyleDefault.Background(color).Foreground(color)
		s.SetContent(x+6+i*2, y, ' ', nil, style)
		s.SetContent(x+6+i*2+1, y, ' ', nil, style)
	}
	
	drawText(s, x+8+(maxStepLevel+1)*2, y, tcell.StyleDefault, "More")
}

func drawCell(s tcell.Screen, x, y, level int) {
	color := stepColors[level]
	style := tcell.StyleDefault.Background(color).Foreground(color)
	s.SetContent(x, y, ' ', nil, style)
	s.SetContent(x+1, y, ' ', nil, style)
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
