# Zod v4 Validation System (CRITICAL - MANDATORY)

> **See also**: [`/docs/ERROR_HANDLING.md`](../../docs/ERROR_HANDLING.md) for complete error handling reference

**Critical**: This project uses Zod v4 for runtime type-safe validation. **ALWAYS create Zod schemas for database tables FIRST, then infer TypeScript types from them.**

## Schema-First Development Workflow

1. **ALWAYS create Zod objects** in `src/lib/types/*.schemas.ts` (tables, request bodies, and query schemas). Do not define schemas inline in routes or components.
2. **ALWAYS use `z.infer<>` to create types** - never duplicate type definitions
3. **ALWAYS validate in two places**:
   - **Client/Frontend**: Validate before sending to action or fetching API
   - **Backend**: Validate received request and type-cast using inferred types

## Zod v4 Syntax (IMPORTANT)

**Use `error` parameter, NOT `message`** - The `message` parameter is deprecated in Zod v4:

```typescript
// ✅ CORRECT - Zod v4
z.string().min(8, { error: 'Too short' });
z.email({ error: 'Invalid email' });
z.uuid('Invalid UUID'); // Can pass string directly for single error

// ❌ WRONG - Deprecated syntax
z.string().min(8, 'Too short'); // Will show deprecation warning
z.string().min(8, { message: 'Too short' }); // Deprecated parameter
z.string().uuid({ message: 'Invalid' }); // Use z.uuid() directly
```

## Common Validation Schemas

Located in `src/lib/types/common.schemas.ts`:

```typescript
import {
	uuidSchema, // UUID v4 validation
	emailSchema, // Email with RFC 5322 compliance
	phoneNumberSchema, // US phone numbers
	usernameSchema, // 3-50 chars, alphanumeric
	passwordSchema, // Strong password requirements
	slugSchema, // URL-friendly slugs
	priceSchema, // Decimal with 2 places
	limitSchema, // Pagination limit (1-100)
	offsetSchema, // Pagination offset
	searchQuerySchema, // Search query string
	sortOrderSchema, // 'asc' | 'desc'
	imageUrlSchema // Image URL validation
} from '$lib/types/common.schemas';
```

## Entity Validation Schemas

Schema files mirror database structure:

- `src/lib/types/auth.schemas.ts` - User, Session, etc.
- `src/lib/types/app.schemas.ts` - Main application entities (e.g., Items, Posts)
- `src/lib/types/master.schemas.ts` - Reference data (e.g., Categories, Locations)
- `src/lib/types/common.schemas.ts` - Reusable primitives and helpers

Each entity has multiple schema variants:

- **Base schema**: Full database record (`itemSchema`, `userSchema`)
- **Create schema**: For creating new records (`createItemSchema`, `registerUserSchema`)
- **Update schema**: Partial version for updates (`updateItemSchema`)
- **Response schema**: What API returns (`userResponseSchema`, `itemDetailResponseSchema`)
- **Query schema**: URL params/query strings (`itemQuerySchema`)
- **API response schemas**: Match actual endpoint responses (`apiItemsResponseSchema`)

## Type Inference (MANDATORY PATTERN)

**ALWAYS use `z.infer<>` for TypeScript types** from Zod schemas:

```typescript
// ✅ CORRECT - Single source of truth
export const itemSchema = z.object({
	id: uuidSchema,
	title: z.string().min(1).max(255),
	price: priceSchema
	// ...
});

// Infer types from schemas
export type Item = z.infer<typeof itemSchema>;
export type CreateItem = z.infer<typeof createItemSchema>;
export type UpdateItem = z.infer<typeof updateItemSchema>;

// ❌ WRONG - Don't duplicate type definitions manually
// type Item = { id: string; title: string; ... }
```

## Client-Side Validation (Frontend)

**ALWAYS validate before sending data to actions or API endpoints:**

```typescript
// In Svelte component or +page.server.ts
import { createItemSchema } from '$lib/types/app.schemas';
import { createApiClient } from '$lib/api/client';

// Validate before API call
const api = createApiClient(fetch);
const result = createItemSchema.safeParse(formData);

if (!result.success) {
	// Show validation errors to user
	const errors: Record<string, string[]> = {};
	result.error.issues.forEach((err) => {
		const path = err.path[0] as string;
		if (!errors[path]) errors[path] = [];
		errors[path].push(err.message);
	});
	return fail(400, { message: 'Validation failed', errors });
}

// Call API with validated data (already typed!)
await api.items.create(result.data);
```

## Backend Validation (API Routes)

**ALWAYS validate received requests and type-cast:**

```typescript
import { createItemSchema } from '$lib/types/app.schemas';
import { ApiError, ApiSuccess } from '$lib/server/errors';
import { withApiLogging } from '$lib/server/logger/middleware';

export async function POST(event) {
	return withApiLogging(
		event,
		async ({ requestId }) => {
			// Parse and validate request body
			const body = await event.request.json();

			if (!body) {
				return ApiError.badRequest('Invalid JSON in request body', undefined, requestId);
			}

			const result = createItemSchema.safeParse(body);

			if (!result.success) {
				// Automatically converts ZodError to ApiErrorResponse with field errors
				return ApiError.fromZod(result.error, requestId);
			}

			// result.data is now fully typed as CreateItem!
			const validatedData = result.data;

			// Type-safe database operations
			const [item] = await db
				.insert(items)
				.values({
					...validatedData,
					userId: locals.user.id
				})
				.returning();

			return ApiSuccess.created(item, { requestId });
		},
		{ operation: 'CREATE_ITEM', schema: 'app' }
	);
}
```

## Form Actions Validation Pattern

```typescript
// src/routes/register/+page.server.ts
import { registerUserSchema } from '$lib/types/auth.schemas';
import { fail, redirect } from '@sveltejs/kit';

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const data = Object.fromEntries(formData);

		// Validate using Zod schema
		const result = registerUserSchema.safeParse(data);

		if (!result.success) {
			// Extract field-level errors
			const errors: Record<string, string[]> = {};
			result.error.issues.forEach((err) => {
				const path = err.path[0] as string;
				if (!errors[path]) errors[path] = [];
				errors[path].push(err.message);
			});

			return fail(400, {
				message: 'Please fix the validation errors',
				...data, // Return form data for repopulation
				errors
			});
		}

		// Proceed with validated data (typed!)
		await api.auth.register(result.data);
		return redirect(302, '/dashboard');
	}
};
```

## Validation Testing Patterns

### Unit Testing Schemas

```typescript
import { describe, it, expect } from 'vitest';
import { createItemSchema } from '$lib/types/app.schemas';

describe('createItemSchema', () => {
	it('validates correct item data', () => {
		const validData = {
			title: 'Test Item',
			description: 'A valid description here',
			categoryId: '123e4567-e89b-12d3-a456-426614174000',
			locationId: '123e4567-e89b-12d3-a456-426614174001',
			price: '99.99'
		};

		const result = createItemSchema.safeParse(validData);
		expect(result.success).toBe(true);
	});

	it('rejects invalid data', () => {
		const invalidData = {
			title: 'ab', // Too short
			description: 'short',
			categoryId: 'not-a-uuid',
			locationId: 'not-a-uuid',
			price: '-10.00' // Negative price
		};

		const result = createItemSchema.safeParse(invalidData);
		expect(result.success).toBe(false);
		if (!result.success) {
			expect(result.error.issues.length).toBeGreaterThan(0);
		}
	});
});
```
