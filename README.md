# RAiway

Database migration tool for RAI databases

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
