# Python Test for capnp

This folder contains two files. the server and the client.
You need to edit the `client.py` file to set the CLIENT_ID and CLIENT_SECRET used to generate JWT.

Then, start the server:
```bash
python server.py *:6000
```
Then, run the client:
```
python client.py localhost:6000
```
