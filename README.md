# AudioMeta Core

The AudioMeta adapters share this dependency-free Go 1.18 module. It defines the field registry, typed values, validation, ownership rules, and cumulative editing operations. The adapter supplies container I/O.

`Fields` returns canonical definitions. `Definition.Mapping` reports preferred format storage and read aliases. `Encode` and `Decode` define the fallback wire values. Text, dates, URIs and UUIDs are represented as text; integers, finite decimals and booleans use their decimal or literal representation. Binary values use unpadded URL-safe base64. Structured values use unpadded URL-safe base64 of JSON with the exported Go struct member names and JSON base64 byte data. The version in the surrounding namespace identifies this contract.

New editors stage cumulative edits relative to an open, unmodified input. Maps, slices and picture bytes returned from editor operations belong to the caller. A value is immutable and can be shared. Use `Value.Equal` to compare metadata, including artwork contents. Adapters and their readers require external synchronization.

`contracttest.Run` tests all registry definitions with the same input and expectations. Build the modules together using a local Go workspace while preparing module versions.
