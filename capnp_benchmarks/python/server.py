import capnp
import click
from jose import jwt
schema = capnp.load('../specification.capnp')


class BadAuthentification(Exception):
    pass


class ServerImpl(schema.StoreFactory.Server):

    def __init__(self, key):
        self.key = key

    def _validate_jwt(self, token):
        jwt.decode(token, self.key, algorithms=['ES384'], audience='capnp')
        return True

    def createStore(self, jwt, _context, **kwargs):
        if self._validate_jwt(jwt.payload):
            return StoreImpl()
        raise BadAuthentification()


class NotFound(Exception):
    pass


class StoreImpl(schema.Store.Server):

    def __init__(self):
        self._db = {}

    def _generate_id(self, obj):
        return b'id'

    def set(self, object, **kwargs):
        """
        set object in the store
        """
        id = self._generate_id(object)
        self._db[id] = object.to_dict()
        return schema.ID.new_message(data=id)

    def get(self, id, **kwargs):
        """
        get object from the store
        """
        if id.data not in self._db:
            raise NotFound()
        obj = self._db[id.data]
        return schema.Object.new_message(**obj)


@click.command()
@click.option('--bind', '-b', default='0.0.0.0:6000', help='Capnp rpc bind address')
@click.option('--profile', '-p', is_flag=True, default=False, help='Enable profiling')
def main(bind, profile):
    """
    Runs the server bound to the\
    given address/port ADDRESS may be '*' to bind to all local addresses.\
    :PORT may be omitted to choose a port automatically.
    """

    JWT_SIGNING_KEY = """-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2
7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6
6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv
-----END PUBLIC KEY-----"""

    server = capnp.TwoPartyServer(bind, bootstrap=ServerImpl(JWT_SIGNING_KEY))
    print("server start on {}".format(bind))
    server.run_forever()

if __name__ == '__main__':
    main()
