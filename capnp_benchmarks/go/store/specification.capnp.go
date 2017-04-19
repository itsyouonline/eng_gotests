package store

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
	server "zombiezen.com/go/capnproto2/server"
)

type JWT struct{ capnp.Struct }

// JWT_TypeID is the unique identifier for the type JWT.
const JWT_TypeID = 0x8951bf8c0097aa3d

func NewJWT(s *capnp.Segment) (JWT, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return JWT{st}, err
}

func NewRootJWT(s *capnp.Segment) (JWT, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return JWT{st}, err
}

func ReadRootJWT(msg *capnp.Message) (JWT, error) {
	root, err := msg.RootPtr()
	return JWT{root.Struct()}, err
}

func (s JWT) String() string {
	str, _ := text.Marshal(0x8951bf8c0097aa3d, s.Struct)
	return str
}

func (s JWT) Payload() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s JWT) HasPayload() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JWT) SetPayload(v []byte) error {
	return s.Struct.SetData(0, v)
}

// JWT_List is a list of JWT.
type JWT_List struct{ capnp.List }

// NewJWT creates a new list of JWT.
func NewJWT_List(s *capnp.Segment, sz int32) (JWT_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return JWT_List{l}, err
}

func (s JWT_List) At(i int) JWT { return JWT{s.List.Struct(i)} }

func (s JWT_List) Set(i int, v JWT) error { return s.List.SetStruct(i, v.Struct) }

// JWT_Promise is a wrapper for a JWT promised by a client call.
type JWT_Promise struct{ *capnp.Pipeline }

func (p JWT_Promise) Struct() (JWT, error) {
	s, err := p.Pipeline.Struct()
	return JWT{s}, err
}

type ID struct{ capnp.Struct }

// ID_TypeID is the unique identifier for the type ID.
const ID_TypeID = 0x87b70cfb93ae82ba

func NewID(s *capnp.Segment) (ID, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return ID{st}, err
}

func NewRootID(s *capnp.Segment) (ID, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return ID{st}, err
}

func ReadRootID(msg *capnp.Message) (ID, error) {
	root, err := msg.RootPtr()
	return ID{root.Struct()}, err
}

func (s ID) String() string {
	str, _ := text.Marshal(0x87b70cfb93ae82ba, s.Struct)
	return str
}

func (s ID) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s ID) HasData() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s ID) SetData(v []byte) error {
	return s.Struct.SetData(0, v)
}

// ID_List is a list of ID.
type ID_List struct{ capnp.List }

// NewID creates a new list of ID.
func NewID_List(s *capnp.Segment, sz int32) (ID_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return ID_List{l}, err
}

func (s ID_List) At(i int) ID { return ID{s.List.Struct(i)} }

func (s ID_List) Set(i int, v ID) error { return s.List.SetStruct(i, v.Struct) }

// ID_Promise is a wrapper for a ID promised by a client call.
type ID_Promise struct{ *capnp.Pipeline }

func (p ID_Promise) Struct() (ID, error) {
	s, err := p.Pipeline.Struct()
	return ID{s}, err
}

type Object struct{ capnp.Struct }

// Object_TypeID is the unique identifier for the type Object.
const Object_TypeID = 0xf5f354a2de981d37

func NewObject(s *capnp.Segment) (Object, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Object{st}, err
}

func NewRootObject(s *capnp.Segment) (Object, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Object{st}, err
}

func ReadRootObject(msg *capnp.Message) (Object, error) {
	root, err := msg.RootPtr()
	return Object{root.Struct()}, err
}

func (s Object) String() string {
	str, _ := text.Marshal(0xf5f354a2de981d37, s.Struct)
	return str
}

func (s Object) Id() (ID, error) {
	p, err := s.Struct.Ptr(0)
	return ID{Struct: p.Struct()}, err
}

func (s Object) HasId() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Object) SetId(v ID) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewId sets the id field to a newly
// allocated ID struct, preferring placement in s's segment.
func (s Object) NewId() (ID, error) {
	ss, err := NewID(s.Struct.Segment())
	if err != nil {
		return ID{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Object) Title() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s Object) HasTitle() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Object) TitleBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s Object) SetTitle(v string) error {
	return s.Struct.SetText(1, v)
}

func (s Object) Number() uint32 {
	return s.Struct.Uint32(0)
}

func (s Object) SetNumber(v uint32) {
	s.Struct.SetUint32(0, v)
}

// Object_List is a list of Object.
type Object_List struct{ capnp.List }

// NewObject creates a new list of Object.
func NewObject_List(s *capnp.Segment, sz int32) (Object_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	return Object_List{l}, err
}

func (s Object_List) At(i int) Object { return Object{s.List.Struct(i)} }

func (s Object_List) Set(i int, v Object) error { return s.List.SetStruct(i, v.Struct) }

// Object_Promise is a wrapper for a Object promised by a client call.
type Object_Promise struct{ *capnp.Pipeline }

func (p Object_Promise) Struct() (Object, error) {
	s, err := p.Pipeline.Struct()
	return Object{s}, err
}

func (p Object_Promise) Id() ID_Promise {
	return ID_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type StoreFactory struct{ Client capnp.Client }

func (c StoreFactory) CreateStore(ctx context.Context, params func(StoreFactory_createStore_Params) error, opts ...capnp.CallOption) StoreFactory_createStore_Results_Promise {
	if c.Client == nil {
		return StoreFactory_createStore_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0xd946bae9fd9caea5,
			MethodID:      0,
			InterfaceName: "../specification.capnp:StoreFactory",
			MethodName:    "createStore",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(StoreFactory_createStore_Params{Struct: s}) }
	}
	return StoreFactory_createStore_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type StoreFactory_Server interface {
	CreateStore(StoreFactory_createStore) error
}

func StoreFactory_ServerToClient(s StoreFactory_Server) StoreFactory {
	c, _ := s.(server.Closer)
	return StoreFactory{Client: server.New(StoreFactory_Methods(nil, s), c)}
}

func StoreFactory_Methods(methods []server.Method, s StoreFactory_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xd946bae9fd9caea5,
			MethodID:      0,
			InterfaceName: "../specification.capnp:StoreFactory",
			MethodName:    "createStore",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := StoreFactory_createStore{c, opts, StoreFactory_createStore_Params{Struct: p}, StoreFactory_createStore_Results{Struct: r}}
			return s.CreateStore(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// StoreFactory_createStore holds the arguments for a server call to StoreFactory.createStore.
type StoreFactory_createStore struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  StoreFactory_createStore_Params
	Results StoreFactory_createStore_Results
}

type StoreFactory_createStore_Params struct{ capnp.Struct }

// StoreFactory_createStore_Params_TypeID is the unique identifier for the type StoreFactory_createStore_Params.
const StoreFactory_createStore_Params_TypeID = 0x8b90863d4e2beb85

func NewStoreFactory_createStore_Params(s *capnp.Segment) (StoreFactory_createStore_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return StoreFactory_createStore_Params{st}, err
}

func NewRootStoreFactory_createStore_Params(s *capnp.Segment) (StoreFactory_createStore_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return StoreFactory_createStore_Params{st}, err
}

func ReadRootStoreFactory_createStore_Params(msg *capnp.Message) (StoreFactory_createStore_Params, error) {
	root, err := msg.RootPtr()
	return StoreFactory_createStore_Params{root.Struct()}, err
}

func (s StoreFactory_createStore_Params) String() string {
	str, _ := text.Marshal(0x8b90863d4e2beb85, s.Struct)
	return str
}

func (s StoreFactory_createStore_Params) Jwt() (JWT, error) {
	p, err := s.Struct.Ptr(0)
	return JWT{Struct: p.Struct()}, err
}

func (s StoreFactory_createStore_Params) HasJwt() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s StoreFactory_createStore_Params) SetJwt(v JWT) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewJwt sets the jwt field to a newly
// allocated JWT struct, preferring placement in s's segment.
func (s StoreFactory_createStore_Params) NewJwt() (JWT, error) {
	ss, err := NewJWT(s.Struct.Segment())
	if err != nil {
		return JWT{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// StoreFactory_createStore_Params_List is a list of StoreFactory_createStore_Params.
type StoreFactory_createStore_Params_List struct{ capnp.List }

// NewStoreFactory_createStore_Params creates a new list of StoreFactory_createStore_Params.
func NewStoreFactory_createStore_Params_List(s *capnp.Segment, sz int32) (StoreFactory_createStore_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return StoreFactory_createStore_Params_List{l}, err
}

func (s StoreFactory_createStore_Params_List) At(i int) StoreFactory_createStore_Params {
	return StoreFactory_createStore_Params{s.List.Struct(i)}
}

func (s StoreFactory_createStore_Params_List) Set(i int, v StoreFactory_createStore_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// StoreFactory_createStore_Params_Promise is a wrapper for a StoreFactory_createStore_Params promised by a client call.
type StoreFactory_createStore_Params_Promise struct{ *capnp.Pipeline }

func (p StoreFactory_createStore_Params_Promise) Struct() (StoreFactory_createStore_Params, error) {
	s, err := p.Pipeline.Struct()
	return StoreFactory_createStore_Params{s}, err
}

func (p StoreFactory_createStore_Params_Promise) Jwt() JWT_Promise {
	return JWT_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type StoreFactory_createStore_Results struct{ capnp.Struct }

// StoreFactory_createStore_Results_TypeID is the unique identifier for the type StoreFactory_createStore_Results.
const StoreFactory_createStore_Results_TypeID = 0xd6513a272ff74767

func NewStoreFactory_createStore_Results(s *capnp.Segment) (StoreFactory_createStore_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return StoreFactory_createStore_Results{st}, err
}

func NewRootStoreFactory_createStore_Results(s *capnp.Segment) (StoreFactory_createStore_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return StoreFactory_createStore_Results{st}, err
}

func ReadRootStoreFactory_createStore_Results(msg *capnp.Message) (StoreFactory_createStore_Results, error) {
	root, err := msg.RootPtr()
	return StoreFactory_createStore_Results{root.Struct()}, err
}

func (s StoreFactory_createStore_Results) String() string {
	str, _ := text.Marshal(0xd6513a272ff74767, s.Struct)
	return str
}

func (s StoreFactory_createStore_Results) Store() Store {
	p, _ := s.Struct.Ptr(0)
	return Store{Client: p.Interface().Client()}
}

func (s StoreFactory_createStore_Results) HasStore() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s StoreFactory_createStore_Results) SetStore(v Store) error {
	if v.Client == nil {
		return s.Struct.SetPtr(0, capnp.Ptr{})
	}
	seg := s.Segment()
	in := capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	return s.Struct.SetPtr(0, in.ToPtr())
}

// StoreFactory_createStore_Results_List is a list of StoreFactory_createStore_Results.
type StoreFactory_createStore_Results_List struct{ capnp.List }

// NewStoreFactory_createStore_Results creates a new list of StoreFactory_createStore_Results.
func NewStoreFactory_createStore_Results_List(s *capnp.Segment, sz int32) (StoreFactory_createStore_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return StoreFactory_createStore_Results_List{l}, err
}

func (s StoreFactory_createStore_Results_List) At(i int) StoreFactory_createStore_Results {
	return StoreFactory_createStore_Results{s.List.Struct(i)}
}

func (s StoreFactory_createStore_Results_List) Set(i int, v StoreFactory_createStore_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// StoreFactory_createStore_Results_Promise is a wrapper for a StoreFactory_createStore_Results promised by a client call.
type StoreFactory_createStore_Results_Promise struct{ *capnp.Pipeline }

func (p StoreFactory_createStore_Results_Promise) Struct() (StoreFactory_createStore_Results, error) {
	s, err := p.Pipeline.Struct()
	return StoreFactory_createStore_Results{s}, err
}

func (p StoreFactory_createStore_Results_Promise) Store() Store {
	return Store{Client: p.Pipeline.GetPipeline(0).Client()}
}

type Store struct{ Client capnp.Client }

func (c Store) Get(ctx context.Context, params func(Store_get_Params) error, opts ...capnp.CallOption) Store_get_Results_Promise {
	if c.Client == nil {
		return Store_get_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0xbbea044111afa444,
			MethodID:      0,
			InterfaceName: "../specification.capnp:Store",
			MethodName:    "get",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Store_get_Params{Struct: s}) }
	}
	return Store_get_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}
func (c Store) Set(ctx context.Context, params func(Store_set_Params) error, opts ...capnp.CallOption) Store_set_Results_Promise {
	if c.Client == nil {
		return Store_set_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0xbbea044111afa444,
			MethodID:      1,
			InterfaceName: "../specification.capnp:Store",
			MethodName:    "set",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Store_set_Params{Struct: s}) }
	}
	return Store_set_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Store_Server interface {
	Get(Store_get) error

	Set(Store_set) error
}

func Store_ServerToClient(s Store_Server) Store {
	c, _ := s.(server.Closer)
	return Store{Client: server.New(Store_Methods(nil, s), c)}
}

func Store_Methods(methods []server.Method, s Store_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 2)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xbbea044111afa444,
			MethodID:      0,
			InterfaceName: "../specification.capnp:Store",
			MethodName:    "get",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Store_get{c, opts, Store_get_Params{Struct: p}, Store_get_Results{Struct: r}}
			return s.Get(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xbbea044111afa444,
			MethodID:      1,
			InterfaceName: "../specification.capnp:Store",
			MethodName:    "set",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Store_set{c, opts, Store_set_Params{Struct: p}, Store_set_Results{Struct: r}}
			return s.Set(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// Store_get holds the arguments for a server call to Store.get.
type Store_get struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Store_get_Params
	Results Store_get_Results
}

// Store_set holds the arguments for a server call to Store.set.
type Store_set struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Store_set_Params
	Results Store_set_Results
}

type Store_get_Params struct{ capnp.Struct }

// Store_get_Params_TypeID is the unique identifier for the type Store_get_Params.
const Store_get_Params_TypeID = 0x95ab49cf334ca91b

func NewStore_get_Params(s *capnp.Segment) (Store_get_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_get_Params{st}, err
}

func NewRootStore_get_Params(s *capnp.Segment) (Store_get_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_get_Params{st}, err
}

func ReadRootStore_get_Params(msg *capnp.Message) (Store_get_Params, error) {
	root, err := msg.RootPtr()
	return Store_get_Params{root.Struct()}, err
}

func (s Store_get_Params) String() string {
	str, _ := text.Marshal(0x95ab49cf334ca91b, s.Struct)
	return str
}

func (s Store_get_Params) Id() (ID, error) {
	p, err := s.Struct.Ptr(0)
	return ID{Struct: p.Struct()}, err
}

func (s Store_get_Params) HasId() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Store_get_Params) SetId(v ID) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewId sets the id field to a newly
// allocated ID struct, preferring placement in s's segment.
func (s Store_get_Params) NewId() (ID, error) {
	ss, err := NewID(s.Struct.Segment())
	if err != nil {
		return ID{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Store_get_Params_List is a list of Store_get_Params.
type Store_get_Params_List struct{ capnp.List }

// NewStore_get_Params creates a new list of Store_get_Params.
func NewStore_get_Params_List(s *capnp.Segment, sz int32) (Store_get_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Store_get_Params_List{l}, err
}

func (s Store_get_Params_List) At(i int) Store_get_Params { return Store_get_Params{s.List.Struct(i)} }

func (s Store_get_Params_List) Set(i int, v Store_get_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Store_get_Params_Promise is a wrapper for a Store_get_Params promised by a client call.
type Store_get_Params_Promise struct{ *capnp.Pipeline }

func (p Store_get_Params_Promise) Struct() (Store_get_Params, error) {
	s, err := p.Pipeline.Struct()
	return Store_get_Params{s}, err
}

func (p Store_get_Params_Promise) Id() ID_Promise {
	return ID_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Store_get_Results struct{ capnp.Struct }

// Store_get_Results_TypeID is the unique identifier for the type Store_get_Results.
const Store_get_Results_TypeID = 0x8766b59a2655c2ac

func NewStore_get_Results(s *capnp.Segment) (Store_get_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_get_Results{st}, err
}

func NewRootStore_get_Results(s *capnp.Segment) (Store_get_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_get_Results{st}, err
}

func ReadRootStore_get_Results(msg *capnp.Message) (Store_get_Results, error) {
	root, err := msg.RootPtr()
	return Store_get_Results{root.Struct()}, err
}

func (s Store_get_Results) String() string {
	str, _ := text.Marshal(0x8766b59a2655c2ac, s.Struct)
	return str
}

func (s Store_get_Results) Object() (Object, error) {
	p, err := s.Struct.Ptr(0)
	return Object{Struct: p.Struct()}, err
}

func (s Store_get_Results) HasObject() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Store_get_Results) SetObject(v Object) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewObject sets the object field to a newly
// allocated Object struct, preferring placement in s's segment.
func (s Store_get_Results) NewObject() (Object, error) {
	ss, err := NewObject(s.Struct.Segment())
	if err != nil {
		return Object{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Store_get_Results_List is a list of Store_get_Results.
type Store_get_Results_List struct{ capnp.List }

// NewStore_get_Results creates a new list of Store_get_Results.
func NewStore_get_Results_List(s *capnp.Segment, sz int32) (Store_get_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Store_get_Results_List{l}, err
}

func (s Store_get_Results_List) At(i int) Store_get_Results {
	return Store_get_Results{s.List.Struct(i)}
}

func (s Store_get_Results_List) Set(i int, v Store_get_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Store_get_Results_Promise is a wrapper for a Store_get_Results promised by a client call.
type Store_get_Results_Promise struct{ *capnp.Pipeline }

func (p Store_get_Results_Promise) Struct() (Store_get_Results, error) {
	s, err := p.Pipeline.Struct()
	return Store_get_Results{s}, err
}

func (p Store_get_Results_Promise) Object() Object_Promise {
	return Object_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Store_set_Params struct{ capnp.Struct }

// Store_set_Params_TypeID is the unique identifier for the type Store_set_Params.
const Store_set_Params_TypeID = 0xb936c81790391389

func NewStore_set_Params(s *capnp.Segment) (Store_set_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_set_Params{st}, err
}

func NewRootStore_set_Params(s *capnp.Segment) (Store_set_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_set_Params{st}, err
}

func ReadRootStore_set_Params(msg *capnp.Message) (Store_set_Params, error) {
	root, err := msg.RootPtr()
	return Store_set_Params{root.Struct()}, err
}

func (s Store_set_Params) String() string {
	str, _ := text.Marshal(0xb936c81790391389, s.Struct)
	return str
}

func (s Store_set_Params) Object() (Object, error) {
	p, err := s.Struct.Ptr(0)
	return Object{Struct: p.Struct()}, err
}

func (s Store_set_Params) HasObject() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Store_set_Params) SetObject(v Object) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewObject sets the object field to a newly
// allocated Object struct, preferring placement in s's segment.
func (s Store_set_Params) NewObject() (Object, error) {
	ss, err := NewObject(s.Struct.Segment())
	if err != nil {
		return Object{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Store_set_Params_List is a list of Store_set_Params.
type Store_set_Params_List struct{ capnp.List }

// NewStore_set_Params creates a new list of Store_set_Params.
func NewStore_set_Params_List(s *capnp.Segment, sz int32) (Store_set_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Store_set_Params_List{l}, err
}

func (s Store_set_Params_List) At(i int) Store_set_Params { return Store_set_Params{s.List.Struct(i)} }

func (s Store_set_Params_List) Set(i int, v Store_set_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Store_set_Params_Promise is a wrapper for a Store_set_Params promised by a client call.
type Store_set_Params_Promise struct{ *capnp.Pipeline }

func (p Store_set_Params_Promise) Struct() (Store_set_Params, error) {
	s, err := p.Pipeline.Struct()
	return Store_set_Params{s}, err
}

func (p Store_set_Params_Promise) Object() Object_Promise {
	return Object_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Store_set_Results struct{ capnp.Struct }

// Store_set_Results_TypeID is the unique identifier for the type Store_set_Results.
const Store_set_Results_TypeID = 0xe938761ec42676ae

func NewStore_set_Results(s *capnp.Segment) (Store_set_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_set_Results{st}, err
}

func NewRootStore_set_Results(s *capnp.Segment) (Store_set_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Store_set_Results{st}, err
}

func ReadRootStore_set_Results(msg *capnp.Message) (Store_set_Results, error) {
	root, err := msg.RootPtr()
	return Store_set_Results{root.Struct()}, err
}

func (s Store_set_Results) String() string {
	str, _ := text.Marshal(0xe938761ec42676ae, s.Struct)
	return str
}

func (s Store_set_Results) Id() (ID, error) {
	p, err := s.Struct.Ptr(0)
	return ID{Struct: p.Struct()}, err
}

func (s Store_set_Results) HasId() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Store_set_Results) SetId(v ID) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewId sets the id field to a newly
// allocated ID struct, preferring placement in s's segment.
func (s Store_set_Results) NewId() (ID, error) {
	ss, err := NewID(s.Struct.Segment())
	if err != nil {
		return ID{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Store_set_Results_List is a list of Store_set_Results.
type Store_set_Results_List struct{ capnp.List }

// NewStore_set_Results creates a new list of Store_set_Results.
func NewStore_set_Results_List(s *capnp.Segment, sz int32) (Store_set_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Store_set_Results_List{l}, err
}

func (s Store_set_Results_List) At(i int) Store_set_Results {
	return Store_set_Results{s.List.Struct(i)}
}

func (s Store_set_Results_List) Set(i int, v Store_set_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Store_set_Results_Promise is a wrapper for a Store_set_Results promised by a client call.
type Store_set_Results_Promise struct{ *capnp.Pipeline }

func (p Store_set_Results_Promise) Struct() (Store_set_Results, error) {
	s, err := p.Pipeline.Struct()
	return Store_set_Results{s}, err
}

func (p Store_set_Results_Promise) Id() ID_Promise {
	return ID_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

const schema_d5376380d33c2177 = "x\xda\xa4\x94OH\x14o\x18\xc7\x9f\xe7}g\x9d\x15" +
	"vq_\xc6\x1f\xcb\xaf0Q\xcc\xa8h\xcd\xa4,Q" +
	"V\xc54\xc5j'\x8d\xe88\x8e\xa3\xac\xa8\xb3\xcc\x8c" +
	"\xca\x1e\xa2?`\"u\x10\"\x8a\xe8\x14u\xb0B\xe9" +
	" \xa5\x1d\x02\xe9P\xa7\xb0(0\xa8\xbbv+D$" +
	"b\xe2\x9dev\xc6m7\x84.\xcb0\xfb\xe5\xf3<" +
	"\xcf\xf7y\xbes\xf8\x056\x93\xda@M\x00@n\x0c" +
	"\x14\xd9O\x97\xcfW\xdf[\x18\x98\x02V\x86\x00\x01\x14" +
	"\x01\xea.\x91.\x04\x94\xa6I\x1c\xd0^\xba6w\xeb" +
	"g\xe8\xf9\x14\xb0(\xda\x13\x15\x8d\x1f\xae\xa8\xf5\x1f3" +
	"Bi\x96\xbc\x97\x16\x08\x7fz\xc6\xb5\x9bM\x8f\xef\xdc" +
	"|%O\xe7Q\xae\x90\xcf\xd2\x17G\xb9\xeaP'\xbf" +
	"\x1d<\xd3t}\xe6\x06\xb0\xaal\xd9-\xb2\xc4\xcb\x16" +
	"S.\xd8=\xdb]\xf7\xae\xf3\xc9m\x7f_\xfbi+" +
	"\x17\xd4:\x82i\xe9\xc4L\xf4\xcd\xb1E\xbf@\xce\x08" +
	".:\x82\xb6\x87\xf3\xacEX\x7f\x09,J\xbdv\x00" +
	"\xa54]\x97&)\xef\xe5*\x9d\x92V\xf8\x93=\xd8" +
	"\xb1Y\xb3\xafA\xfe\xe4\xefg\x91.s\xda[\x87\xf6" +
	"h\xee\xfe\xaf\xb5\xa5\xf6\xd5\\Z\xdd\x1a\xadDi\xcb" +
	"\xc1m\xd0\x0e\xa9B\xe0\xb8\xb9\xf1\xea\xd7{\xc6\x8f\xaf" +
	"\xf9\x9b+\x16\x1cW\xff\x138\xae\xbe\xec\xee\xd7\x07\xbd" +
	"?6@\x8e\xa2\xdf,\xc7\xa2\xa3\xc2w\xa9\x85s\xa4" +
	"&a\x1e\x0e\xd9\xb1X\x8d\x99\xd2\xd4\xa40\x90T\x15" +
	"+\xa9\x8f\xc6T%5\x9aj\xe8\xb1tC\x8b\x0dj" +
	"V\xd59\xcd\x1c\x1b\xb6L\x00Y\xa0\x02\x80\x80\x00," +
	"\xdc\x00 \x07)\xca\xa5\x04\xe3z\xdf\x90\xa6Z\x18\xf1" +
	"\x0a\x03b\x040\xcb&9\xecNlK \xfaq\x07" +
	"<\\I\xbfb)\x18\x06\x82\xe1\xbf \xba\xe8\x85\xde" +
	"\x1cF\xab\xc7\xb8\x9cR\xd2\xc3\xba\xd2\xff\x07&\x90o" +
	"\xcavE\xb5t#\x1dS\x0dM\xb14\xe7UU<" +
	"\xa1\x18\xca\x88\xe9\xe7Wz|qh\xc2\xc2\x88{\x92" +
	"9\xd3\x16v2\xc3\xdc\xe6\xe3.\x0fJ\x93\xfd\x18\xf1" +
	"\"\xb1#\xaa\x99\x9f\xfaO\xdb\xe9)\xe1hnn\x90" +
	"\x06\x00\xb2iA7\xce\xac\xb6\x12\x08\xdb+\xa2\x17\x14" +
	"t\x8f\x92\xfd\xcf\xff\x0b\x8b\xe2\xa0f5\xa3h\xf2\xdf" +
	"\x04z\xc5\x8av\xba\x00\xf7\xea\xfcc\x1d\xf1\xc6*7" +
	"\xb9\x0a\x99\x97D@d\xbe\xa9h\xa1B\xa2n\xa43" +
	"\x97\xc3\x87s\xbf\x15\xe8\x86\x94\xb1> \xacX\xb4\xdd" +
	"f@\xd4\x0dm\xfb\x10\x85w\x917+;\xdeq\xee" +
	"&\xce\x96;\xbb\xe3\xdd\x86\xb2\xb8\x93\x1c\xd7LQ\xee" +
	"&\xc8\x10K\x91\xbf\xec\xe4\xd6\xb4Q\x94\x13\x04\x91\x94" +
	"\"\x01`\xa7\xf9\x15\x9c\xa2(\xf7\x16\xa8[n%\xad" +
	"a\x0dC@0\x04\x18\x1f\x1d\x1b\xe9\xd3\x0c\x0c\x02\xc1" +
	" \xe0\xef\x00\x00\x00\xff\xffB\xcd\x89/"

func init() {
	schemas.Register(schema_d5376380d33c2177,
		0x8766b59a2655c2ac,
		0x87b70cfb93ae82ba,
		0x8951bf8c0097aa3d,
		0x8b90863d4e2beb85,
		0x95ab49cf334ca91b,
		0xb936c81790391389,
		0xbbea044111afa444,
		0xd6513a272ff74767,
		0xd946bae9fd9caea5,
		0xe938761ec42676ae,
		0xf5f354a2de981d37)
}
