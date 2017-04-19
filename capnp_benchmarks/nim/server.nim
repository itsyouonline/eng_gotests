import reactor, capnp, caprpc, specification

type 
  MyStore = ref object of RootObj
    jwt*: string
  
  MyStoreFactory = ref object of RootObj

proc get(self: MyStore, id:ID): Future[Product] =
  let completer = newCompleter[Product]()
  completer.complete(Product(id:id))
  return completer.getFuture

proc set(self: MyStore, obj: Product): Future[ID] =
  let completer = newCompleter[ID]()
  completer.complete(obj.id)
  return completer.getFuture

proc toCapServer(self: MyStore): CapServer =
  return toGenericCapServer(self.asStore)


proc createStore(self: MyStoreFactory, jwt: JWT): Future[Store] =
  # validate jwt
  let completer = newCompleter[Store]()
  completer.complete(MyStore().asStore)
  return completer.getFuture

proc toCapServer(self: MyStoreFactory): CapServer =
  return toGenericCapServer(self.asStoreFactory)

proc main() {.async.} =
  let server = await createTcpServer(6000, "127.0.0.1")
  let myStoreFactory = new(MyStoreFactory)

  asyncFor conn in server.incomingConnections:
    discard newTwoPartyServer(conn.BytePipe, myStoreFactory.toCapServer())

when isMainModule:
  main().runMain
