import capnp
import time
import click
import requests
schema = capnp.load('../specification.capnp')


def get_jwt(client_id, client_secret):
    url = "https://itsyou.online/v1/oauth/access_token?grant_type=client_credentials&client_id={client_id}&client_secret={client_secret}&response_type=id_token&aud=capnp"
    resp = requests.post(url.format(client_id=client_id, client_secret=client_secret))

    if resp.status_code != 200:
        return None

    token = resp.content
    return token


def validate(store):
    obj = schema.Object.new_message(number=42, title='Hello world')
    print("set object to the store")
    promise = store.set(obj)
    id = promise.wait().id

    print("get object from store")
    obj2 = store.get(id).wait().object

    if obj2.to_dict() == obj.to_dict():
        print("object are the same")
    else:
        print('error, object differ')


def benchark(store, n):
    start = time.time()
    promises = []
    for i in range(n):
        obj = schema.Object.new_message(number=i, title='Hello world')
        promises.append(store.set(obj))

    for promise in promises:
        id = promise.wait().id
        obj = store.get(id).wait().object

    end = time.time()
    return end - start


@click.command()
@click.option('--bind', '-b', default='localhost:6000', help='Capnp rpc bind address')
@click.option('--profile', '-p', is_flag=True, default=False, help='Enable profiling')
def main(bind, profile):
    CLIENT_ID = '9Hawa4Z_Ki9Z8f5gPgIqR4qspx40'
    CLIENT_SECRET = 'slKjOo1YYSDEN22ZrJU904_8EHdq'

    capnp_client = capnp.TwoPartyClient(bind)
    storeFactory = capnp_client.bootstrap().cast_as(schema.StoreFactory)

    print("get store")
    jwt = schema.JWT.new_message(payload=get_jwt(CLIENT_ID, CLIENT_SECRET))
    promise = storeFactory.createStore(jwt)
    store = promise.wait().store

    print("validation test:")
    validate(store)

    print("\nperformances: ")
    for n in [1000, 10000, 50000]:
        duration = benchark(store, n)
        print("set/get {} objects in {} sec".format(n, duration))

if __name__ == '__main__':
    main()
