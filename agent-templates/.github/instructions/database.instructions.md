# Multi-Schema Database Architecture

> **See also**: [`/docs/DATABASE.md`](../../docs/DATABASE.md) for complete database architecture reference

**Critical**: This project uses PostgreSQL schemas to separate concerns. All SQL must be schema-qualified.

## Schema Structure

```typescript
// src/lib/server/db/schema/*.schema.ts
export const authSchema = pgSchema('auth'); // user, session
export const masterSchema = pgSchema('master'); // categories, locations
export const appSchema = pgSchema('app'); // items, posts, etc.

// Example table definition
export const user = authSchema.table('user', {
	id: uuid('id').primaryKey(),
	email: varchar('email', { length: 255 }).notNull().unique()
	// ...
});

export const items = appSchema.table('items', {
	id: uuid('id').defaultRandom().primaryKey(),
	userId: uuid('user_id').references(() => user.id, { onDelete: 'cascade' }) // Cross-schema FK
	// ...
});
```

## Schema Responsibilities

- **`auth`**: User authentication and session management
- **`master`**: Reference data (categories, locations) - relatively static
- **`app`**: Application transactional data (items, posts, core business entities)
- **`public`**: Migration tracking (`migration_history` table)

**Why**: Security isolation, separate backup/restore strategies, easier permission management, logical data separation.

## SQL Migration System (NOT Drizzle Kit Migrations)

This project uses **custom SQL migrations** stored in `migrations/migrate/` and `migrations/seed/`.

### File Naming Convention

```
<timestamp>_<description>.<up|down>.sql

Examples:
- 20251124_000001_create_migration_history.up.sql
- 20251124_000001_create_migration_history.down.sql
```

### Critical Commands

```bash
# Database setup (preferred workflow)
docker compose up -d         # Start PostgreSQL
npm run migrate:up           # Run all migrations (creates schemas + tables)
npm run seed:up              # Insert demo data

# Migration management
npm run generate:migration -- --type=migrate --name=add_column --schema=app
npm run migrate:down         # Rollback all migrations
npm run seed:down            # Rollback all seeds

# Development tools
npm run db:studio            # Browse database with Drizzle Studio
```

### Migration Requirements

1. **Schema qualification required**: Always use `auth.user`, `app.items`, `master.categories` in SQL
2. **Idempotency**: Use `IF NOT EXISTS`, `IF EXISTS`, `ON CONFLICT DO NOTHING`
3. **Cross-schema FKs work**: `REFERENCES auth.user(id)` from `app.items` is valid
4. **First migration creates schemas**: `000000_create_schemas.up.sql` must run before all others

### Example Migration

```sql
-- migrations/migrate/20251124_000005_create_items_table.up.sql
CREATE TABLE IF NOT EXISTS app.items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES auth.user(id) ON DELETE CASCADE,
  category_id UUID NOT NULL REFERENCES master.categories(id),
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  price DECIMAL(10, 2),
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- migrations/migrate/20251124_000005_create_items_table.down.sql
DROP TABLE IF EXISTS app.items;
```

## Database Schema Update Workflow

When adding new fields or tables:

### 1. Update Drizzle Schema

```typescript
// src/lib/server/db/schema/app.schema.ts
export const items = appSchema.table('items', {
	// ... existing fields
	isPromoted: boolean('is_promoted').default(false).notNull() // NEW FIELD
});
```

### 2. Create Zod Schema

```typescript
// src/lib/types/app.schemas.ts
export const itemSchema = z.object({
	// ... existing fields
	isPromoted: z.boolean().default(false) // NEW FIELD
});

// Update create/update schemas as needed
export const updateItemSchema = z.object({
	// ... existing fields
	isPromoted: z.boolean().optional() // NEW FIELD
});
```

### 3. Generate Migration

```bash
npm run generate:migration -- --type=migrate --name=add_is_promoted_to_items --schema=app
```

This creates:

- `migrations/migrate/TIMESTAMP_add_is_promoted_to_items.up.sql`
- `migrations/migrate/TIMESTAMP_add_is_promoted_to_items.down.sql`

### 4. Edit Migration Files

```sql
-- .up.sql
ALTER TABLE app.items
ADD COLUMN is_promoted BOOLEAN DEFAULT FALSE NOT NULL;

-- .down.sql
ALTER TABLE app.items
DROP COLUMN is_promoted;
```

### 5. Run Migration

```bash
npm run migrate:down  # Rollback to test down migration
npm run migrate:up    # Apply migration
```

### 6. Update API Endpoints

Update any affected API endpoints to handle the new field:

```typescript
// src/routes/api/items/+server.ts
const result = updateItemSchema.safeParse(body);
// The new isPromoted field is now validated automatically
```

## Database Queries

```typescript
import { db } from '$lib/server/db';
import { items, categories, user } from '$lib/server/db/schema';
import { eq, and } from 'drizzle-orm';

// Cross-schema join (app.items + master.categories + auth.user)
const result = await db
	.select()
	.from(items)
	.leftJoin(categories, eq(items.categoryId, categories.id))
	.leftJoin(user, eq(items.userId, user.id))
	.where(eq(items.isActive, true));
```

## Common Pitfalls

1. **Schema qualification**: SQL files MUST use `auth.user`, not just `user`
2. **Migration order**: Schema creation (`000000_create_schemas`) must be first
3. **Drizzle Kit migrations**: We use custom SQL migrations, NOT `drizzle-kit generate`
4. **Client-side DB access**: Never import `$lib/server/db` in `.svelte` files

## Database Testing

```bash
# Test migration can be applied and rolled back
npm run migrate:down
npm run migrate:up

# Verify schema in database
npm run db:studio

# Or use psql
docker exec -it [container-name] psql -U root -d [db_name]
\dn  # List schemas
\dt auth.*  # List tables in auth schema
```
