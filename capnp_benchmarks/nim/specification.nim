import capnp, capnp/gensupport, collections/iface

import reactor, caprpc, caprpc/rpcgensupport
# file: specification.capnp

type
  JWT* = ref object
    payload*: string

  ID* = ref object
    data*: string

  Product* = ref object
    id*: ID
    title*: string
    number*: uint32

  StoreFactory* = distinct Interface
  StoreFactory_CallWrapper* = ref object of CapServerWrapper

  StoreFactory_createStore_Params* = ref object
    jwt*: JWT

  StoreFactory_createStore_Result* = ref object
    store*: Store

  Store* = distinct Interface
  Store_CallWrapper* = ref object of CapServerWrapper

  Store_get_Params* = ref object
    id*: ID

  Store_get_Result* = ref object
    `object`*: Product

  Store_set_Params* = ref object
    obj*: Product

  Store_set_Result* = ref object
    id*: ID



makeStructCoders(JWT, [], [
  (payload, 0, PointerFlag.none, true)
  ], [])

makeStructCoders(ID, [], [
  (data, 0, PointerFlag.none, true)
  ], [])

makeStructCoders(Product, [
  (number, 0, 0, true)
  ], [
  (id, 0, PointerFlag.none, true),
  (title, 1, PointerFlag.text, true)
  ], [])

interfaceMethods StoreFactory:
  toCapServer(): CapServer
  createStore(jwt: JWT): Future[Store]

proc getIntefaceId*(t: typedesc[StoreFactory]): uint64 = return 15656406668745158309'u64

miscCapMethods(StoreFactory, StoreFactory_CallWrapper)

proc capCall*[T: StoreFactory](cap: T, id: uint64, args: AnyPointer): Future[AnyPointer] =
  case int(id):
    of 0:
      let argObj = args.castAs(StoreFactory_createStore_Params)
      let retVal = cap.createStore(argObj.jwt)
      return wrapFutureInSinglePointer(StoreFactory_createStore_Result, store, retVal)
    else: raise newException(NotImplementedError, "not implemented")

proc createStore*[T: StoreFactory_CallWrapper](self: T, jwt: JWT): Future[Store] =
  return getFutureField(self.cap.call(15656406668745158309'u64, 0, toAnyPointer(StoreFactory_createStore_Params(jwt: jwt))).castAs(StoreFactory_createStore_Result), store)

makeStructCoders(StoreFactory_createStore_Params, [], [
  (jwt, 0, PointerFlag.none, true)
  ], [])

makeStructCoders(StoreFactory_createStore_Result, [], [
  (store, 0, PointerFlag.none, true)
  ], [])

interfaceMethods Store:
  toCapServer(): CapServer
  get(id: ID): Future[Product]
  set(obj: Product): Future[ID]

proc getIntefaceId*(t: typedesc[Store]): uint64 = return 13540639907158926404'u64

miscCapMethods(Store, Store_CallWrapper)

proc capCall*[T: Store](cap: T, id: uint64, args: AnyPointer): Future[AnyPointer] =
  case int(id):
    of 0:
      let argObj = args.castAs(Store_get_Params)
      let retVal = cap.get(argObj.id)
      return wrapFutureInSinglePointer(Store_get_Result, `object`, retVal)
    of 1:
      let argObj = args.castAs(Store_set_Params)
      let retVal = cap.set(argObj.obj)
      return wrapFutureInSinglePointer(Store_set_Result, id, retVal)
    else: raise newException(NotImplementedError, "not implemented")

proc get*[T: Store_CallWrapper](self: T, id: ID): Future[Product] =
  return getFutureField(self.cap.call(13540639907158926404'u64, 0, toAnyPointer(Store_get_Params(id: id))).castAs(Store_get_Result), `object`)

proc set*[T: Store_CallWrapper](self: T, obj: Product): Future[ID] =
  return getFutureField(self.cap.call(13540639907158926404'u64, 1, toAnyPointer(Store_set_Params(obj: obj))).castAs(Store_set_Result), id)

makeStructCoders(Store_get_Params, [], [
  (id, 0, PointerFlag.none, true)
  ], [])

makeStructCoders(Store_get_Result, [], [
  (`object`, 0, PointerFlag.none, true)
  ], [])

makeStructCoders(Store_set_Params, [], [
  (obj, 0, PointerFlag.none, true)
  ], [])

makeStructCoders(Store_set_Result, [], [
  (id, 0, PointerFlag.none, true)
  ], [])


