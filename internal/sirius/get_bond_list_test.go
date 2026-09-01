package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestApiClient_GetBondList_Returns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
	{
        "pages": {
            "current": 1,
            "total": 2
        },
        "total": 26,
        "bonds": [
            {
                "id": 13,
                "caseReferenceNumber": "12345678",
                "clientFirstName": "Joseph",
                "clientLastName": "Smith",
                "companyName": "Company Ltd",
                "bondReferenceNumber": "BOND-1",
                "bondAmount": 101,
                "bondIssuedDate" : "2025-01-01T00:00:00+00:00",
                "client":{"id":63},
                "bondStatus":{"handle":"MATCH","label":"Match"},
                "deputyNames": ["Angela White", "Gary Black"]
            }
        ]
    }`

	params := BondListParams{
		Team:    model.Team{Id: 13},
		Page:    1,
		PerPage: 25,
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
				BondAmount:          101,
				BondIssuedDate:      model.NewDate("01/01/2025"),
				BondClient: model.Client{
					Id: 63,
				},
				BondStatus: model.RefData{
					Label:  "Match",
					Handle: "MATCH",
				},
				Deputies: []string{"Angela White", "Gary Black"},
			},
		},
		Pages: model.PageInformation{
			PageCurrent: 1,
			PageTotal:   2,
		},
		TotalBonds: 26,
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

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	bondList, err := client.GetBondList(getContext(nil), BondListParams{
		Team:    model.Team{Id: 13},
		Page:    1,
		PerPage: 25,
	})

	expectedResponse := BondList{
		Bonds:      nil,
		Pages:      model.PageInformation{},
		TotalBonds: 0,
	}

	assert.Equal(t, expectedResponse, bondList)

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/bonds/without-orders?limit=25&page=1",
		Method: http.MethodGet,
	}, err)
}

func TestGetBondList_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("User exists").
		UponReceiving("A request for bonds without orders").
		WithRequest("GET", "/supervision-api/v1/bonds/without-orders", func(b *consumer.V4RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
			b.Query("limit", matchers.S("25"))
			b.Query("page", matchers.S("1"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.MapMatcher{
				"pages": matchers.StructMatcher{
					"current": matchers.Like(1),
					"total":   matchers.Like(2),
				},
				"total": matchers.Like(26),
				"bonds": matchers.EachLike(matchers.StructMatcher{
					"id":                  matchers.Like(13),
					"caseReferenceNumber": matchers.Like("12345678"),
					"clientFirstName":     matchers.Like("Joseph"),
					"clientLastName":      matchers.Like("Smith"),
					"companyName":         matchers.Like("Company Ltd"),
					"bondReferenceNumber": matchers.Like("BOND-1"),
					"bondAmount":          matchers.Like(101),
					"bondIssuedDate":      matchers.Like("2025-01-01T00:00:00+00:00"),
					"client": matchers.StructMatcher{
						"id": matchers.Like(63),
					},
					"bondStatus": matchers.StructMatcher{
						"handle": matchers.Like("MATCH"),
						"label":  matchers.Like("Match"),
					},
					"deputyNames": matchers.EachLike("Angela White", 1),
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			bonds, err := client.GetBondList(getContext(nil), BondListParams{
				Team:    model.Team{Id: 1},
				Page:    1,
				PerPage: 25,
			})
			assert.NoError(t, err)

			assert.EqualValues(t, BondList{
				Bonds: []model.Bond{
					{
						Id:                  13,
						CourtRef:            "12345678",
						FirstName:           "Joseph",
						LastName:            "Smith",
						CompanyName:         "Company Ltd",
						BondReferenceNumber: "BOND-1",
						BondAmount:          101,
						BondIssuedDate:      model.NewDate("01/01/2025"),
						BondClient: model.Client{
							Id: 63,
						},
						BondStatus: model.RefData{
							Label:  "Match",
							Handle: "MATCH",
						},
						Deputies: []string{"Angela White", "Gary Black"},
					},
				},
				Pages: model.PageInformation{
					PageCurrent: 1,
					PageTotal:   2,
				},
				TotalBonds: 26,
			}, bonds)
			return nil
		})

	assert.NoError(t, err)
}
