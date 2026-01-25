package controllers

import (
	"net/http"
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type GeoController struct {
	BaseController
	app     *app.App
	geoRepo *repositories.GeoRepository
}

func NewGeoController(app *app.App) *GeoController {
	return &GeoController{
		app:     app,
		geoRepo: repositories.NewGeoRepository(app),
	}
}

// ListCities 取得城市列表（含區域）
// @Summary 取得城市列表
// @Tags Geo
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.GeoCityResource}
// @Router /api/v1/geo/cities [get]
func (ctl *GeoController) ListCities(ctx *gin.Context) {
	cities, err := ctl.geoRepo.ListCities(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get cities",
		})
		return
	}

	var cityResources []resources.GeoCityResource
	for _, city := range cities {
		var districtResources []resources.GeoDistrictResource
		for _, district := range city.Districts {
			districtResources = append(districtResources, resources.GeoDistrictResource{
				ID:     district.ID,
				CityID: district.CityID,
				Name:   district.Name,
			})
		}
		cityResources = append(cityResources, resources.GeoCityResource{
			ID:        city.ID,
			Name:      city.Name,
			Districts: districtResources,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   cityResources,
	})
}
