package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
)

type BondMetaData struct {
	BondMetaData []model.AssigneeAndCount `json:"ecmCount"`
}

type BondList struct {
	Bonds      []model.Bond          `json:"bonds"`
	Pages      model.PageInformation `json:"pages"`
	TotalBonds int                   `json:"total"`
}

type BondListParams struct {
	Team    model.Team
	Page    int
	PerPage int
}

type bondClientResponse struct {
	ID int `json:"id"`
}

type bondResponse struct {
	ID                  int                 `json:"id"`
	CourtRef            string              `json:"caseReferenceNumber"`
	FirstName           string              `json:"clientFirstName"`
	LastName            string              `json:"clientLastName"`
	CompanyName         string              `json:"companyName"`
	BondReferenceNumber string              `json:"bondReferenceNumber"`
	BondAmount          int                 `json:"bondAmount"`
	BondIssuedDate      string              `json:"bondIssuedDate"`
	BondClient          *bondClientResponse `json:"client,omitempty"`
	BondStatus          *refDataResponse    `json:"bondStatus,omitempty"`
	Deputies            []string            `json:"deputyNames"`
}

func (r bondResponse) model() (model.Bond, error) {
	bondIssuedDate, err := parseModelDate(r.BondIssuedDate)
	if err != nil {
		return model.Bond{}, err
	}

	bondClient := model.Client{}
	if r.BondClient != nil {
		bondClient.Id = r.BondClient.ID
	}

	return model.Bond{
		Id:                  r.ID,
		CourtRef:            r.CourtRef,
		FirstName:           r.FirstName,
		LastName:            r.LastName,
		CompanyName:         r.CompanyName,
		BondReferenceNumber: r.BondReferenceNumber,
		BondAmount:          r.BondAmount,
		BondIssuedDate:      bondIssuedDate,
		BondClient:          bondClient,
		BondStatus:          r.BondStatus.model(),
		Deputies:            r.Deputies,
	}, nil
}

type bondListResponse struct {
	Bonds      []bondResponse          `json:"bonds"`
	Pages      pageInformationResponse `json:"pages"`
	TotalBonds int                     `json:"total"`
}

func (r bondListResponse) model() (BondList, error) {
	var bonds []model.Bond
	if r.Bonds != nil {
		bonds = make([]model.Bond, 0, len(r.Bonds))
		for _, bond := range r.Bonds {
			mappedBond, err := bond.model()
			if err != nil {
				return BondList{}, err
			}
			bonds = append(bonds, mappedBond)
		}
	}

	return BondList{
		Bonds:      bonds,
		Pages:      r.Pages.model(),
		TotalBonds: r.TotalBonds,
	}, nil
}

func (c *ApiClient) GetBondList(ctx Context, params BondListParams) (BondList, error) {
	var v BondList

	endpoint := fmt.Sprintf("/v1/bonds/without-orders?limit=%d&page=%d", params.PerPage, params.Page)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)

	if err != nil {
		c.logErrorRequest(req, err)
		return v, err
	}

	resp, err := c.http.Do(req)

	if err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}
	defer unchecked(resp.Body.Close)

	if resp.StatusCode == http.StatusUnauthorized {
		c.logResponse(req, resp, err)
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		c.logResponse(req, resp, err)
		return v, newStatusError(resp)
	}

	var response bondListResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logResponse(req, resp, err)
		return v, err
	}

	v, err = response.model()
	if err != nil {
		c.logResponse(req, resp, err)
		return BondList{}, err
	}

	return v, nil
}
