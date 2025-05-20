# schedulerd

A stateless, rule-based scheduling engine that assigns workloads to workers based on declarative filter and scoring rules.

Built in Go. Powered by [CEL (Common Expression Language)](https://opensource.google/projects/cel). Designed for production.

---

## What It Does

- Accepts a **workload** and a list of **candidate workers**
- Applies user-defined **filter rules** to eliminate invalid candidates
- Applies **scoring rules** to rank remaining candidates
- Returns a ranked list of the best-suited workers
- Stateless by design — no persistence, no side effects

---

## Example

### Input JSON (`input.json`)
```json
{
  "workload": {
    "labels": {
      "duration": 4
    }
  },
  "workers": [
    {
      "id": "alice",
      "labels": {
        "capacity": 8,
        "load": 2,
        "location": "remote"
      }
    },
    {
      "id": "bob",
      "labels": {
        "capacity": 4,
        "load": 3,
        "location": "onsite"
      }
    }
  ]
}
```

### Rules YAML (`rules.yaml`)
```yaml
filters:
  type: and
  rules:
    - type: expression
      expression: "worker.capacity - worker.load >= workload.duration"

scorers:
  type: weighted_avg
  rules:
    - weight: 0.7
      type: expression
      expression: "(worker.capacity - worker.load) / worker.capacity"
    - weight: 0.3
      type: expression
      expression: "worker.location == 'remote' ? 1.0 : 0.0"
```

---

## CLI Usage

```bash
schedulerd eval \
  --config rules.yaml \
  --input input.json
```

### Output
```json
[
  {
    "worker_id": "alice",
    "score": 0.91
  }
]
```

---

## REST API

```bash
curl -X POST http://localhost:8080/schedule \
  -H "Content-Type: application/json" \
  -d @input.json
```

Returns the same ranked list of worker candidates.

---

## 📄 Features (Planned or In Progress)

- [x] CEL-based expression evaluation
- [x] Composable rule logic (and/or/not)
- [x] Score aggregators (avg, max, weighted_avg, etc.)
- [ ] CLI support (`eval` mode)
- [ ] REST API (`/schedule`)
- [ ] Rule trace/debug output
- [ ] Docker image

---

## Philosophy

- **Stateless**: The engine doesn’t store state — inputs and outputs only.
- **Composable**: Filters and scorers are built from simple parts.
- **Extensible**: Add new rules, plugins, or output modes easily.
- **Safe**: Uses CEL for sandboxed evaluation, no panic risk.

---

## License

[MIT](./LICENSE)

---

## Credits

Built by [Your Name or Organization].  
Inspired by scheduling systems in Kubernetes, Prometheus, and classical rule engines.
