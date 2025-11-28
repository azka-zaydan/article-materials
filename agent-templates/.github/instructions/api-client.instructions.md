# API Client Architecture (Type-Safe Pattern)

> **See also**: [`/docs/API.md`](../../docs/API.md) for complete API endpoint reference

This project uses a **type-safe API client factory pattern** to make requests to internal API endpoints.

## Benefits

- **Automatic validation** before sending requests
- **Type safety** throughout the request/response cycle
- **Consistent error handling**
- **SSR-compatible** (uses fetch from request event)

## Using the API Client

**ALWAYS use `createApiClient(fetch)` in server-side code**:

```typescript
// src/routes/+page.server.ts
import { createApiClient } from '$api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const api = createApiClient(fetch);

	// Type-safe API calls with automatic validation
	const items = await api.items.getAll({ limit: 10 });
	const categories = await api.categories.getAll();

	// Parallel fetching for better performance
	const [userData, metadata] = await Promise.all([api.users.me(), api.metadata.get()]);

	return { items, categories, userData, metadata };
};
```

## API Client Structure

Located in `src/lib/api/`:

- `client.ts` - Factory function that creates API client
- `base.ts` - Base fetch wrapper with error handling
- `auth.ts` - Authentication endpoints
- `[resource].ts` - Resource CRUD operations (e.g. listings, posts, products)

## How API Methods Work

**All API methods validate input** before sending:

```typescript
// src/lib/api/items.ts
export function createItemsApi(fetch: typeof globalThis.fetch) {
	return {
		async create(data: unknown): Promise<ItemDetailResponse> {
			// Validates using createItemSchema before sending
			const validatedData = createItemSchema.parse(data);

			return apiFetch<ItemDetailResponse>(fetch, '/api/items', {
				method: 'POST',
				body: JSON.stringify(validatedData)
			});
		},

		async getAll(params?: ItemQuery): Promise<ApiItemsResponse> {
			const validatedParams = params ? itemQuerySchema.parse(params) : {};
			const queryString = buildQueryString(validatedParams);
			return apiFetch<ApiItemsResponse>(fetch, `/api/items${queryString}`);
		}
	};
}
```

## Adding New API Endpoints

### 1. Create Zod Schemas

```typescript
// src/lib/types/comments.schemas.ts
// Centralize ALL schemas here, including request bodies and query params.
export const createCommentSchema = z.object({
	itemId: uuidSchema,
	content: z.string().min(1).max(500)
});

export const commentsByItemQuerySchema = z.object({
	itemId: uuidSchema
});

export const commentResponseSchema = z.object({
	id: uuidSchema,
	itemId: uuidSchema,
	userId: uuidSchema,
	content: z.string(),
	createdAt: z.coerce.date()
});

export type CreateComment = z.infer<typeof createCommentSchema>;
export type CommentsByItemQuery = z.infer<typeof commentsByItemQuerySchema>;
export type CommentResponse = z.infer<typeof commentResponseSchema>;
```

### 2. Create API Client Module

```typescript
// src/lib/api/comments.ts
import { createCommentSchema, type CommentResponse } from '$lib/types/comments.schemas';
import { apiFetch } from './base';

export function createCommentsApi(fetch: typeof globalThis.fetch) {
	return {
		async create(data: unknown): Promise<CommentResponse> {
			const validatedData = createCommentSchema.parse(data);
			return apiFetch<CommentResponse>(fetch, '/api/comments', {
				method: 'POST',
				body: JSON.stringify(validatedData)
			});
		},

		async getByItem(itemId: string): Promise<CommentResponse[]> {
			return apiFetch<CommentResponse[]>(fetch, `/api/comments?itemId=${itemId}`);
		}
	};
}
```

### 3. Add to Main Client

```typescript
// src/lib/api/client.ts
import { createCommentsApi } from './comments';

export function createApiClient(fetch: typeof globalThis.fetch) {
	return {
		auth: createAuthApi(fetch),
		items: createItemsApi(fetch),
		comments: createCommentsApi(fetch) // Add new API
	};
}
```

### 4. Create Backend API Route

```typescript
// src/routes/api/comments/+server.ts
import { createCommentSchema } from '$lib/types/comments.schemas';
import { ApiError, ApiSuccess } from '$lib/server/errors';
import { withApiLogging } from '$lib/server/logger/middleware';

export async function POST(event) {
	return withApiLogging(
		event,
		async ({ requestId }) => {
			const body = await event.request.json();
			const result = createCommentSchema.safeParse(body);

			if (!result.success) {
				return ApiError.fromZod(result.error, requestId);
			}

			// Create comment in database
			const [comment] = await db.insert(comments).values(result.data).returning();

			return ApiSuccess.created(comment, { requestId });
		},
		{ operation: 'CREATE_COMMENT', schema: 'app' }
	);
}
```

## Migration from Raw Fetch

### Checklist

- [ ] Identify all `fetch('/api/...')` calls in `+page.server.ts` files
- [ ] Import `createApiClient` from `$lib/api/client`
- [ ] Create API client: `const api = createApiClient(fetch)`
- [ ] Replace `fetch()` with typed API call
- [ ] Remove manual `response.json()` calls
- [ ] Remove manual error handling (client handles it)
- [ ] Update type annotations to use schema-inferred types
- [ ] Test the endpoint

### Example Migration

**Before (raw fetch):**

```typescript
export const load: PageServerLoad = async ({ fetch, params }) => {
	const response = await fetch(`/api/items/${params.id}`);
	if (!response.ok) {
		throw error(404, 'Item not found');
	}
	const item = await response.json();
	return { item };
};
```

**After (API client):**

```typescript
import { createApiClient } from '$lib/api/client';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const api = createApiClient(fetch);
	const item = await api.items.getById(params.id);
	return { item };
};
```

## Error Handling in API Client

The base `apiFetch` function automatically:

1. Sets `Content-Type: application/json`
2. Checks `response.ok`
3. Parses JSON response
4. Throws errors with message from API

```typescript
// src/lib/api/base.ts
export async function apiFetch<T>(
	fetch: typeof globalThis.fetch,
	url: string,
	init?: RequestInit
): Promise<T> {
	const response = await fetch(url, {
		...init,
		headers: {
			'Content-Type': 'application/json',
			...init?.headers
		}
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ message: 'Request failed' }));
		throw new Error(error.message || `API error: ${response.status}`);
	}

	return response.json() as Promise<T>;
}
```

## Query String Building

Use `buildQueryString` helper for type-safe URL params:

```typescript
import { buildQueryString } from '$lib/api/base';

const params = { limit: 20, offset: 0, categoryId: '123' };
const queryString = buildQueryString(params);
// Returns: "?limit=20&offset=0&categoryId=123"
```
