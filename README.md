# Shipbox

## Use Case 1 - CSV Composition

### foo.shpbx

```csv
{{with sheet}}
a,b,{{ .Table "bar" }}
c,{{ .Table "baz" }},d
{{ .Table "buzz" }},e,f
{{end}}
```

### bar.shpbx

```csv
bara,barb
barc,bard
```

### baz.shpbx

```csv
baza,bazb,bazc
bazd,baze,bazf 
bazg,bazh,bazi 
```

### buzz.shpbx

```csv
buzza
```

### output.json

```json
{
  "foo": [
    [
      "a", "b", {
        "bar": [
          ["bara", "barb"],
          ["barc", "bard"]
        ]
      }
    ],
    [
      "c", {
        "baz": [
          "baza", "bazb", "bazc",
          "bazd", "baze", "bazf",
          "bazg", "bazh", "bazi"
        ]
      }, "d"
    ],
    [
      {
        "buzz": [
          ["buzza"]
        ]
      }, "e", "f" 
    ]
  ]
}
```

### Command

```
shipbox -i foo.shpbx -i bar.shpbx -i baz.shpbx -i buzz.shpbx -o output.json
```

