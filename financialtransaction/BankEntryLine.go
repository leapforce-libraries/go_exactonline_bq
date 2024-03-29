package exactonline_bq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_bigquery "cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"

	errortools "github.com/leapforce-libraries/go_errortools"
	financialtransaction "github.com/leapforce-libraries/go_exactonline_new/financialtransaction"
	bigquery "github.com/leapforce-libraries/go_google/bigquery"
	types "github.com/leapforce-libraries/go_types"
)

type BankEntryLine struct {
	SoftwareClientLicenseGuid_ string
	Created_                   time.Time
	Modified_                  time.Time
	ID                         string
	Account                    string
	AccountCode                string
	AccountName                string
	AmountDC                   float64
	AmountFC                   float64
	AmountVATFC                float64
	Asset                      string
	AssetCode                  string
	AssetDescription           string
	CostCenter                 string
	CostCenterDescription      string
	CostUnit                   string
	CostUnitDescription        string
	Created                    _bigquery.NullTimestamp
	Creator                    string
	CreatorFullName            string
	Date                       _bigquery.NullTimestamp
	Description                string
	Division                   int32
	Document                   string
	DocumentNumber             int32
	DocumentSubject            string
	EntryID                    string
	EntryNumber                int32
	ExchangeRate               float64
	GLAccount                  string
	GLAccountCode              string
	GLAccountDescription       string
	LineNumber                 int32
	Modified                   _bigquery.NullTimestamp
	Modifier                   string
	ModifierFullName           string
	Notes                      string
	OffsetID                   string
	OurRef                     int32
	Project                    string
	ProjectCode                string
	ProjectDescription         string
	Quantity                   float64
	VATCode                    string
	VATCodeDescription         string
	VATPercentage              float64
	VATType                    string
}

func getBankEntryLine(c *financialtransaction.BankEntryLine, softwareClientLicenseGuid string) BankEntryLine {
	t := time.Now()

	return BankEntryLine{
		softwareClientLicenseGuid,
		t, t,
		c.ID.String(),
		c.Account.String(),
		c.AccountCode,
		c.AccountName,
		c.AmountDC,
		c.AmountFC,
		c.AmountVATFC,
		c.Asset.String(),
		c.AssetCode,
		c.AssetDescription,
		c.CostCenter,
		c.CostCenterDescription,
		c.CostUnit,
		c.CostUnitDescription,
		bigquery.DateToNullTimestamp(c.Created),
		c.Creator.String(),
		c.CreatorFullName,
		bigquery.DateToNullTimestamp(c.Date),
		c.Description,
		c.Division,
		c.Document.String(),
		c.DocumentNumber,
		c.DocumentSubject,
		c.EntryID.String(),
		c.EntryNumber,
		c.ExchangeRate,
		c.GLAccount.String(),
		c.GLAccountCode,
		c.GLAccountDescription,
		c.LineNumber,
		bigquery.DateToNullTimestamp(c.Modified),
		c.Modifier.String(),
		c.ModifierFullName,
		c.Notes,
		c.OffsetID.String(),
		c.OurRef,
		c.Project.String(),
		c.ProjectCode,
		c.ProjectDescription,
		c.Quantity,
		c.VATCode,
		c.VATCodeDescription,
		c.VATPercentage,
		c.VATType,
	}
}

func (service *Service) WriteBankEntryLines(bucketHandle *storage.BucketHandle, softwareClientLicenseGuid string, lastModified *time.Time) ([]*storage.ObjectHandle, int, interface{}, *errortools.Error) {
	if bucketHandle == nil {
		return nil, 0, nil, nil
	}

	objectHandles := []*storage.ObjectHandle{}
	var w *storage.Writer

	call := service.FinancialTransactionService().NewGetBankEntryLinesCall(lastModified)

	rowCount := 0
	batchRowCount := 0
	batchSize := 10000

	for true {
		bankEntryLines, e := call.Do()
		if e != nil {
			return nil, 0, nil, e
		}

		if bankEntryLines == nil {
			break
		}

		if batchRowCount == 0 {
			guid := types.NewGuid()
			objectHandle := bucketHandle.Object((&guid).String())
			objectHandles = append(objectHandles, objectHandle)

			w = objectHandle.NewWriter(context.Background())
		}

		for _, tl := range *bankEntryLines {
			batchRowCount++

			b, err := json.Marshal(getBankEntryLine(&tl, softwareClientLicenseGuid))
			if err != nil {
				return nil, 0, nil, errortools.ErrorMessage(err)
			}

			// Write data
			_, err = w.Write(b)
			if err != nil {
				return nil, 0, nil, errortools.ErrorMessage(err)
			}

			// Write NewLine
			_, err = fmt.Fprintf(w, "\n")
			if err != nil {
				return nil, 0, nil, errortools.ErrorMessage(err)
			}
		}

		if batchRowCount > batchSize {
			// Close and flush data
			err := w.Close()
			if err != nil {
				return nil, 0, nil, errortools.ErrorMessage(err)
			}
			w = nil

			fmt.Printf("#BankEntryLines flushed: %v\n", batchRowCount)

			rowCount += batchRowCount
			batchRowCount = 0
		}
	}

	if w != nil {
		// Close and flush data
		err := w.Close()
		if err != nil {
			return nil, 0, nil, errortools.ErrorMessage(err)
		}

		rowCount += batchRowCount
	}

	fmt.Printf("#BankEntryLines: %v\n", rowCount)

	return objectHandles, rowCount, BankEntryLine{}, nil
}
