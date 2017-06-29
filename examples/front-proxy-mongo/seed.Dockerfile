FROM mongo

COPY apis/example.json /init.json
CMD mongoimport --host janus-database --db janus --collection api_specs --type json --file /init.json --jsonArray
