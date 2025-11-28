# Svelte 5 Patterns (CRITICAL - No Svelte 4 Syntax)

> **See also**: [Svelte 5 Official Docs](https://svelte.dev/docs/svelte/overview) and [`/docs/DEVELOPMENT.md`](../../docs/DEVELOPMENT.md)

**NEVER use Svelte 4 syntax**: No `export let`, `<slot>`, or `$:` reactive statements.

## Props

```svelte
<script lang="ts">
	// ✅ Svelte 5 runes
	let { data, children } = $props();

	// ❌ Never use Svelte 4 syntax
	// export let data;
</script>
```

## State & Reactivity

```svelte
<script lang="ts">
	let count = $state(0); // Reactive state
	let doubled = $derived(count * 2); // Computed value

	$effect(() => {
		// Side effects
		console.log('Count:', count);
	});

	// ❌ Never use $: reactive statements
	// $: doubled = count * 2;  // OLD WAY
</script>
```

## Children/Slots

```svelte
<script lang="ts">
	let { children } = $props();
</script>

<div>
	{@render children()}
	<!-- ✅ Svelte 5 -->
	<!-- <slot /> ❌ Svelte 4 -->
</div>
```

## Event Handlers

```svelte
<script lang="ts">
	let count = $state(0);

	function increment() {
		count++;
	}
</script>

<!-- ✅ Svelte 5 - lowercase event handlers -->
<button onclick={increment}>Count: {count}</button>

<!-- ❌ Svelte 4 - on:click -->
<!-- <button on:click={increment}>Count: {count}</button> -->
```

## Bindings

```svelte
<script lang="ts">
	let value = $state('');
</script>

<!-- ✅ Both work in Svelte 5 -->
<input bind:value />
<input bind:value />
```

## Conditional Rendering

```svelte
{#if condition}
	<div>Shown when true</div>
{:else if otherCondition}
	<div>Alternative</div>
{:else}
	<div>Default</div>
{/if}
```

## List Rendering

```svelte
<script lang="ts">
	let { items } = $props();
</script>

{#each items as item (item.id)}
	<div>{item.name}</div>
{/each}
```

## Snippets (Reusable Markup)

```svelte
<script lang="ts">
	let { data } = $props();
</script>

{#snippet card(title, content)}
	<div class="card">
		<h3>{title}</h3>
		<p>{content}</p>
	</div>
{/snippet}

{@render card('Title 1', 'Content 1')}
{@render card('Title 2', 'Content 2')}
```

## Component Composition

### Parent Component

```svelte
<script lang="ts">
	import Card from '$ui/common/display/Card.svelte';
	let { items } = $props();
</script>

{#each items as item (item.id)}
	<Card title={item.title}>
		<p>{item.description}</p>
	</Card>
{/each}
```

### Child Component (Card.svelte)

```svelte
<script lang="ts">
	let { title, children } = $props();
</script>

<div class="card">
	<h3>{title}</h3>
	<div class="content">
		{@render children()}
	</div>
</div>
```

## Form Handling

```svelte
<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageProps } from './$types';

	let { form }: PageProps = $props();

	// Extract field errors using $derived
	const emailErrors = $derived((form?.errors as Record<string, string[]>)?.email?.join(', '));
</script>

<form method="post" use:enhance>
	<input
		name="email"
		type="email"
		value={form?.email ?? ''}
		class={emailErrors ? 'border-red-500' : ''}
	/>
	{#if emailErrors}
		<p class="text-red-600">{emailErrors}</p>
	{/if}

	<button type="submit">Submit</button>
</form>
```

## Navigation

```svelte
<script lang="ts">
	import { goto } from '$app/navigation';

	function handleClick() {
		goto('/dashboard');
	}
</script>

<button onclick={handleClick}>Go to Dashboard</button>
```

## Type Safety

```svelte
<script lang="ts">
	import type { PageProps } from './$types';
	import type { Item } from '$lib/types/app.schemas';

	let { data }: PageProps = $props();

	// Type-safe derived values
	const items = $derived(data.items as Item[]);
</script>
```

## Common Svelte 5 Patterns

### Toggle State

```svelte
<script lang="ts">
	let isOpen = $state(false);
	const toggle = () => (isOpen = !isOpen);
</script>

<button onclick={toggle}>Toggle</button>
{#if isOpen}
	<div>Content</div>
{/if}
```

### Filtered List

```svelte
<script lang="ts">
	let { items } = $props();
	let searchQuery = $state('');

	const filteredItems = $derived(
		items.filter((item) => item.name.toLowerCase().includes(searchQuery.toLowerCase()))
	);
</script>

<input bind:value={searchQuery} placeholder="Search..." />

{#each filteredItems as item (item.id)}
	<div>{item.name}</div>
{/each}
```

### Async Data Loading

```svelte
<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	// Data is already loaded by +page.server.ts
	// No need for onMount or async logic
</script>

{#each data.items as item (item.id)}
	<div>{item.title}</div>
{/each}
```

## MCP Server Integration

**ALWAYS use Svelte MCP before writing components:**

1. Call `list-sections` → discover Svelte 5 docs
2. Call `get-documentation` → fetch API references
3. Write component using current patterns
4. Call `svelte-autofixer` → validate component
5. Only call `playground-link` if user requests it
