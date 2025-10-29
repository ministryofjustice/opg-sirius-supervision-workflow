package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestApiClient_GetBondList_Returns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
	[
		{
			"id": 13,
			"caseReferenceNumber": "12345678",
			"clientFirstName": "Joseph",
			"clientLastName": "Smith",
			"companyName": "Company Ltd",
			"bondReferenceNumber": "BOND-1",
			"bondAmount": 1.1,
			"bondIssuedDate" : "2025-01-01T00:00:00+00:00"
		}
	]`

	params := BondListParams{
		Team: model.Team{Id: 13},
	}

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := BondList{
		Bonds: []model.Bond{
			{
				Id:                  13,
				CourtRef:            "12345678",
				FirstName:           "Joseph",
				LastName:            "Smith",
				CompanyName:         "Company Ltd",
				BondReferenceNumber: "BOND-1",
				BondAmount:          1.1,
				BondIssuedDate:      model.NewDate("01/01/2025"),
			},
		},
	}

	bondList, err := client.GetBondList(getContext(nil), params)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResponse, bondList)
}

func TestApiClient_GetBondList_Returns500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	bondList, err := client.GetBondList(getContext(nil), BondListParams{
		Team: model.Team{Id: 13},
	})

	expectedResponse := BondList{
		Bonds: nil,
	}

	assert.Equal(t, expectedResponse, bondList)

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/bonds-without-orders",
		Method: http.MethodGet,
	}, err)
}
