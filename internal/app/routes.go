package app

import "net/http"

const (
	AccessControlAllowOrigin = "Access-Control-Allow-Origin"
	ApplicationJson          = "application/json"
	ContentType              = "Content-Type"
)

func (s *HezzlWebService) ConfigureRouter() {
	s.get("/goods/list", s.goodList)
	s.post("/good/create", s.createGood)
	s.patch("/good/update", s.updateGood)
	s.delete("/good/remove", s.deleteGood)
	s.patch("/good/reprioritiize", s.reprioritiize)
}

func (s *HezzlWebService) post(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, patchPostMethod(http.MethodPost, handler))
}

func (s *HezzlWebService) patch(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, patchPostMethod(http.MethodPatch, handler))
}

func (s *HezzlWebService) get(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, deleteGetMethod(http.MethodGet, handler))
}

func (s *HezzlWebService) delete(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	s.router.HandleFunc(path, deleteGetMethod(http.MethodDelete, handler))
}

func patchPostMethod(method string, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method && r.Header.Get(ContentType) == ApplicationJson {
			w.Header().Set(AccessControlAllowOrigin, "*")
			w.Header().Set(ContentType, ApplicationJson)
			handler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func deleteGetMethod(method string, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			w.Header().Set(AccessControlAllowOrigin, "*")
			w.Header().Set(ContentType, ApplicationJson)
			handler(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}
