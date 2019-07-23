package dashserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/shiningrush/dolphin"
	"github.com/shiningrush/dolphin/dist/dolphinui"
)

// TaskStatus is return to web dashboard
type TaskStatus struct {
	TaskName         string `json:"taskName"`
	TaskType         string `json:"taskType"`
	PlanTime         string `json:"planTime"`
	TaskState        string `json:"taskState"`
	ResetBeforeBegin string `json:"resetBeforeBegin"`
	LastExecuteTime  string `json:"lastExecuteTime"`
	LastExecuteState string `json:"lastExecuteState"`
	LastExecuteCost  string `json:"lastExecuteCost"`
}

// Start server
func Start(dashServerAddr string) *http.Server {
	mux := http.NewServeMux()
	serveDolphinUIFile(mux)
	mux.Handle("/GetAllTasks", handlerWrapper(serveAllTasks))
	mux.Handle("/StartTask", handlerWrapper(startTask))
	mux.Handle("/StopTask", handlerWrapper(stopTask))
	mux.Handle("/ExecuteTask", handlerWrapper(executeTask))
	mux.Handle("/ResetTask", handlerWrapper(resetTask))

	log.Println("dash server listen at :" + dashServerAddr)

	hs := &http.Server{Addr: dashServerAddr, Handler: mux}
	go func() {
		if err := hs.ListenAndServe(); err != nil {
			log.Println("We got a error when init dashserver, error:" + err.Error())
		}
	}()

	return hs
}

func serveDolphinUIFile(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    dolphinui.Asset,
		AssetDir: dolphinui.AssetDir,
		Prefix:   "thirdparty/dashboard",
	})
	prefix := "/dashboard/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func serveAllTasks(w http.ResponseWriter, r *http.Request) error {
	var result []*TaskStatus
	allTasks := dolphin.GetAllTasks()
	for k, v := range allTasks {
		newTaskStatus := &TaskStatus{
			TaskName:         k,
			TaskType:         v.Type.ToString(),
			PlanTime:         "(Cron format) " + v.PlanTime,
			TaskState:        v.State.ToString(),
			ResetBeforeBegin: strconv.FormatBool(v.ResetBeforeBegin),
			LastExecuteTime:  v.LastExecuteTime.Local().Format("2006-01-02 15:04:05"),
			LastExecuteState: v.LastExecuteState,
			LastExecuteCost:  strconv.Itoa(v.LastExecuteCost) + "s",
		}
		if newTaskStatus.LastExecuteState == "" {
			newTaskStatus.LastExecuteState = "Everything is fine."
		}

		result = append(result, newTaskStatus)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TaskName < result[j].TaskName
	})

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err
	}

	w.Write(jsonResult)
	return nil
}

func startTask(w http.ResponseWriter, r *http.Request) error {
	if v, ok := r.URL.Query()["taskname"]; ok {
		if err := dolphin.Start(v[0]); err != nil {
			return err
		}
	} else {
		return errors.New("No taskname ")
	}

	return nil
}

func stopTask(w http.ResponseWriter, r *http.Request) error {
	if v, ok := r.URL.Query()["taskname"]; ok {
		if err := dolphin.Stop(v[0]); err != nil {
			return err
		}
	} else {
		return errors.New("No taskname ")
	}

	return nil
}
func executeTask(w http.ResponseWriter, r *http.Request) error {
	if v, ok := r.URL.Query()["taskname"]; ok {
		if err := dolphin.Execute(v[0]); err != nil {
			return err
		}
	} else {
		return errors.New("No taskname ")
	}

	return nil
}

func resetTask(w http.ResponseWriter, r *http.Request) error {
	if v, ok := r.URL.Query()["taskname"]; ok {
		if err := dolphin.Reset(v[0]); err != nil {
			return err
		}
	} else {
		return errors.New("No taskname ")
	}

	return nil
}

type handlerWrapper func(w http.ResponseWriter, r *http.Request) error

func (h handlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		http.Error(w, "{\"msg\":\" "+err.Error()+"\"}", 500)
	}
}
