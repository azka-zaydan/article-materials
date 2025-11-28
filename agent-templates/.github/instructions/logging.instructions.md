# Logging System (Winston)

> **See also**: [`/docs/LOGGER.md`](../../docs/LOGGER.md) for complete logging documentation

All HTTP requests are automatically logged via `src/hooks.server.ts` with request IDs, timing, and user context.

## Quick Reference

### API Endpoints

```typescript
import { withApiLogging } from '$lib/server/logger/middleware';

export async function GET(event) {
	return withApiLogging(
		event,
		async ({ requestId, userId }) => {
			// Custom headers: return Response directly
			return json({ data }, { headers: { 'cache-control': 'public, max-age=60' } });
		},
		{ operation: 'GET_DATA', schema: 'app' }
	);
}
```

### Database Queries

```typescript
import { withQueryLogging } from '$lib/server/logger/db';

const users = await withQueryLogging('get_users', () => db.select().from(users), {
	requestId,
	schema: 'auth',
	slowThreshold: 500
});
```

### Database Transactions

```typescript
import { logTransaction } from '$lib/server/logger/db';

const result = await logTransaction(
	db,
	'create_item_with_images',
	async (tx) => {
		const item = await tx.insert(items).values(data).returning();
		await tx.insert(images).values(imageData);
		return item[0];
	},
	{ requestId, schema: 'app' }
);
```

### Manual Logging

```typescript
import { log } from '$lib/server/logger';

log.error('Error message', error); // Errors with stack traces
log.warn('Warning', { duration: 1500 }); // Warnings
log.info('Business event', { userId, operation: 'USER_CREATED' });
log.debug('Debug info', { step: 1 }); // Filtered in production
```

## Log Levels

- `log.error()` - Errors with stack traces
- `log.warn()` - Slow queries, 4xx responses
- `log.info()` - Business events (user created, item posted)
- `log.http()` - HTTP requests/responses (auto-handled by hooks)
- `log.debug()` - Verbose debugging (filtered in production)

## Security

**Auto-redacted fields**: password, token, sessionToken, apiKey, secret

**Never log sensitive data manually.**

## Log Formats

- **Development**: Colorized console with pretty-print
- **Production**: JSON for log aggregation

## Context Requirements

**ALWAYS pass** `requestId` and `schema` to logging wrappers:

```typescript
// ✅ Good
withApiLogging(
	event,
	async ({ requestId }) => {
		/* ... */
	},
	{ operation: 'OP', schema: 'app' }
);

// ❌ Bad - missing context
withApiLogging(
	event,
	async () => {
		/* ... */
	},
	{}
);
```

## Log Files

Logs are written to:

- `logs/combined.log` - All logs
- `logs/error.log` - Error-level logs only
- Console - Development only

## Troubleshooting

### Slow Query Detection

Queries taking longer than threshold (default 200ms) are automatically logged as warnings:

```typescript
const data = await withQueryLogging('query_name', () => db.select().from(table), {
	requestId,
	schema: 'app',
	slowThreshold: 100 // Custom threshold in ms
});
```

### Request Tracing

Every request gets a unique `requestId` that appears in all logs:

```
[2025-11-25 10:30:15] info: HTTP GET /api/items [req-123456]
[2025-11-25 10:30:15] info: DB query: get_items completed in 45ms [req-123456]
[2025-11-25 10:30:15] info: HTTP GET /api/items → 200 (78ms) [req-123456]
```

Use `requestId` to trace all operations for a single request.

## See Also

- `/docs/LOGGER_SUMMARY.md` - Overview
- `/docs/LOGGER_WRAPPERS.md` - Detailed usage
- `/docs/LOGGER_FIXES.md` - Recent changes
