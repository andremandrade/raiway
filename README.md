# RAiway

Database migration tool for RAI databases

## How to use

* Build this project: `./build`
* Create your migration script and add it to the folder that you prefer
  * Migration files should follow the standard `{migration_number}.yaml`, example:
    * `1.yaml`
    * `2.yaml`
  * The `migration_number` is used to define the migrations sequence
* Run `./raiway -p [PROFILE] -d [DATABASE] -e [ENGINE] -m [MIGRATIONS_FOLDER_PATH]`

## How does RAiway work

![RAiway-Database-Migrator](https://github.com/andremandrade/raiway/assets/6182479/9d110c6f-15eb-4a70-b563-067d94bb7916)

## Migration script

A migration script is a YAML file with an array of `Operation` in the root.

The properties of an `Operation` depends on the `type` property that defines the REL operation to be executed. The supported operations are:

* `load-csv`: Load a CSV file from a local file to a REL model
* `load-models`: Load local files as REL nodels
* `delete-models`: Delete existing models
* `update`: Execute queries in `read-only=false` mode  from a local file or an inline query
* `disable-ics`: Disable integrity constraints
* `enable-ics`: Enable integrity constraints

### Example of a migration script

```yaml
- type: load-csv
  hostMode: local
  filePath: example/data/people.csv
  modelName: people_csv
  delimiter: '|'
  quotechar: '~'
  escapechar: '\\'
  scheme:
    ID: int

- type: load-csv
  hostMode: local
  filePath: example/data/company.csv
  modelName: company_csv
  scheme:
    company_id: int

- name: Entity models # optional - this information goes to log
  type: load-models
  files: 
  - example/rel/model/person.rel
  - example/rel/model/company.rel
  prefix: example

- name: Deleting models to move to entities folder
  type: delete-models
  models:
  - example/person
  - example/company

- name: Re-creating entities models into entity folder
  type: load-models
  files: 
  - example/rel/model/person.rel
  - example/rel/model/company.rel
  prefix: example/entities
  
- type: enable-ics

- type: disable-ics

- type: update
  query: | # inline query
    def insert:jobs = loaded_jobs

- type: update
  filePath: example/rel/update/create-system-user.rel
```
