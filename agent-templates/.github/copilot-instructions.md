# Copilot Instructions

> **Note**: This is the main instruction file. Detailed guides are in `.github/instructions/`.

## Project Overview

[Project Name] built with SvelteKit 2 + Svelte 5, Drizzle ORM, and PostgreSQL with **multi-schema architecture**. Features Lucia authentication, custom SQL migration system, Winston logging, type-safe API client, and schema-first validation.

## Tech Stack

- **Framework**: SvelteKit 2 with Svelte 5 (runes API) + TypeScript strict mode
- **Styling**: Tailwind CSS 4 via Vite plugin (`@tailwindcss/vite`)
- **Database**: PostgreSQL + Drizzle ORM with postgres-js driver
- **Auth**: Lucia + Argon2 password hashing
- **Validation**: Zod v4 for runtime type-safe validation (schema-first approach)
- **Logging**: Winston with custom formatters and request tracking
- **Deployment**: Node.js adapter (`@sveltejs/adapter-node`)

## MANDATORY RULES

### Before Code Changes

1. **ALWAYS consult MCP servers first**:
   - Svelte MCP: `list-sections` → `get-documentation` → `svelte-autofixer`
   - Context7 MCP: `resolve-library-id` → `get-library-docs` for library APIs
2. **Check existing documentation** in `/docs` and `.github/instructions/`
3. **Verify schema qualification** for all SQL operations (use `auth.user`, not `user`)

### During Code Changes

4. **Use Svelte 5 runes exclusively** - See [svelte5.instructions.md](instructions/svelte5.instructions.md)
5. **Use Zod v4 validation** - See [validation.instructions.md](instructions/validation.instructions.md)
6. **Use API client pattern** - See [api-client.instructions.md](instructions/api-client.instructions.md)
7. **Wrap API endpoints with logging** - See [logging.instructions.md](instructions/logging.instructions.md)
8. **Wrap database queries with logging** - See [database.instructions.md](instructions/database.instructions.md)
9. **Follow design system principles** - See [ui-design.instructions.md](instructions/ui-design.instructions.md)
10. **Avoid unsafe casts**: Never use `as any`. Prefer schema-inferred types (`z.infer<>`) and narrow values via Zod `.safeParse()`/`.parse()`.
11. **Fixing bugs? Consult docs first**: When fixing bugs/errors, ALWAYS use Svelte MCP and Context7 before coding.
12. **Follow DRY principle**: Extract duplicate code into reusable functions. Place server-only utilities in `$lib/server/`, shared utilities in `$lib/utils/`, and component logic in `$lib/components/`.

### After Code Changes

13. **ALWAYS run validation**:

```bash
npm run format    # Auto-format code
npm run lint      # Check linting
npm run check     # TypeScript + Svelte validation
```

14. **Update documentation** for new features and notable fixes
15. **Test migrations** if database schema changed - See [database.instructions.md](instructions/database.instructions.md)
16. **Verify UI follows design system** - Check against [ui-design.instructions.md](instructions/ui-design.instructions.md)

## Core Patterns (Quick Reference)

### Schema-First Validation

```typescript
// 1. Create Zod schema in src/lib/types/*.schemas.ts
export const resourceSchema = z.object({ ... });
export type Resource = z.infer<typeof resourceSchema>;

// 2. Validate client-side before API call
const result = createResourceSchema.safeParse(formData);

// 3. Validate backend in API endpoint
const result = createResourceSchema.safeParse(body);
if (!result.success) {
  return ApiError.fromZod(result.error, requestId);
}
```

See [validation.instructions.md](instructions/validation.instructions.md) for complete guide.

### Type-Safe API Client

```typescript
// Always use createApiClient(fetch)
import { createApiClient } from '$lib/api/client';

export const load: PageServerLoad = async ({ fetch }) => {
	const api = createApiClient(fetch);
	const data = await api.resource.getAll({ limit: 10 });
	return { data };
};
```

See [api-client.instructions.md](instructions/api-client.instructions.md) for adding new endpoints.

### Multi-Schema Database

```typescript
// Use schema-qualified names
export const user = authSchema.table('user', { ... });
export const items = appSchema.table('items', {
  userId: uuid('user_id').references(() => user.id)  // Cross-schema FK
});
```

See [database.instructions.md](instructions/database.instructions.md) for schema update workflow.

### Svelte 5 Components

```svelte
<script lang="ts">
	let { data, children } = $props(); // ✅ Svelte 5
	let count = $state(0); // ✅ Reactive state
	let doubled = $derived(count * 2); // ✅ Computed

	// ❌ Never: export let, $:, <slot>
</script>
```

See [svelte5.instructions.md](instructions/svelte5.instructions.md) for all patterns.

## Project Structure

```
src/
├── lib/
│   ├── api/            # Type-safe API client ($api)
│   ├── components/     # Reusable Svelte components ($components, $ui)
│   │   ├── ui/
│   │   │   ├── common/     # Base/reusable (actions, display, forms)
│   │   │   ├── forms/      # Domain-specific forms
│   │   │   ├── [feature]/  # Feature-specific components
│   │   │   └── layout/     # Layout components ($ui/layout)
│   │   └── index.ts    # Component exports
│   ├── server/         # SERVER-ONLY ($server, $db, $schema)
│   │   ├── auth.ts, errors.ts, logger/
│   │   └── db/         # Database ($db)
│   │       └── schema/ # Drizzle schemas ($schema)
│   ├── types/          # Zod schemas & inferred types ($types)
│   └── utils/          # Shared utilities ($utils)
├── routes/
│   ├── api/            # JSON API endpoints
│   ├── (pages)/        # Application pages
│   └── +layout.svelte  # Root layout
└── hooks.server.ts     # Session validation

migrations/
├── migrate/            # Database schema migrations
└── seed/               # Demo data seeds

.github/
├── copilot-instructions.md  # This file
└── instructions/            # Detailed guides
    ├── validation.md        # Zod v4 validation
    ├── api-client.md        # API client usage
    ├── database.md          # Multi-schema DB
    ├── svelte5.md           # Svelte 5 patterns
    ├── workflows.md         # Development workflows
    └── logging.md           # Winston logging
```

## Path Aliases

**ALWAYS use path aliases** instead of relative imports:

```typescript
// ✅ CORRECT - Use path aliases
import { createApiClient } from '$api/client';
import type { Resource } from '$types/app.schemas';
import { db } from '$db';
import { items } from '$schema';
import { Button } from '$components';
import Navigation from '$ui/layout/Navigation.svelte';

// ❌ WRONG - Don't use relative paths
import { createApiClient } from '../../../lib/api/client';
```

**Available aliases**:

- `$components` → Component index exports
- `$ui` → Direct component imports
- `$ui/layout` → Layout components
- `$api` → API client functions
- `$types` → Zod schemas & types
- `$server` → Server-only code
- `$db` → Database client
- `$schema` → Drizzle schemas
- `$utils` → Utilities

## Quick Start

```bash
cp .env.example .env
docker compose up -d
npm run migrate:up
npm run seed:up
npm run dev  # → http://localhost:5173
```

See [workflows.instructions.md](instructions/workflows.instructions.md) for complete setup.

## Common Pitfalls

1. **Schema qualification**: SQL files MUST use `auth.user`, not just `user`
2. **Svelte 4 syntax**: No `export let`, `<slot>`, or `$:` reactive statements
3. **Raw fetch calls**: Use API client instead of `fetch('/api/...')`
4. **Type duplication**: Use `z.infer<>` instead of manual type definitions
5. **Missing validation**: Validate on both client and backend
6. **Client-side DB access**: Never import `$lib/server/*` in `.svelte` files

## Documentation

### AI Agent Instructions

- [validation.instructions.md](instructions/validation.instructions.md) - Zod v4 validation system
- [api-client.instructions.md](instructions/api-client.instructions.md) - Type-safe API client
- [database.instructions.md](instructions/database.instructions.md) - Multi-schema database
- [svelte5.instructions.md](instructions/svelte5.instructions.md) - Svelte 5 patterns
- [workflows.instructions.md](instructions/workflows.instructions.md) - Development workflows
- [logging.instructions.md](instructions/logging.instructions.md) - Winston logging
- [ui-design.instructions.md](instructions/ui-design.instructions.md) - Design system guide

### Developer Documentation

- `/docs/API.md` - API endpoint reference
- `/docs/DATABASE.md` - Database architecture
- `/docs/DEVELOPMENT.md` - Development guide
- `/docs/ERROR_HANDLING.md` - Error handling patterns
- `/docs/LOGGER.md` - Logging system details
- `/docs/UI_STYLE_GUIDE.md` - [Project Name] design system guide
- `/docs/COMPONENTS.md` - Component library reference

## Authentication

Session-based auth using Lucia with SHA-256 hashed tokens. Session validation happens in `src/hooks.server.ts` on every request.

**Key Files**:

- `src/lib/server/auth.ts` - Lucia setup
- `src/hooks.server.ts` - Session validation
- `src/routes/api/auth/*` - Auth endpoints

## Recent Changes

1. **Type-safe API client** (`src/lib/api/`) with validation
2. **Zod v4 validation schemas** for all entities
3. **Standardized error handling** (`ApiError`, `ApiSuccess`)
4. **Migrated to API client pattern** in all `+page.server.ts`

### Migration from Raw Fetch

**Old**:

```typescript
const response = await fetch('/api/resource');
const data = await response.json();
```

**New**:

```typescript
const api = createApiClient(fetch);
const data = await api.resource.getAll();
```

See [api-client.instructions.md](instructions/api-client.instructions.md) for migration checklist.
