package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
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

		switch r.Method {
		case http.MethodGet:
			bondList, err := client.GetBondList(ctx, sirius.BondListParams{
				Team: app.SelectedTeam,
			})
			if err != nil {
				return err
			}

			vars.BondList = bondList

			return tmpl.Execute(w, vars)

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
