from pymongo import MongoClient

MONGO_URI = "mongodb://localhost:27017/"
DB_NAME = "TracerApp"
COLLECTIONS = ["events", "activities", "properties"]

client = MongoClient(MONGO_URI)
db = client[DB_NAME]

existing_collections = db.list_collection_names()

if set(COLLECTIONS).issubset(set(existing_collections)):
    print("Resetting existing collections...")
    for collection in COLLECTIONS:
        db[collection].delete_many({})
else:
    print("Creating new collections...")
    for collection in COLLECTIONS:
        db.create_collection(collection)

print("Database reset successfully.")
