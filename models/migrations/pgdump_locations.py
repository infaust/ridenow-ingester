import json

def escape(str):
	return str.replace("'", "''").replace("\\''", "''")


insert = ("INSERT INTO location VALUES ({id}, '{name}', '{description}', '{relevant}', point({lat}, {lon}), '{url}', '{country}');")
with open('msw_locations.json') as jsonfile:
    data = json.load(jsonfile)
    print("-- This dump has been auto-generated using pgdump_locations.py\n\n")
    for row in data:
        line = insert.format(id=row['_id'],
        	                 name=escape(row['name']),
        	                 description=escape(row['description']),
                             relevant='false',
        	                 lat=row['lat'],
        	                 lon=row['lon'],
        	                 url=escape(row['url']),
        	                 country=escape(row['country']['iso']))
        print(line)
