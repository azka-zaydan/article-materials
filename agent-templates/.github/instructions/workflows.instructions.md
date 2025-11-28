# Development Workflows

> **See also**: [`/docs/DEVELOPMENT.md`](../../docs/DEVELOPMENT.md) for complete development guide

## Local Setup

```bash
cp .env.example .env              # Set DATABASE_URL
docker compose up -d              # Start PostgreSQL
npm run migrate:up                # Create all schemas + tables
npm run seed:up                   # Insert demo data
npm run dev                       # Start dev server → http://localhost:5173
```

## Database Credentials (from compose.yaml)

- User: `root`
- Password: `[password]`
- Database: `[db_name]`
- Port: `5432`

## Quality Checks

**ALWAYS run in this order after code changes:**

```bash
npm run format    # Auto-format code
npm run lint      # Check linting
npm run check     # TypeScript + Svelte validation
```

## Testing Workflows

### Migration Testing

```bash
npm run migrate:down  # Rollback
npm run migrate:up    # Reapply
npm run db:studio     # Visual check
```

### API Testing

```bash
# Start dev server
npm run dev

# Test endpoint with curl
curl -X POST http://localhost:5173/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"testuser","email":"test@example.com","password":"Test123!@#"}'
```

## MCP Server Workflow

### Svelte MCP Server

**ALWAYS consult before writing Svelte components:**

1. Call `list-sections` → discover relevant Svelte 5 / SvelteKit docs
2. Call `get-documentation` → fetch detailed API references
3. Write code using current patterns
4. Call `svelte-autofixer` → validate Svelte components
5. Only call `playground-link` if user explicitly requests it

### Context7 MCP Server (Upstash)

**For library documentation:**

1. Call `resolve-library-id` → get Context7 ID (e.g., "drizzle-orm", "lucia-auth")
2. Call `get-library-docs` with library ID + topic → fetch current docs
3. Apply changes using latest API patterns

**Why**: Prevents deprecated APIs, ensures Svelte 5 compliance, avoids outdated patterns.

## Adding New Features

### 1. Create Validation Schema

```typescript
// src/lib/types/feature.schemas.ts
import { z } from 'zod';
import { uuidSchema } from './common.schemas';

export const featureSchema = z.object({
	id: uuidSchema,
	name: z.string().min(1).max(255),
	isActive: z.boolean().default(true)
});

export const createFeatureSchema = featureSchema.omit({ id: true });
export type Feature = z.infer<typeof featureSchema>;
export type CreateFeature = z.infer<typeof createFeatureSchema>;
```

### 2. Create Database Schema

```typescript
// src/lib/server/db/schema/app.schema.ts
export const features = appSchema.table('features', {
	id: uuid('id').defaultRandom().primaryKey(),
	name: varchar('name', { length: 255 }).notNull(),
	isActive: boolean('is_active').default(true).notNull()
});
```

### 3. Generate & Run Migration

```bash
npm run generate:migration -- --type=migrate --name=create_features_table --schema=app
# Edit the generated .up.sql and .down.sql files
npm run migrate:up
```

### 4. Create API Client

```typescript
// src/lib/api/features.ts
import { createFeatureSchema } from '$lib/types/feature.schemas';
import { apiFetch } from './base';

export function createFeaturesApi(fetch: typeof globalThis.fetch) {
	return {
		async create(data: unknown) {
			const validatedData = createFeatureSchema.parse(data);
			return apiFetch(fetch, '/api/features', {
				method: 'POST',
				body: JSON.stringify(validatedData)
			});
		}
	};
}
```

### 5. Create API Endpoint

```typescript
// src/routes/api/features/+server.ts
import { createFeatureSchema } from '$lib/types/feature.schemas';
import { ApiError, ApiSuccess } from '$lib/server/errors';
import { withApiLogging } from '$lib/server/logger/middleware';

export async function POST(event) {
	return withApiLogging(
		event,
		async ({ requestId }) => {
			const body = await event.request.json();
			const result = createFeatureSchema.safeParse(body);

			if (!result.success) {
				return ApiError.fromZod(result.error, requestId);
			}

			const [feature] = await db.insert(features).values(result.data).returning();

			return ApiSuccess.created(feature, { requestId });
		},
		{ operation: 'CREATE_FEATURE', schema: 'app' }
	);
}
```

### 6. Update Documentation

- API changes → update `/docs/API.md`
- Database changes → update `/docs/DATABASE.md`
- New patterns → create new file in `/docs/`
- ❌ NEVER create `.md` files outside `/docs/` (except README.md in root)

### 7. Test & Validate

```bash
npm run format
npm run lint
npm run check
```

## Git Workflow

```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes
git add .
git commit -m "feat: add new feature with validation"

# Push and create PR
git push origin feature/new-feature
```

## Debugging

### Check Database State

```bash
# Connect to database
docker exec -it [container-name] psql -U root -d [db_name]

# List schemas
\dn

# List tables in schema
\dt auth.*
\dt app.*
\dt master.*
```

### Check Logs

```bash
# Application logs are in logs/ directory
tail -f logs/combined.log
tail -f logs/error.log
```

### Drizzle Studio

```bash
# Visual database browser
npm run db:studio
# Opens at http://localhost:4983
```

## Production Build

```bash
npm run build
node build
```

## Environment Variables

Required in `.env`:

```bash
DATABASE_URL="postgres://root:[password]@localhost:5432/[db_name]"
NODE_ENV="development"  # or "production"
```
