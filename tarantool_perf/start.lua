-- listen to port 3302
box.cfg{listen = 3302}

-- create our space
-- it only IN MEMORY space
box.schema.space.create("xxx", { if_not_exists = true , temporary=true })

-- create index with string as key
box.space.xxx:create_index('id', {parts = {1, 'int'}, type='HASH'})

-- give access to 'guest' user without password, to all operations
-- you don't want this kind of access in production
box.schema.user.grant('guest', 'read,write,execute', 'universe')

-- just test
box.space.xxx:put{-1, "anu"}
