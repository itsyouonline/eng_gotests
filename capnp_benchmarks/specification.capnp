@0xd5376380d33c2177;

struct JWT {
    payload @0 :Data;
}

struct ID {
    data @0 :Data;
}

struct Object {
    id @0 :ID;
    title @1 :Text;
    number @2 :UInt32;
}

interface StoreFactory {
    createStore @0 (jwt: JWT) -> (store :Store);
}

interface Store {
    get @0 (id :ID) -> (object :Object);
    set @1 (object :Object) -> (id :ID);
}
