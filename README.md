# SVDB – Serverless VectorDB

> **A Serverless‑first, free vector database, simple to deploy, simple to manage.**

[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/danielealbano/svdb.svg)](https://github.com/danielealbano/svdb/issues)

---

SVDB ( **S**erverless **V**ector **D**ata**b**ase ) is an open‑source vector database designed from the ground up for
**serverless environments** - yet simple enough to run on your laptop:

* **Cloud‑native, serverless‑first** – deploy to Google Cloud Run & Functions, store shards on GCS.
* **Zero‑ops local dev** – identical Docker containers for local and cloud deployments.
* **Elastic shards & workers** – automatic shard rollover and worker provisioning.
* **Multiple collections** – isolate datasets while sharing the same infrastructure.
* **Backed by [USearch](https://github.com/unum-cloud/usearch/)** – Smaller & Faster Single-File
    Similarity Search & Clustering Engine for Vectors.
* **Simple CRUD & similarity search** – `insert`, `search`, `exists`, `get`, `delete`.
* **gRPC & HTTP APIs** – gRPC for high throughput, HTTP for easy integration.
* **Google BigQuery integration** – use SVDB via Google Cloud Functions in BigQuery.
* **BSD 3‑Clause licensed** – permissive for commercial and academic use.

> **Status ⚠️** SVDB is in early development – APIs and behaviour may change without notice.

---

## Table of contents

1. [How it works](#how-it-works)
2. [CLI](#cli)
3. [Collections, shards & workers](#collections-shards--workers)
4. [Road‑map](#road-map)
5. [Contributing](#contributing)
6. [License](#license)

---

## How it works

```
+-------------+     gRPC / HTTP     +-----------------------------+
| client / BQ |  ─────────────────▶ | engine frontend (Cloud Run) |
+-------------+                     +-----------------------------+
                                                   │
                                                   │
                                                   ▼
                                      +-------------------------+
                                      | worker (Cloud Function) |
                                      +-------------------------+
                                                 │   ▲ 
                                                 ▼   │
                                     Google Cloud Storage (S3 API)
```

1. The `engine‑frontend` receives query and write requests, maintaining cluster metadata (collections, shard state,
   worker addresses).
2. It routes the request to one or more **workers** responsible for the target shards.
3. The `engine-worker` load the shard files in memory from **GCS (S3‑compatible)**, execute the vector operation via,
   **USearch** and stream partial results back.
4. `engine‑frontend` merges/limits the results and returns them to the caller.
5. When a shard reaches its max size the worker marks it *sealed* and notifies `engine‑frontend` over gRPC; the frontend
   spins up a **new worker + shard** and re‑routes subsequent writes.

---

## CLI

The `svdb` CLI wraps the gRPC API and common admin operations (deployment commands are coming soon and will hide all Google Cloud details).

| command       | description                        |
|---------------|------------------------------------|
| `collections` | list / create / drop collections   |
| `insert`      | upsert a vector                    |
| `search`      | similarity search                  |
| `get` / `has` | retrieve or test presence of a key |
| `delete`      | remove a key                       |
| `len`         | number of items in a collection    |
| `status`      | show cluster status                |

Run `svdb --help` or `svdb <command> --help` for full syntax.

---

## Collections, shards & workers

* **Collection** – logical namespace / table for vectors with the same dimensionality.
* **Shard** – immutable binary index file backed by GCS and powered by USearch.
* **Worker** – Cloud Run service that mounts exactly one shard and executes reads/writes.

When a write would overflow a shard:
1. The worker seals the shard and tells `engine‑frontend` over gRPC.
2. The frontend allocates a new shard and starts a new worker.
3. Any buffered or subsequent writes are retried against the fresh worker.

This design keeps workers stateless & horizontally scalable while maintaining append‑only durability.

---

## Road‑map

* [x] engine-worker with a local backend
* [ ] engine-frontend with a local backend
* [ ] engine-worker with a google-cloud backend
* [ ] engine-frontend with a google-cloud backend
* [ ] CLI
* [ ] BigQuery external routines
* [ ] Collection‑level TTL
* [ ] Terraform / Pulumi modules
* [ ] Kubernetes helm charts (for non‑serverless deployments)

See [open issues](https://github.com/danielealbano/svdb/issues) and feel free to suggest improvements.

---

## Contributing

We ❤️ pull requests!  Start by reading [CONTRIBUTING.md](CONTRIBUTING.md) and opening a discussion if you plan a
substantial change.

---

## License

SVDB is released under the **BSD 3‑Clause License**.  See [LICENSE](LICENSE) for the full text.
