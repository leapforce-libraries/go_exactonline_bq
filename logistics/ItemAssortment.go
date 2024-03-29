package exactonline_bq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/storage"

	errortools "github.com/leapforce-libraries/go_errortools"
	logistics "github.com/leapforce-libraries/go_exactonline_new/logistics"
	types "github.com/leapforce-libraries/go_types"
)

type ItemAssortment struct {
	SoftwareClientLicenseGuid_ string
	Created_                   time.Time
	Modified_                  time.Time
	ID                         string
	Code                       int32
	Description                string
	Division                   int32
	//Properties  []ItemAssortmentPropertyBQ
}

func getItemAssortment(c *logistics.ItemAssortment, softwareClientLicenseGuid string) ItemAssortment {
	t := time.Now()

	return ItemAssortment{
		softwareClientLicenseGuid,
		t, t,
		c.ID.String(),
		c.Code,
		c.Description,
		c.Division,
		//properties,
	}
}

func (service *Service) WriteItemAssortments(bucketHandle *storage.BucketHandle, softwareClientLicenseGuid string, lastModified *time.Time) ([]*storage.ObjectHandle, int, interface{}, *errortools.Error) {
	if bucketHandle == nil {
		return nil, 0, nil, nil
	}

	objectHandles := []*storage.ObjectHandle{}
	var w *storage.Writer

	params := logistics.GetItemAssortmentsCallParams{
		ModifiedAfter: lastModified,
	}

	call := service.LogisticsService().NewGetItemAssortmentsCall(&params)

	rowCount := 0
	batchRowCount := 0
	batchSize := 10000

	for true {
		itemAssortments, e := call.Do()
		if e != nil {
			return nil, 0, nil, e
		}

		if itemAssortments == nil {
			break
		}

		if batchRowCount == 0 {
			guid := types.NewGuid()
			objectHandle := bucketHandle.Object((&guid).String())
			objectHandles = append(objectHandles, objectHandle)

			w = objectHandle.NewWriter(context.Background())
		}

		for _, tl := range *itemAssortments {
			batchRowCount++

			b, err := json.Marshal(getItemAssortment(&tl, softwareClientLicenseGuid))
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

			fmt.Printf("#ItemAssortments flushed: %v\n", batchRowCount)

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

	fmt.Printf("#ItemAssortments: %v\n", rowCount)

	return objectHandles, rowCount, ItemAssortment{}, nil
}
