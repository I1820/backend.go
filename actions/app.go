package actions

import (
	"context"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/unrolled/secure"
	validator "gopkg.in/go-playground/validator.v9"
	resty "gopkg.in/resty.v1"

	"github.com/I1820/types"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var db *mgo.Database
var validate *validator.Validate

// HTTP clients
var pmclient = resty.New().SetHostURL(envy.Get("PM_URL", "http://127.0.0.1:1375")).SetError(types.Error{})
var dmclient = resty.New().SetHostURL(envy.Get("DM_URL", "http://127.0.0.1:1373")).SetError(types.Error{})
var wfclient = resty.New().SetHostURL(envy.Get("WF_URL", "http://127.0.0.1:6976")).SetError(types.Error{})

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_backend_session",
		})
		// Automatically redirect to SSL
		app.Use(forcessl.Middleware(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		// If no content type is sent by the client
		// the application/json will be set, otherwise the client's
		// content type will be used.
		app.Use(contenttype.Add("application/json"))

		// Create mongodb connection
		url := envy.Get("DB_URL", "mongodb://172.18.0.1:27017")
		client, err := mgo.NewClient(url)
		if err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB new client error: %s", err)
		}
		if err := client.Connect(context.Background()); err != nil {
			buffalo.NewLogger("fatal").Fatalf("DB connection error: %s", err)
		}
		db = client.Database("i1820")

		// validator
		validate = validator.New()

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Slash issue
		app.Muxer().StrictSlash(true)
		// Routes
		app.GET("/about", AboutHandler)

		// swagger ui
		app.ServeFiles("/swagger", http.Dir("swagger"))

		api := app.Group("/api/v1")
		{
			// auth routes contains login, logout and signup
			auth := api.Group("/auth")
			{
				ar := AuthResource{}
				auth.POST("/login", ar.Login)
				auth.POST("/register", ar.Signup)
				auth.POST("/refresh", ar.Refresh)
			}

			// health routes
			health := api.Group("/health")
			{
				health.GET("/pm", PMHealthHandler)
				health.GET("/wf", WFHealthHandler)
				health.GET("/dm", DMHealthHandler)
			}

			// user additional info
			ar := AdditionalsResource{}
			ainfo := api.Group("/info")
			ainfo.Use(AuthMiddleware)
			{
				ainfo.GET("/{key}", ar.Show)
				ainfo.POST("/{key}", ar.Create)
			}

			// proxies to pm
			api.Resource("projects", ProjectsResource{}).Use(AuthMiddleware)
			pg := api.Group("projects/{project_id}")
			pg.Use(AuthMiddleware)
			{
				tr := ThingsResource{}
				pg.Resource("things", tr)
				pg.ANY("things/{thing_id}/tokens/{path:.*}", tr.Tokens)
				pg.ANY("things/{thing_id}/assets/{path:.*}", tr.Assets)
				pg.ANY("things/{thing_id}/connectivities/{path:.*}", tr.Connectivities)

				pg.ANY("things/{thing_id}/queries/{path:.*}", QueryHandler)
			}

			// proxies to wf
			api.POST("wf/{service}", WFHandler)
		}
	}

	return app
}
