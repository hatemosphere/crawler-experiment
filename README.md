### Run instructions

```
$ git clone https://github.com/hatemosphere/crawler-experiment
$ cd crawler-experiment && docker build . -t crawler1337
$ docker run crawler1337 [domain]
```
Format for domain argument: `Scheme+Host` (eg. "https://example.com")

### TODOs

1. Add links between pages using graphs
