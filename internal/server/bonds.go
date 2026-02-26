package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
)

type BondsClient interface {
	GetBondList(sirius.Context, sirius.BondListParams) (sirius.BondList, error)
}

type BondsPage struct {
	BondList sirius.BondList
	ListPage
}

func (bp BondsPage) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "bonds",
		SelectedTeam:    bp.App.SelectedTeam.Selector,
		SelectedPerPage: bp.PerPage,
		SelectedSort:    bp.Sort,
	}
}

func bonds(client BondsClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsAllocationsTeam() {
			page := ClientTasksPage{ListPage: ListPage{PerPage: 25}}
			return Redirect{Path: page.CreateUrlBuilder().GetTeamUrl(app.SelectedTeam)}
		}

		vars := BondsPage{}
		vars.App = app
		vars.UrlBuilder = vars.CreateUrlBuilder()

		params := r.URL.Query()
		perPageOptions := []int{25, 50, 100}
		vars.PerPage = paginate.GetRequestedElementsPerPage(params.Get("per-page"), perPageOptions)

		bondList, err := client.GetBondList(ctx, sirius.BondListParams{
			Team: app.SelectedTeam,
		})
		if err != nil {
			return err
		}

		vars.BondList = bondList

		vars.Pagination = paginate.Pagination{
			CurrentPage:     bondList.Pages.PageCurrent,
			TotalPages:      bondList.Pages.PageTotal,
			TotalElements:   bondList.TotalBonds,
			ElementsPerPage: vars.PerPage,
			ElementName:     "bonds",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}

		return tmpl.Execute(w, vars)
	}
}
