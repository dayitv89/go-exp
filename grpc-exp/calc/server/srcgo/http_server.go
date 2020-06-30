package srcgo

import "net/http"

//WebServerMux default grpc implementation; actually it doesn't require just implementing to show demo
func WebServerMux(grpcPort string) error {
	mux, cancel, err := setupMux(grpcPort)
	defer cancel()
	if err != nil {
		return err
	}
	return http.ListenAndServe(":3001", mux)
}
