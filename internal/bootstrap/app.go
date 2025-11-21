package bootstrap

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heaveless/dbz-api/internal/application/character"
	"github.com/heaveless/dbz-api/internal/delivery/http"
	"github.com/heaveless/dbz-api/internal/delivery/http/handler"
	"github.com/heaveless/dbz-api/internal/infrastructure/api"
	"github.com/heaveless/dbz-api/internal/infrastructure/breaker"
	"github.com/heaveless/dbz-api/internal/infrastructure/repositoy"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Env *Env
	Db  *mongo.Client
	Svr *gin.Engine
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Db = NewDatabase(app.Env)

	collection := app.Db.Database(app.Env.DBName).Collection("characters")
	dbCollection := breaker.NewMongoDbCollection(collection)

	dbBreaker := breaker.NewDbCollectionWithBreaker(dbCollection, 3*time.Second)
	httpBreaker := breaker.NewHttpWithBreaker(3 * time.Second)

	characterRepo := repositoy.NewCharacterRepository(dbBreaker)
	characterApi := api.NewCharacterApi(app.Env.ApiUri, httpBreaker)

	characterService := character.NewCharacterService(characterRepo, characterApi)

	characterHandler := handler.NewCharacterHandler(characterService)

	app.Svr = http.NewServer(characterHandler)

	return *app
}

func (app *Application) CloseDbConnection() {
	CloseDatabaseConnection(app.Db)
}
