db = db.getSiblingDB('post');
// create and move to your new database
db.createUser({
'user': "mongo",
'pwd': "mongo",
'roles': [{
    'role': 'dbOwner',
    'db': 'post'}]});
// user created
db.createCollection('post');
// add new collection