package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/alvarowolfx/golang-voice-iot/mearm"
)

type HttpRobotArm struct {
	robotArm *mearm.RobotArm
	server   *http.Server
	port     string
}

type HttpResponse struct {
	Message string                 `json:"message"`
	State   map[string]interface{} `json:"state"`
}

func NewHttpRobotArm(robotArm *mearm.RobotArm, port string) *HttpRobotArm {
	return &HttpRobotArm{
		robotArm: robotArm,
		port:     port,
	}
}

func (hra *HttpRobotArm) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/arm/{servo:(?:elbow|grip|base|shoulder)}/{value:[0-9]+}", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			json.NewEncoder(res).Encode(HttpResponse{
				Message: "error",
				State: map[string]interface{}{
					"message": "Invalid method. Only accepting PUT requests",
				},
			})
		}

		vars := mux.Vars(req)
		servo := vars["servo"]
		value := vars["value"]

		intValue, err := strconv.Atoi(value)

		if err != nil {
			json.NewEncoder(res).Encode(HttpResponse{
				Message: "error",
				State: map[string]interface{}{
					"message": err.Error(),
				},
			})
		}

		if intValue < 0 || intValue > 180 {
			json.NewEncoder(res).Encode(HttpResponse{
				Message: "error",
				State: map[string]interface{}{
					"message": "Invalid value. Must be between 0 and 180 degree",
				},
			})
		}

		switch servo {
		case "elbow":
			hra.robotArm.SetElbow(intValue)
			break
		case "shoulder":
			hra.robotArm.SetShoulder(intValue)
			break
		case "base":
			hra.robotArm.SetBase(intValue)
			break
		case "grip":
			if intValue > 0 {
				hra.robotArm.OpenGrip()
			} else {
				hra.robotArm.CloseGrip()
			}
			break
		}

		pos := hra.robotArm.RobotArmPos()

		json.NewEncoder(res).Encode(HttpResponse{
			Message: "ok",
			State: map[string]interface{}{
				"pos": pos,
			},
		})
	})

	hra.server = &http.Server{Addr: ":" + hra.port, Handler: r}
	go hra.server.ListenAndServe() // Http server blocks execution
}

func (hra HttpRobotArm) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	hra.server.Shutdown(ctx)
}
