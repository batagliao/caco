package services

import (
	"context"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const sheet_id = "1dVe1zoi9_g0dE5POQ0g0L-lLH8RfQG02HvWtqfa_dlU"

const (
	start_date = 0
	end_date   = 1
	name       = 2
	slack_id   = 3
)

type SpreadsheetService struct {
	sheetservice *sheets.Service
}

type GetRowByDateResult struct {
	StartDate time.Time
	EndDate   time.Time
	Name      string
	SlackID   string
}

func NewSpredsheetService(ctx context.Context) *SpreadsheetService {
	sheetsService, errSheet := sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if errSheet != nil {
		log.Fatal(errSheet)
	}
	return &SpreadsheetService{
		sheetservice: sheetsService,
	}
}

func (s *SpreadsheetService) GetRowByDate(date time.Time) (*GetRowByDateResult, error) {
	sh, err := s.sheetservice.Spreadsheets.Get(sheet_id).IncludeGridData(true).Do()
	if err != nil {
		return nil, err
	}
	data := sh.Sheets[0].Data[0]

	for idx, r := range data.RowData {
		if idx == int(data.StartRow) {
			continue
		}

		startDateTxt := r.Values[start_date].FormattedValue
		if startDateTxt == "" {
			return nil, nil
		}
		startDate, err := time.Parse("02/01/2006", startDateTxt)
		if err != nil {
			return nil, err
		}

		endDateTxt := r.Values[end_date].FormattedValue
		if endDateTxt == "" {
			return nil, nil
		}
		endDate, err := time.Parse("02/01/2006", endDateTxt)
		if err != nil {
			return nil, err
		}

		// println(date.String())
		// println(startDate.Before(date))
		// println(startDate.Equal(date))
		// println(endDate.After(date))
		// println(endDate.Equal(date))

		// é o resgistro que queremos?
		if (startDate.Before(date) || startDate.Equal(date)) &&
			(endDate.After(date) || endDate.Equal(date)) {
			result := &GetRowByDateResult{
				StartDate: startDate,
				EndDate:   endDate,
				Name:      r.Values[name].FormattedValue,
				SlackID:   r.Values[slack_id].FormattedValue,
			}
			return result, nil
		}
	}

	return nil, nil
}
