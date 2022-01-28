package api

import (
	_ "embed"
	"github.com/numary/ledger/pkg/api/controllers"
	"github.com/numary/ledger/pkg/api/middlewares"
	"github.com/numary/ledger/pkg/api/routes"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// API struct
type API struct {
	handler *gin.Engine
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

// NewAPI
func NewAPI(
	routes *routes.Routes,
) *API {
	gin.SetMode(gin.ReleaseMode)

	cc := cors.DefaultConfig()
	cc.AllowAllOrigins = true
	cc.AllowCredentials = true
	cc.AddAllowHeaders("authorization")

	h := &API{
		handler: routes.Engine(cc),
	}

	return h
}

type Config struct {
	StorageDriver string
	LedgerLister  controllers.LedgerLister
	HttpBasicAuth string
	Version       string
}

func Module(cfg Config) fx.Option {
	return fx.Options(
		controllers.ProvideVersion(func() string {
			return cfg.Version
		}),
		controllers.ProvideStorageDriver(func() string {
			return cfg.StorageDriver
		}),
		controllers.ProvideLedgerLister(func() controllers.LedgerLister {
			return cfg.LedgerLister
		}),
		middlewares.ProvideHTTPBasic(func() string {
			return cfg.HttpBasicAuth
		}),
		middlewares.Module,
		routes.Module,
		controllers.Module,
		fx.Provide(NewAPI),
	)
}