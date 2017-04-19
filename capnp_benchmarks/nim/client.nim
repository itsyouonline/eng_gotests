import strutils, times

import reactor
import capnp, caprpc, specification

proc benchmark(store: Store, n: int) {.async.} =
   
  var promises = newSeq[Future[Product]]()
  for i in 0..<n:
    let prod = Product(number:i.uint32, title: "Hello world", id:ID(data:i.intToStr))
    let prom = store.set(prod).then((id:ID) => store.get(id))
    promises.add(prom)

  for i in 0..<n:
    let prom = promises[i]
    discard await prom
  
proc main() {.async.} =
  # create store factory
  let sys = newTwoPartyClient(await connectTcp("127.0.0.1:6000"))
  let obj = await sys.bootstrap()
  
  let jwt = JWT(payload:"12345")
  let store = await obj.castAs(StoreFactory).createStore(jwt)
 
 
  echo("starting benchmark")
  for n in [1000, 10000, 50000]:
    let start = getTime()
    await benchmark(store, n)
    echo(n," get/seet took ", toSeconds(getTime()) - toSeconds(start), " seconds")

when isMainModule:
  main().runMain()
