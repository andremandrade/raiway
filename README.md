# RAiway

Database migration tool for RAI databases

## How does RAiway work

![RAiway-Database-Migrator](https://github.com/andremandrade/raiway/assets/6182479/9d110c6f-15eb-4a70-b563-067d94bb7916)

## Migration script

### YAML scheme descriptor

* `script`: Root tag
  * Value: array of `Operation`

* `Operation`
  * Properties:
    * `type`: `string`
      * Value:
        * `load-csv`
        * `load-models`
        * `delete-models`
        * `enable-ics`
        * `disable-ics`
        * `update`
