package main

import (
	"docs.greenitglobe.com/despiegk/gotests/capnp_benchmarks/go/store"

	log "github.com/Sirupsen/logrus"
)

//storeFactory is a local implementation of StoreFactory
type storeFactory struct{}

func (sf storeFactory) CreateStore(call store.StoreFactory_createStore) (err error) {
	jwtParam, err := call.Params.Jwt()
	if err != nil {
		return
	}
	jwt, err := jwtParam.Payload()
	if err != nil {
		return
	}
	//TODO: validate jwt

	//Create a new locally implemented Store capability
	ss := store.Store_ServerToClient(storeServer{jwt: string(jwt)})
	return call.Results.SetStore(ss)
}

//storeServer is a local implementation of Store
type storeServer struct {
	jwt string
}

func (ss storeServer) Get(call store.Store_get) error {
	log.Debugln("Get called on store authorized to", ss.jwt, "with params", call.Params)
	return nil
}

func (ss storeServer) Set(call store.Store_set) error {
	log.Debugln("Set called on store authorized to", ss.jwt, "with params", call.Params)
	return nil
}
