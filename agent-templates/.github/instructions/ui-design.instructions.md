# UI Design & Components (CRITICAL - MANDATORY)

> **See also**: [`/docs/UI_STYLE_GUIDE.md`](../../docs/UI_STYLE_GUIDE.md) and [`/docs/COMPONENTS.md`](../../docs/COMPONENTS.md)

**Critical**: This project uses a **strict design system**. ALL UI components and pages MUST follow this design language.

## Design Philosophy

**Consistency & Minimalism** - If it doesn't serve a function, remove it.

### Core Principles

1. **Strict adherence**: Do not deviate from the defined utility classes.
2. **Simple borders**: Use `[border-color]` (e.g., `border-gray-300`) for everything.
3. **Standard links**: ALWAYS `[link-classes]` (e.g., `text-blue-600 hover:underline`).
4. **Primary branding**: `[primary-color]` for primary actions and branding.
5. **Typography**: `[default-text-size]` (e.g., `text-sm`) is the default for most content.
6. **Border Radius**: Use consistent `[radius-class]` (e.g., `rounded` or `rounded-lg`).
7. **Backgrounds**: `[bg-color]` (e.g., `bg-white`), no arbitrary colored backgrounds.

## Mandatory Styling Rules

### Links (ALWAYS)

```html
<!-- ✅ CORRECT -->
<a href="/path" class="text-blue-600 hover:underline">link text</a>

<!-- ❌ WRONG -->
<a href="/path" class="text-purple-600">link text</a>
<a href="/path" class="underline">link text</a>
```

### Borders (ALWAYS)

```html
<!-- ✅ CORRECT -->
<div class="border border-gray-300">...</div>

<!-- ❌ WRONG -->
<div class="border border-gray-400">...</div>
<div class="border-2 border-purple-500">...</div>
```

### Buttons

```html
<!-- ✅ PRIMARY -->
<button class="px-4 py-2 bg-purple-700 text-white hover:bg-purple-800 rounded text-sm">
	submit
</button>

<!-- ✅ SECONDARY -->
<button class="px-4 py-2 border border-gray-400 bg-white hover:bg-gray-100 text-gray-900 text-sm">
	cancel
</button>

<!-- ❌ WRONG -->
<button class="px-6 py-3 bg-gradient-to-r from-purple-500 rounded-xl shadow-lg">
	Fancy Button
</button>
```

### Form Inputs

```html
<!-- ✅ CORRECT -->
<input class="w-full px-3 py-1 border border-gray-300 text-sm" />

<!-- ❌ WRONG -->
<input class="rounded-lg focus:ring-4 focus:ring-purple-500 shadow-sm" />
```

### Text Sizes

```html
<!-- ✅ DEFAULT SIZES -->
text-3xl - Logo only text-2xl - Page titles text-xl - Section headers text-lg - Subsection headers
text-sm - DEFAULT for body text, buttons, inputs text-xs - Metadata, captions

<!-- ❌ WRONG -->
<p class="text-base">Normal paragraph</p>
<!-- Too large -->
<button class="text-lg">Submit</button>
<!-- Too large -->
```

## Component Usage

### ALWAYS Use These Components

Located in `src/lib/components/`:

**Common UI Components** (`ui/common/`):

- `Button.svelte` (`ui/common/actions/`) - All buttons
- `Input.svelte` (`ui/common/forms/`) - All form inputs
- `Card.svelte` (`ui/common/display/`) - Bordered containers
- `Badge.svelte` (`ui/common/display/`) - Status labels

**Layout Components** (`ui/layout/`):

- `Navigation.svelte` - Site header
- `Container.svelte` - Max-width containers
- `PageHeader.svelte` - Page titles with breadcrumbs

### Button Component

```svelte
<script>
	import Button from '$ui/common/actions/Button.svelte';
	// Or use index export:
	// import { Button } from '$components';
</script>

<!-- Primary action -->
<Button variant="primary">Submit</Button>

<!-- Secondary action -->
<Button variant="secondary">Cancel</Button>

<!-- Dangerous action -->
<Button variant="danger">Delete</Button>

<!-- Link-style button -->
<Button variant="link">Learn More</Button>
```

**Variants**:

- `primary` - Main actions
- `secondary` - Alternative actions
- `danger` - Destructive actions
- `link` - Tertiary actions

### Input Component

```svelte
<script>
	import Input from '$ui/common/forms/Input.svelte';
	// Or use index export:
	// import { Input } from '$components';

	let email = $state('');
</script>

<!-- With label -->
<Input label="Email address" type="email" bind:value={email} placeholder="you@example.com" />

<!-- With error -->
<Input label="Email address" type="email" bind:value={email} error="Please enter a valid email" />
```

### Card Component

```svelte
<script>
	import Card from '$ui/common/display/Card.svelte';
	// Or use index export:
	// import { Card } from '$components';
</script>

<!-- Default (with border) -->
<Card>
	<h2>Content</h2>
</Card>

<!-- Without border -->
<Card border={false}>
	<p>Borderless</p>
</Card>

<!-- Custom padding -->
<Card padding="lg">
	<p>Large padding</p>
</Card>
```

### Navigation Component

```svelte
<script>
	import Navigation from '$ui/layout/Navigation.svelte';
	// Or use index export:
	// import { Navigation } from '$components';
</script>

<Navigation />
```

## Common Patterns

### Page Layout

```svelte
<script>
	import Navigation from '$ui/layout/Navigation.svelte';
</script>

<div class="min-h-screen bg-white">
	<div class="max-w-6xl mx-auto px-4 py-6">
		<Navigation />

		<!-- Page content -->
		<main>
			<!-- Your content here -->
		</main>
	</div>
</div>
```

### Form Layout

```svelte
<script>
	import Input from '$ui/common/forms/Input.svelte';
	import Button from '$ui/common/actions/Button.svelte';
	// Or use index exports:
	// import { Input, Button } from '$components';

	let email = $state('');
	let password = $state('');
</script>

<div class="max-w-md mx-auto">
	<h1 class="text-2xl font-bold mb-6">Sign in to your account</h1>

	<form class="space-y-4">
		<Input label="Email address" type="email" bind:value={email} />

		<Input label="Password" type="password" bind:value={password} />

		<Button type="submit" class="w-full">Sign in</Button>
	</form>

	<p class="mt-4 text-sm text-gray-600">
		Don't have an account?
		<a href="/register" class="text-blue-600 hover:underline">Create one</a>
	</p>
</div>
```

### List Display

```svelte
<div class="space-y-2">
	{#each items as item (item.id)}
		<div class="py-2 border-b border-gray-200">
			<a href="/items/{item.id}" class="text-blue-600 hover:underline">
				{item.title}
			</a>
			<div class="text-sm text-gray-600">
				<span class="font-semibold">${item.price}</span>
				<span class="mx-2">·</span>
				<span>{item.location}</span>
			</div>
		</div>
	{/each}
</div>
```

### Breadcrumbs

```svelte
<nav class="mb-4 text-sm text-gray-600">
	<a href="/" class="text-blue-600 hover:underline">home</a>
	<span class="mx-2">›</span>
	<a href="/categories/electronics" class="text-blue-600 hover:underline"> electronics </a>
	<span class="mx-2">›</span>
	<span class="text-gray-900">smartphones</span>
</nav>
```

## Anti-Patterns (NEVER DO THIS)

### ❌ Inconsistent Styling

```html
<!-- NO arbitrary rounded corners -->
<div class="rounded-[13px]">...</div>  ❌

<!-- NO arbitrary shadows -->
<div class="shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)]">...</div>  ❌

<!-- NO arbitrary gradients -->
<div class="bg-gradient-to-r from-pink-500 to-yellow-500">...</div>  ❌
```

### ❌ Wrong Colors

```html
<!-- NO colored text except links/errors/branding -->
<p class="text-purple-600">Regular text</p>
❌

<!-- NO colored borders except errors/focus -->
<div class="border-purple-500">...</div>
❌
```

### ❌ Wrong Text Sizes

```html
<!-- NO large default text -->
<p class="text-lg">Body paragraph</p>
❌ <button class="text-base">Submit</button> ❌
```

## Writing Style

### Text Case

- **Lowercase/Sentence case**: Preferred for most UI elements
- **Title Case**: Only for proper nouns or page titles
- **UPPERCASE**: Avoid unless for specific design elements

### Placeholders

```html
<!-- ✅ CORRECT -->
<input placeholder="search [project]" />
<input placeholder="you@example.com" />

<!-- ❌ WRONG -->
<input placeholder="Search for items, categories, locations..." />
```

## Color Usage Guide

### When to Use Each Color

**Primary Color** - ONLY for:

- Primary action buttons
- Active state indicators
- Logo/Branding

**Link Color** - ALWAYS for:

- ALL hyperlinks
- Link-style buttons
- Breadcrumb links

**Error Color** - ONLY for:

- Error messages
- Error input borders
- Danger buttons

**Neutral Colors** - Everything else:

- Borders
- Text
- Backgrounds

## Validation Checklist

Before committing UI changes, verify:

- [ ] All links use standard link classes
- [ ] All borders use standard border classes
- [ ] No arbitrary shadows, gradients, or transitions (unless specified)
- [ ] Text size matches the design system
- [ ] Buttons use component variants
- [ ] Forms use Input component
- [ ] No colored backgrounds (except button states)
- [ ] Spacing is consistent

## Examples

### ✅ GOOD - Standard Style

```svelte
<div class="max-w-6xl mx-auto px-4 py-6">
	<h1 class="text-2xl font-bold mb-6">My Items</h1>

	<div class="space-y-2">
		{#each items as item (item.id)}
			<div class="py-2 border-b border-gray-200">
				<a href="/items/{item.id}" class="text-blue-600 hover:underline">
					{item.title}
				</a>
				<div class="text-sm text-gray-600">
					<span class="font-semibold">${item.price}</span>
					<span class="mx-2">·</span>
					<span>{item.location}</span>
				</div>
			</div>
		{/each}
	</div>
</div>
```

### ❌ WRONG - Inconsistent Style

```svelte
<div class="container mx-auto px-6 py-12">
	<h1
		class="text-4xl font-extrabold mb-8 bg-gradient-to-r from-purple-600 to-blue-500 bg-clip-text text-transparent"
	>
		My Items
	</h1>
    <!-- ... random styles ... -->
</div>
```

## When in Doubt

**Ask yourself**:

1. "Does this follow the strict design rules?"
2. "Is this component using the established patterns?"
3. "Can I simplify this?"

**Remember**: Consistency is key. Embrace the system.
