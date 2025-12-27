# Adequate

Minimal template for monolithic web applications.

## Philosophy

**One binary, one deployment, no unnecessary complexity.**

This template implements three core principles:

1. **Stateless monolith**: a single Go executable that serves both the API and pre-built static files. Trivial deployment, one-command rollback, horizontal scalability by replication.

2. **Clear front/back separation**: the backend is the source of truth (data + API), the frontend is static (pre-generated HTML/JS/CSS). No server-side rendering after login, no business logic in the client.

3. **Built-in observability**: structured logs, graceful shutdown, health metrics, and tracing of slow requests. Everything is measurable in a single pipeline.

This approach removes premature microservices, fragile SSR, and cascading configurations. It favors predictability, visibility, and stability.

## Tech stack

| Component | Technology | Role |
|-----------|-------------|------|
| **Backend** | Go | REST API, HTTP server, DB migrations |
| **Frontend** | Svelte + SvelteKit | User interface (static build) |
| **Database** | SQLite (WAL mode) | Embedded storage, versioned migrations |
| **SQL generation** | SQLC | Type-safe queries from `.sql` files |
| **Frontend build** | Vite + TypeScript | Bundler and dev server with HMR |
| **Styling** | Tailwind CSS | CSS utilities with a Vite plugin |

## License

MIT
