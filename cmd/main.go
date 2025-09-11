package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/juanrojas09/gocourse_course/pkg/bootstrap"
	"github.com/juanrojas09/gocourse_course/pkg/handler"
)

func main() {

	log := bootstrap.InitLogger()
	//inicializo db
	db, err := bootstrap.InitDb()
	handleStartupErrors(err)

	//inicializamos instancias del dominio de curso
	endpoints := bootstrap.InitCourses(db, log)

	ctx := context.Background()
	//hacemos uso del handler
	h := handler.NewHttpServer(ctx, endpoints)

	addr := os.Getenv("API_URL_BASE") + ":" + os.Getenv("API_PORT")
	srv := http.Server{
		Handler: setupHeadersAndCors(h),
		Addr:    addr,
	}

	/* Creas un canal de errores: errChn := make(chan error)
	Inicias una goroutine que ejecuta el servidor y envía cualquier error al canal:
	errChn <- srv.ListenAndServe()
	El hilo principal se queda esperando en err = <-errChn hasta que la goroutine envíe un error (o el servidor termine).
	Cuando recibe el error, lo asigna a err y sigue la ejecución (por ejemplo, para validar si es nil y cortar si hay error).
	*/
	errChn := make(chan error)
	go func() {
		log.Println("listening in", addr)
		errChn <- srv.ListenAndServe()
	}()

	err = <-errChn
	handleStartupErrors(err)
}

func handleStartupErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// funcion que en cada req, setea los valores por defecto del cors y headers y recibe la peticion y escribe la respuesta con el serveHTTP
// es importante entender que usa el http.HandlerFunc pq de ahi saca la informacion de la peticion con los objetos w y r
// luego dentro de la funcion anonima ejecuta el seteo de los headers y demas para servirlo
// posterior a eso todo se destruye
func setupHeadersAndCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", "GET<POST,PATCH,OPTIONS,HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})

}
