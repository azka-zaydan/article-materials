# AI Agent Guide - [Project Name]

> **Critical**: This is the main [Project Name] application. Read ALL instructions before making code changes.

## 🚨 MANDATORY PRE-FLIGHT CHECKLIST

**Before making ANY code changes, you MUST:**

1. **Read copilot-instructions.md**: `.github/copilot-instructions.md` - Main project rules and standards
2. **Consult Svelte MCP**: Use `list-sections` → `get-documentation` → `svelte-autofixer` workflow
3. **Check relevant instructions**: Read applicable files from `.github/instructions/` based on your task
4. **Verify documentation**: Check `/docs` for existing patterns and API references
5. **Follow DRY principle**: Extract reusable functions, never duplicate code
6. **Follow design system**: ALL UI MUST match [Project Name]'s design aesthetic

## 📚 Project Documentation Structure

### Core Instructions (`.github/instructions/`)

**ALWAYS read these BEFORE coding:**

- **`api-client.instructions.md`** - Type-safe API client pattern (MANDATORY for API calls)
- **`database.instructions.md`** - Multi-schema PostgreSQL architecture (CRITICAL for DB changes)
- **`logging.instructions.md`** - Winston logging system (REQUIRED for all API/DB operations)
- **`svelte5.instructions.md`** - Svelte 5 runes patterns (NO Svelte 4 syntax allowed)
- **`validation.instructions.md`** - Zod v4 schema-first validation (MANDATORY for all inputs)
- **`ui-design.instructions.md`** - [Project Name] design system (CRITICAL for all UI work)
- **`workflows.instructions.md`** - Development workflows and quality checks

### Developer Documentation (`/docs`)

**Reference these for detailed information:**

- **`API.md`** - Complete API endpoint reference
- **`COMPONENTS.md`** - Reusable component library
- **`DATABASE.md`** - Database schema and migration system
- **`DEVELOPMENT.md`** - Development guide and patterns
- **`ERROR_HANDLING.md`** - Standardized error responses
- **`LOGGER.md`** - Winston logger usage and best practices
- **`UI_STYLE_GUIDE.md`** - Complete design system (MUST READ for UI)

## 🎯 Project-Specific Rules

### [Project Name] Design System (CRITICAL)

**EVERY UI component MUST follow these rules:**

1. **Links**: ALWAYS `[link-classes]` (e.g., `text-blue-600 hover:underline`)
2. **Borders**: ALWAYS `[border-classes]` (e.g., `border-gray-300`)
3. **Primary Branding**: `[primary-color]` for primary actions and logo ONLY
4. **Typography**: `text-sm` is default for body text (adjust as per design system)
5. **Border Radius**: Use consistent rounding (e.g., minimal `rounded` or `rounded-lg`)
6. **Backgrounds**: Clean backgrounds, avoid clutter
7. **Style Consistency**: NO deviation from the defined style guide

**Read `ui-design.instructions.md` BEFORE creating ANY UI component.**

### Database Architecture

1. **Multi-Schema**: Uses `auth`, `master`, and `app` PostgreSQL schemas (or as configured)
2. **Schema Qualification**: ALWAYS use `auth.user`, `app.[table]`, `master.[table]` in SQL
3. **Cross-Schema FKs**: Supported
4. **Custom Migrations**: SQL-based migrations (NOT Drizzle Kit)

### Critical Patterns

**API Client Usage** (MANDATORY):

```typescript
// ✅ CORRECT - Use type-safe API client
import { createApiClient } from '$api/client';

const api = createApiClient(fetch);
const data = await api.resource.getAll();

// ❌ WRONG - Never use raw fetch
const response = await fetch('/api/resource');
const data = await response.json();
```

**Validation Pattern** (MANDATORY):

```typescript
// ✅ CORRECT - Schema-first with Zod v4
import { createResourceSchema } from '$types/app.schemas';
import { ApiError } from '$server/errors';

const result = createResourceSchema.safeParse(data);
if (!result.success) {
	return ApiError.fromZod(result.error, requestId);
}
const validated = result.data; // Fully typed!

// ❌ WRONG - Manual validation
if (!data.title || data.title.length < 3) {
	return fail(400, { message: 'Invalid title' });
}
```

**Logging Pattern** (MANDATORY):

```typescript
// ✅ CORRECT - Use wrappers with context
import { withApiLogging } from '$server/logger/middleware';
import { withQueryLogging } from '$server/logger/db';

export async function POST(event) {
	return withApiLogging(
		event,
		async ({ requestId }) => {
			const data = await withQueryLogging('get_data', () => db.select().from(table), {
				requestId,
				schema: 'app'
			});
			return { data };
		},
		{ operation: 'GET_DATA', schema: 'app' }
	);
}
```

**Svelte 5 Pattern** (MANDATORY):

```svelte
<script lang="ts">
	// ✅ CORRECT - Svelte 5 runes
	let { data, children } = $props();
	let count = $state(0);
	let doubled = $derived(count * 2);

	// ❌ WRONG - Svelte 4 syntax (NEVER USE)
	// export let data;
	// let count = 0;
	// $: doubled = count * 2;
	// <slot />
</script>

<!-- ✅ CORRECT - Lowercase event handlers -->
<button onclick={increment}>Click</button>

<!-- ❌ WRONG - on: prefix -->
<button on:click={increment}>Click</button>
```

## 🔧 Available MCP Tools

### Svelte MCP Server

**ALWAYS use this workflow for Svelte development:**

#### 1. `list-sections`
- **When**: Start of ANY Svelte/SvelteKit task
- **Purpose**: Discover all available documentation sections

#### 2. `get-documentation`
- **When**: After analyzing list-sections results
- **Purpose**: Retrieve full documentation content

#### 3. `svelte-autofixer`
- **When**: AFTER writing ANY Svelte component
- **Purpose**: Validate Svelte 5 compliance and best practices
- **Requirement**: MUST return 0 issues before sending code to user

#### 4. `playground-link`
- **When**: User explicitly requests it
- **Purpose**: Generate Svelte Playground link

### Context7 MCP Server (Upstash)

**For library documentation:**

#### 1. `resolve-library-id`
- **When**: Need documentation for external library
- **Purpose**: Get Context7-compatible library ID

#### 2. `get-library-docs`
- **When**: After resolving library ID
- **Purpose**: Fetch up-to-date library documentation

## 🛠️ Development Workflow

### Before Writing Code

1. **Read Instructions**:
   - Main: `.github/copilot-instructions.md`
   - Specific: Relevant files from `.github/instructions/`
   - UI Design: `ui-design.instructions.md`
   - Docs: Check `/docs` for existing patterns

2. **Consult MCP Servers**:
   - Svelte MCP for UI work
   - Context7 MCP for library APIs

3. **Verify Patterns**:
   - API Client: Check `src/lib/api/`
   - Validation: Check `src/lib/types/*.schemas.ts`
   - Components: Check `src/lib/components/`

### During Development

1. **Create Schemas First** (Validation):
   - Define Zod schemas in `src/lib/types/*.schemas.ts`
   - Infer types with `z.infer<typeof schema>`

2. **Use API Client Pattern**:
   - Add methods to existing modules in `src/lib/api/`
   - Update `src/lib/api/client.ts`

3. **Add Logging** (MANDATORY):
   - Wrap API endpoints with `withApiLogging()`
   - Wrap DB queries with `withQueryLogging()`

4. **Follow Design System** (MANDATORY):
   - Use existing components from `src/lib/components/`
   - Import via path aliases
   - Follow `UI_STYLE_GUIDE.md`

5. **Follow Svelte 5 Patterns**:
   - Use `$props()`, `$state()`, `$derived()`
   - Use `{@render children()}`

### After Writing Code

**MANDATORY VALIDATION SEQUENCE:**

1. **Svelte Autofixer** (if Svelte component)
2. **Format Code**: `npm run format`
3. **Lint Code**: `npm run lint`
4. **Type Check**: `npm run check`

**DO NOT** consider code complete until all validation passes.

## 🚫 Anti-Patterns (NEVER DO)

### Code Anti-Patterns

❌ **Raw fetch calls** (use API client)
❌ **Missing logging** (wrap all API/DB operations)
❌ **Manual validation** (use Zod schemas)
❌ **Svelte 4 syntax** (use Svelte 5 runes)
❌ **Type duplication** (infer from schemas)
❌ **Unqualified schema names in SQL** (always use `schema.table`)

### UI Anti-Patterns

❌ **Ignoring design system**
❌ **Inconsistent styling**
❌ **Hardcoded colors/sizes** (use utility classes defined in style guide)

### Process Anti-Patterns

❌ **Skipping validation**
❌ **Ignoring MCP docs**
❌ **Not reading instructions**

## 📖 Documentation Priority

**When starting a task, read in this order:**

1. **`.github/copilot-instructions.md`**
2. **`ui-design.instructions.md`** (if UI work)
3. **Relevant `.github/instructions/*.md`**
4. **MCP Server Docs**
5. **`UI_STYLE_GUIDE.md`**
6. **`/docs` files**

**When uncertain, READ MORE DOCUMENTATION before coding.**

## ✅ Success Criteria

Code is production-ready when:

1. ✅ All relevant instructions read and followed
2. ✅ MCP servers consulted
3. ✅ Svelte autofixer returns 0 issues
4. ✅ Code formatted, linted, and type-checked
5. ✅ All patterns followed (API client, validation, logging)
6. ✅ Design system followed
7. ✅ No duplicate code
8. ✅ Documentation updated

## 🎯 Remember

- **Read first, code second**
- **Use MCP servers**
- **Follow patterns**
- **Validate everything**
- **Log comprehensively**
- **Stay DRY**
- **Type safety**
- **Svelte 5 only**
- **Strict design adherence**

**When in doubt, ask the user or read more documentation.**
