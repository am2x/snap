# Snap Roadmap

A phased plan to evolve Snap from a `pg_dump` wrapper into a native, database-agnostic
backup/restore engine — and to use that journey to sharpen database fundamentals
(MVCC, snapshot isolation, wire protocols, WAL, catalogs) alongside production-grade Go.

Each task has a checkbox, a concrete deliverable, and — where relevant — the
underlying concept it's meant to teach. Phases are ordered by dependency: finish
(or mostly finish) a phase before starting the next.

---

## Phase 0 — Hygiene & Hardening (fix what exists)

Goal: make the current codebase safe, testable, and honest before building on it.
No new features in this phase.

### 0.1 Security fixes
- [ ] Remove SQL injection in `restore`: `CREATE DATABASE %s` is built with `fmt.Sprintf`
      from a user-supplied flag. Validate the identifier (or quote it properly) before
      interpolating. *Concept: identifiers can't be parameterized — learn why, and what
      safe quoting of identifiers vs. values means.*
- [ ] Stop persisting plaintext passwords. `connect` writes credentials to `config.env`
      via `viper.WriteConfigAs`. Options (pick one, document the tradeoff in NOTES.md):
      - support `~/.pgpass`-style file with `0600` permissions,
      - read password interactively via `golang.org/x/term` (no echo),
      - integrate OS keychain (e.g. `zalando/go-keyring`).
- [ ] Never put passwords in process env of child commands where avoidable
      (`PGPASSWORD` is visible in `/proc` on Linux). Prefer connection-string or
      pgpass approaches. (Moot once Phase 2 removes shell-outs entirely.)

### 0.2 Error handling & logging discipline
- [ ] Replace every `_ =` swallowed error (`cmd.Flags().GetString`, `LoadConfig`)
      with checked errors.
- [ ] Remove `log.Fatalf` from command `Run` functions — return errors instead.
      Switch commands from `Run:` to `RunE:` so Cobra surfaces errors and exit codes
      consistently.
- [ ] Add a `--verbose` persistent flag that controls the logrus level; default to
      warn-level, human-readable output.

### 0.3 Structure & config cleanup
- [ ] Delete dead code: unused stubs in `application.go` (or move them to a design
      doc), commented-out blocks, the unused `pgx v3` dependency (Phase 2 will add
      `pgx/v5` fresh).
- [ ] Flag bug: `restoreDbCmd` and `connectDbCmd` share the package-level `name`
      variable — bound flags stomp each other. Give each command its own flag
      variables or read flags via `cmd.Flags().GetString` only.
- [ ] `connect` writes config *before* verifying the connection succeeds — reorder
      so a failed ping doesn't persist bad config.
- [ ] Introduce an `internal/` layout — expanded in **0.3a** below.
- [ ] Support multiple named connection profiles (`snap connect --profile prod`),
      stored under `~/.config/snap/`.

### 0.3a Package layout — detailed plan

#### Why root-level `internal/`, and what's wrong with the current tree

Current state (as of starting this task):

```
snap/
├── main.go                 — package main, calls cmd.Execute()
├── cmd/
│   ├── root.go             — root command + ALL flag wiring in one init()
│   ├── database.go         — connect, ssh-connect, backup, restore (500+ lines)
│   └── internal/           — ⚠ misplaced: empty config/ and engine/ dirs
├── helpers/
│   └── dialect.go          — identifier quoting (good code, wrong home)
└── utils/                  — ⚠ empty grab-bag dir
```

Three placement problems to fix:

1. **`cmd/internal/` must move to repo-root `internal/`.** The Go compiler
   enforces `internal/` visibility relative to its *parent* directory: code under
   `cmd/internal/...` is importable only by packages under `cmd/`. That happens
   to work today, but the moment `main.go` (or a future `pkg/`, or tests outside
   `cmd/`) needs config or engine types, the import is a compile error. Root-level
   `internal/` means "importable by this module only", which is the intent.
   Reference: https://go.dev/doc/go1.4#internal-packages and
   https://go.dev/doc/modules/layout ("Package or command with supporting packages").
2. **Kill `helpers/` and `utils/`.** Package names should describe *what they
   provide*, not *what they are* (`helpers`, `utils`, `common`, `misc` all fail
   this). `dialect.go` is database identifier handling → it belongs in
   `internal/db`. The empty `utils/` dir should be deleted before it attracts
   code. References: https://go.dev/blog/package-names,
   Dave Cheney "Practical Go" §2 (package naming):
   https://dave.cheney.net/practical-go/presentations/practical-go.html#_identifiers.
3. **`cmd/` as the Cobra package is fine — but decide once.** Cobra's generator
   uses `cmd/` as a package of commands (what you have). The other convention —
   `cmd/snap/main.go` per-binary entrypoints — matters only for multi-binary
   repos. Keep the current shape; note the alternative and move on. Do NOT adopt
   the full `golang-standards/project-layout` tree — it is community lore, not
   official, and oversized for a single-binary CLI (see the maintainer discussion
   in https://github.com/golang-standards/project-layout/issues/117).

#### Target tree (end of Phase 0)

```
snap/
├── main.go                       — entrypoint: cmd.Execute() only. Move the
│                                   timezone hack out (or delete it — forcing
│                                   time.Local to Africa/Nairobi will stamp
│                                   backups with local time; prefer time.UTC
│                                   in filenames, formatted at display time)
├── cmd/                          — Cobra layer ONLY: declare flags, parse, call internal/*
│   ├── root.go                   — rootCmd, Execute(), persistent --verbose flag,
│   │                               SilenceUsage: true (see error-fix list below)
│   ├── connect.go                — connectCmd: flags → config.Profile → engine.Ping
│   ├── backup.go                 — backupCmd: flags → engine.Backup
│   ├── restore.go                — restoreCmd: flags → engine.Restore
│   ├── ssh.go                    — sshConnectCmd (until Phase 4 replaces it)
│   └── version.go                — versionCmd
├── internal/
│   ├── config/
│   │   ├── config.go             — Profile struct {Name, Server, Port, User, Password},
│   │   │                           Load(path)/Save(path) — the ONLY file that
│   │   │                           imports viper; everything else takes a Profile
│   │   ├── validate.go           — Profile.Validate(): non-empty name, numeric port 1–65535
│   │   └── config_test.go        — Save→Load round-trip via t.TempDir()
│   ├── engine/
│   │   ├── engine.go             — orchestration API the commands call:
│   │   │                             Ping(ctx, Profile) error
│   │   │                             Backup(ctx, Profile, BackupOptions) error
│   │   │                             Restore(ctx, Profile, RestoreOptions) error
│   │   │                           context.Context threaded from day one —
│   │   │                           Phase 1 cancellation slots in with no refactor
│   │   ├── backup.go             — today: wraps pg_dump via exec; Phase 2 swaps
│   │   │                           the body for pgx, signatures unchanged
│   │   └── restore.go            — today: wraps psql via exec
│   └── db/
│       ├── dialect.go            — moved verbatim from helpers/dialect.go
│       │                           (package db; call sites: db.GetDialect)
│       └── dialect_test.go       — table-driven: names with ", `, ], unicode,
│                                   empty string; assert quoted output
├── backups/                      — data, unchanged (gitignore it — dumps don't
│                                   belong in history; backups/testdb.sql is
│                                   already committed, remove it)
├── Makefile                      — currently 0 bytes; fill in (0.4)
└── .golangci.yml                 — see error-fix list: current filename is wrong
```

Layering rule to enforce from now on: `cmd → internal/engine → internal/db`,
and `cmd → internal/config`. Nothing imports `cmd`; `internal/db` imports no
sibling. If an import needs to point "sideways" or "up", the type is in the
wrong package.

#### Migration steps (each is one small commit)

1. `git mv helpers/dialect.go internal/db/dialect.go`; change `package helpers`
   → `package db`; fix the import in `cmd/database.go`
   (`github.com/AllanM007/snap/internal/db`); `rmdir helpers utils`;
   `go build ./...`.
2. Add `internal/db/dialect_test.go` (first real test in the repo — proves the
   test wiring before anything depends on it).
3. `git mv cmd/internal/config internal/config` and same for `engine` (they're
   empty, so this is just directory moves); create `internal/config/config.go`
   by extracting `DbConfig` + `LoadConfig` out of `cmd/database.go`; rename
   `DbConfig` → `Profile` (it will soon be one of several named profiles).
4. Add `Profile.Validate()` + `config_test.go`.
5. Create `internal/engine` with `Ping/Backup/Restore` and move the `exec`
   bodies from the command closures into it. The Cobra `RunE` funcs shrink to:
   parse flags → build `Profile`/options → call engine → return error.
6. Split `cmd/database.go` into per-command files (`connect.go`, `backup.go`,
   `restore.go`, `ssh.go`). Move each command's flag registration out of
   `root.go`'s `init()` into an `init()` in its own file — that's what makes
   the shared-variable bug (below) structurally impossible to reintroduce.
7. Delete `application.go` remnants if any survive, plus stale commented code.

Order matters: 1–2 are risk-free warm-ups; 5 is the real extraction; 6 is
mechanical once 5 is done. Run `go build ./... && go test ./...` between every
step — no step should leave the tree broken.

#### Error fixes found in the current WIP diff (fold into this task)

- [ ] **`restore`: swallowed error.** `cmd/database.go` now has
      `restoreDBName, err := cmd.Flags().GetString("name"); if err != nil { return nil }`
      — returns `nil` on error, silently succeeding. Must be `return err`. (This
      is the exact bug class 0.2 exists to kill; worth a linter — `nilerr` in
      golangci-lint catches it.)
- [ ] **`restore`: reused `exec.Cmd`.** The create-then-retry path calls
      `restoreCommand.Run()` on a command that already ran via
      `restoreCommand.CombinedOutput()`. An `exec.Cmd` is single-use — the retry
      always fails with `exec: already started`. Build a fresh
      `exec.Command(...)` for the retry. (Pre-existing bug, but the refactor into
      `engine/restore.go` is the moment to fix it. The whole pipe dance can also
      go: run create, wait for it, then run restore — nothing needs the pipe.)
- [ ] **`connect`: still `log.Fatalf` inside `RunE`.** `log.Fatalf` calls
      `os.Exit(1)` — deferred cleanups don't run and the RunE error path is
      bypassed. Return the error. Same in `backup` (`Failed to create backup
      file`, `Command execution failed`).
- [ ] **`connect`: `viper.WriteConfigAs("config.env")` error ignored**, and still
      written before the ping (your TODO). Both resolved naturally by step 3/5:
      `config.Save` returns an error, and the engine call order becomes
      ping-then-save.
- [ ] **Double error reporting.** `ssh-connect` does `log.Errorf(...)` then
      `return err` — with RunE, Cobra prints the error again. Pick one channel:
      return errors (wrapped with `fmt.Errorf("...: %w", err)` for context) and
      let the top level print. Set `SilenceUsage: true` on rootCmd so runtime
      errors don't dump the help text, and optionally `SilenceErrors: true` with
      a single `fmt.Fprintln(os.Stderr, err)` in `Execute()`. Reference:
      https://github.com/spf13/cobra/blob/main/site/content/user_guide.md
- [ ] **Lint config filename.** `.golangci-lint.yaml` is not a name golangci-lint
      auto-detects — it looks for `.golangci.{yml,yaml,toml,json}`. Rename, and
      fill the currently-empty file. Useful starter set for this codebase:
      `errcheck`, `govet`, `staticcheck`, `nilerr`, `errorlint`, `noctx`,
      `unparam`, `revive`. Reference: https://golangci-lint.run/usage/configuration/
- [ ] **Empty `Makefile`** (0 bytes) — minimum viable targets:
      `build` (`go build -o bin/snap .`), `test` (`go test ./...`),
      `lint` (`golangci-lint run`), `vet`, `tidy`.
- [ ] **`main.go` timezone override** — `time.Local = loc` mutates a process
      global and will surprise every timestamp in the app; drop it and use
      explicit locations where display formatting needs one.

#### References for this task (in reading order)

1. https://go.dev/doc/modules/layout — official module layout guidance (short;
   the single source of truth for where `internal/` goes).
2. https://go.dev/blog/package-names — why `helpers`/`utils` are anti-patterns
   and what to name things instead.
3. https://dave.cheney.net/practical-go/presentations/practical-go.html —
   §2 (naming) and §5 (API design); directly applicable to the engine API shape.
4. Ben Johnson, "Standard Package Layout"
   https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1 —
   the "root package for domain types, subpackages for dependencies" idea behind
   putting `Profile` in `config` and keeping viper quarantined there.
5. https://github.com/spf13/cobra/blob/main/site/content/user_guide.md — RunE,
   SilenceUsage/SilenceErrors, per-file command organization.
6. https://golangci-lint.run/usage/configuration/ — config discovery + linter list.

**Exit criteria for 0.3a:** `cmd/*.go` contain no `exec.Command` and no viper
calls; `go build ./...` clean; dialect + config tests pass; `helpers/`, `utils/`,
`cmd/internal/` gone.

### 0.4 Testing foundation
- [ ] Add first unit tests (Testify is already in the README promise): config
      load/save round-trip, identifier validation, flag parsing.
- [ ] Add a `Makefile` (or `Taskfile`): `build`, `test`, `lint`, `vet`.
- [ ] Add `golangci-lint` config and fix what it finds.
- [ ] Set up GitHub Actions CI: build + vet + lint + unit tests on push.

**Exit criteria:** no secrets on disk in plaintext, no swallowed errors, CI green,
`go vet`/linter clean.

---

## Phase 1 — Make the Wrapper Excellent (UX before internals)

Goal: even while still shelling out, make Snap feel like a real tool. This phase is
cheap and keeps the project usable throughout the rewrite.

- [ ] `snap backup`: add `--output/-o` flag (default `backups/<name>-<timestamp>.sql`)
      — timestamped files instead of silently overwriting the previous backup.
- [ ] Add gzip compression: `--compress` writes `.sql.gz` via `gzip.Writer` wrapped
      around the file. *Concept: composing `io.Writer`s — you'll reuse this pattern
      for the native engine.*
- [ ] Add progress reporting: wrap the output file in a counting writer, render
      bytes-written + elapsed (simple ticker; a TUI library is optional).
- [ ] Add `context.Context` + signal handling: Ctrl-C should cancel the child
      process cleanly (`exec.CommandContext`) and remove the partial backup file.
      *Concept: cancellation propagation — foundational for everything later.*
- [ ] `snap restore`: add `--file/-f` to restore an arbitrary dump; confirm before
      restoring into an existing non-empty database (`--yes` to skip).
- [ ] `snap list`: list local backups with size, age, and target database.
- [ ] `snap verify` (v0): after backup, run the dump through a cheap sanity check
      (non-empty, ends with expected `pg_dump` trailer, gzip integrity if compressed).
- [ ] Version command: inject version via `-ldflags "-X ..."` at build time instead
      of the hardcoded string.

**Exit criteria:** demo-quality UX — timestamped compressed backups with progress,
clean cancellation, list/verify commands.

---

## Phase 2 — Native Postgres Engine (the core learning phase)

Goal: remove `pg_dump`/`psql` entirely for the default path. Talk to Postgres over
the wire with `pgx/v5`. This is where the database fundamentals live — take it slow
and read the docs referenced in each task.

### 2.1 Connection layer
- [ ] Add `jackc/pgx/v5` + `pgxpool`. Implement `internal/db/postgres` with
      connect/ping using a proper connection string built from the profile.
- [ ] Replace `snap connect`'s shell-out ping with a real `Ping(ctx)` +
      `SELECT version()`; print server version and current database.
      *Concept: the Postgres startup/auth handshake (read: pgx docs on auth,
      SCRAM-SHA-256).*

### 2.2 Schema introspection (the catalog deep-dive)
- [ ] Enumerate tables: query `information_schema.tables` /
      `pg_catalog.pg_class` for user tables in dependency-safe order.
- [ ] For each table, reconstruct DDL: columns with types/defaults/nullability
      (`information_schema.columns`), primary keys, unique constraints, indexes
      (`pg_indexes`), sequences and their current values, foreign keys
      (`pg_constraint`), comments.
- [ ] Build a dependency graph from FKs and topologically sort it so restore can
      create/load tables in a valid order. Detect cycles → defer FK creation to a
      post-data step (this is exactly what `pg_dump` does — verify by reading one
      of its dumps). *Concept: system catalogs, DDL round-tripping, graph ordering.*
- [ ] Unit-test DDL generation against fixture schemas (golden files).

### 2.3 Consistent data export
- [ ] Open a transaction with `ISOLATION LEVEL REPEATABLE READ READ ONLY` and run
      the *entire* backup (schema read + all table data) inside it.
      *Concept: MVCC and snapshot isolation — why this gives a consistent
      point-in-time backup without locking writers. Write a short explainer in
      `docs/consistency.md` in your own words; that's the proof you've got it.*
- [ ] Export table data using the COPY protocol (`conn.PgConn().CopyTo(ctx, w,
      "COPY tbl TO STDOUT")`) streaming straight into your compressing writer —
      constant memory regardless of table size. *Concept: the COPY subprotocol vs.
      row-by-row SELECT; measure the throughput difference and record it.*
- [ ] Define the archive format v1: a tar (or zip) containing `manifest.json`
      (snap version, server version, timestamp, table list, row counts, checksums),
      `schema/pre-data.sql`, `data/<table>.copy.gz` per table, `schema/post-data.sql`
      (indexes, FKs, sequence setvals).
- [ ] Checksum each data file (SHA-256) into the manifest as it streams
      (`io.MultiWriter` with a hasher). Upgrade `snap verify` to check them.

### 2.4 Native restore
- [ ] Restore pipeline: create database if missing (using the safe identifier
      handling from 0.1) → apply pre-data DDL → `CopyFrom` each table's data →
      apply post-data DDL (indexes/FKs last — *concept: why building indexes after
      load is dramatically faster; measure it*) → set sequences.
- [ ] Wrap restore in a transaction where possible; document which DDL cannot be
      transactional and how you handle partial failure (drop-and-retry vs. resume).
- [ ] Parallelize data load across tables with an `errgroup` bounded by
      `--jobs/-j N` connections. *Concept: connection pooling, and Go worker-pool
      patterns with proper ctx cancellation.*

### 2.5 Correctness harness (the payoff)
- [ ] Integration tests with `testcontainers-go`: start Postgres, load a fixture
      schema + data, `snap backup`, restore into a *second* fresh container, then
      diff: row counts, per-table checksums (`SELECT md5(string_agg(...))` or
      client-side hash), schema dump comparison.
- [ ] Torture fixture: text with embedded quotes/newlines/NULL bytes-adjacent
      cases, bytea, arrays, enums, numeric extremes, timestamps with zones, a
      self-referencing FK, and an FK cycle. Every one of these will break a naive
      implementation — that's the point.
- [ ] Consistency-under-writes test: run a writer goroutine mutating rows during
      backup; assert the restored data matches a single snapshot (no torn state).
      *This test failing/passing is the single best MVCC lesson in the project.*
- [ ] Benchmark: seed ~5–10GB, record backup/restore wall time and max RSS in
      `docs/benchmarks.md`; compare against `pg_dump`/`pg_restore` honestly.

**Exit criteria:** `snap backup`/`snap restore` need no Postgres client binaries,
pass the torture suite, survive the concurrent-writes test, and hold constant
memory on a 10GB table.

---

## Phase 3 — Database-Agnostic Core

Goal: make "database agnostic" true via a driver interface, and prove it with a
second engine.

- [ ] Extract the driver interface from the Postgres implementation (interface
      *after* one concrete implementation, not before). Sketch:
      ```go
      type Driver interface {
          Ping(ctx context.Context) error
          Snapshot(ctx context.Context) (Snapshot, error) // consistent read view
      }
      type Snapshot interface {
          Schema(ctx context.Context) (*Schema, error)
          ExportTable(ctx context.Context, t Table, w io.Writer) (Stats, error)
          Close() error
      }
      type Restorer interface {
          Prepare(ctx context.Context, s *Schema) error
          ImportTable(ctx context.Context, t Table, r io.Reader) error
          Finalize(ctx context.Context, s *Schema) error // indexes, FKs, sequences
      }
      ```
- [ ] Define a portable logical data format alongside the native COPY format
      (NOTES.md already suggests csv/json): typed column metadata + streamed rows,
      so a Postgres backup can be restored into MySQL. Document the type-mapping
      table (`docs/type-mapping.md`) — *this is where you learn how differently
      databases model types (e.g. Postgres `text` vs MySQL `TEXT/VARCHAR`,
      timezone handling, boolean vs tinyint).*
- [ ] Implement the MySQL driver (`go-sql-driver/mysql`):
      - consistent snapshot via `START TRANSACTION WITH CONSISTENT SNAPSHOT`
        (InnoDB) — *compare its semantics with Postgres REPEATABLE READ and write
        up the differences*,
      - schema introspection via `information_schema`,
      - bulk load via `LOAD DATA LOCAL INFILE` streaming (the driver supports
        reader registration) or batched multi-row INSERTs; measure both.
- [ ] SQLite driver (optional but cheap and instructive): file-based, uses
      `sqlite3 .backup`-equivalent semantics via the serialize API or VACUUM INTO.
- [ ] Cross-engine integration test: Postgres → portable format → MySQL and back;
      diff data allowing for documented type-mapping lossiness.
- [ ] `snap backup --db postgres|mysql|sqlite` selects the driver; auto-detect
      from the profile.

**Exit criteria:** two real drivers behind one interface; a cross-engine
restore passes CI; type-mapping doc exists.

---

## Phase 4 — Remote & Storage Targets

Goal: make Snap useful beyond localhost. (This replaces the current `ssh-connect`
command, which only runs a hardcoded `docker ls`.)

- [ ] Storage backend interface (`io`-centric): `Put(ctx, key, r io.Reader)`,
      `Get`, `List`, `Delete`. Implement `local/` first (current behavior).
- [ ] S3 backend (AWS SDK v2, works with MinIO for local testing via
      testcontainers): multipart streaming upload directly from the backup
      pipeline — never buffer the archive on local disk unless asked.
      *Concept: streaming end-to-end; your Phase 1/2 `io.Writer` composition pays
      off here.*
- [ ] SSH tunnel support done properly: use `golang.org/x/crypto/ssh` (already an
      indirect dep) to dial a tunnel and run the *native* engine through it —
      `snap backup --profile prod` where the profile specifies a bastion. No
      remote shell commands, no remote binaries needed. *Concept: dialers and
      `net.Conn` plumbing; pgx accepts a custom DialFunc.*
- [ ] Docker convenience: detect `container: true` profiles and connect via the
      container's published port (document why exec-ing into the container is no
      longer needed once the engine is native).
- [ ] Retention policy: `snap prune --keep-last N --keep-daily D --keep-weekly W`
      against any storage backend.
- [ ] Client-side encryption at rest: `age` (filippo.io/age) encryption of the
      archive stream, key from file or passphrase. Document the threat model in
      `docs/security.md`.

**Exit criteria:** backup straight to S3/MinIO through an SSH tunnel with
encryption, and prune retention — all in one process, no external binaries.

---

## Phase 5 — Incremental & Point-in-Time (the WAL deep-dive)

Goal: the most advanced phase; each item is independently valuable. This is
senior-level database territory.

- [ ] Incremental logical backup v0: per-table watermark strategy (updated-at or
      PK-range based) with an honest doc about what it can and cannot guarantee
      (deletes! schema changes!). *Concept: why logical incrementals are hard,
      and why real tools use the WAL instead.*
- [ ] Read up on Postgres WAL and logical decoding; write `docs/wal-notes.md`
      summarizing: what the WAL is, LSNs, replication slots, logical vs physical
      replication. (Reading task — do it before the next two items.)
- [ ] Continuous change capture: create a logical replication slot
      (`pgoutput` or `wal2json`) and consume changes via pgx's replication
      support into an append-only changelog file per backup. *Concept: WAL,
      replication slots, exactly the mechanism behind CDC tools like Debezium.*
- [ ] Point-in-time restore v0: restore base backup + replay changelog up to a
      target LSN/timestamp.
- [ ] `snap watch`: long-running mode — periodic base backups + continuous WAL
      streaming, with slot-lag monitoring (a forgotten slot bloats the server's
      disk; surface that metric loudly). *Concept: the operational failure modes
      of replication slots.*
- [ ] Scheduling: document cron/systemd-timer usage rather than building a
      scheduler daemon (keep the tool composable, Unix-style).

**Exit criteria:** a demo script that provably restores a database to "as of
14:32:05" after simulated data loss.

---

## Phase 6 — Release Quality

Goal: make it something other people can install, trust, and read.

- [ ] Rewrite README: what/why, install, quickstart GIF, architecture diagram
      (mermaid), comparison table vs `pg_dump` (be honest about limitations).
- [ ] `docs/` index for all the concept write-ups produced in earlier phases —
      these documents are the portfolio artifact as much as the code.
- [ ] goreleaser: cross-platform binaries, homebrew tap, checksums, changelog.
- [ ] Semantic versioning + tagged releases; `snap version` shows commit/date.
- [ ] Structured JSON log mode (`--log-format json`) for cron use.
- [ ] Exit-code contract documented (0 ok, 1 error, 2 usage, 3 verify-failed …).
- [ ] Fuzz tests on parser-ish surfaces (identifier validation, manifest parsing)
      with Go's native fuzzing.
- [ ] Optional: metrics hook (Prometheus textfile or push) for backup duration,
      size, and last-success timestamp — the numbers an on-call engineer needs.

---

## Suggested order of attack

Phases 0 → 1 are a weekend or two each. Phase 2 is the heart of the project —
budget the most time there and don't skip 2.5 (the correctness harness is where
backup tools earn trust). Phase 3 proves the abstraction, Phase 4 makes it
practical, Phase 5 makes it impressive. Phase 6 can start any time after Phase 2
and grow incrementally.

Rule of thumb throughout: **every phase ends with a test that could have caught
the bugs of the previous phase**, and every concept learned gets a short write-up
in `docs/` — explaining it in your own words is the difference between having
used a feature and having learned a fundamental.
