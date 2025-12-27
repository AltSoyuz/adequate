<script lang="ts">
	import Footer from '$lib/components/footer.svelte';
	import Header from '$lib/components/header.svelte';

	let props = $props();
	let tabs = props.tabs ?? [
		{ id: 'home', label: 'Home', path: '/' },
		{ id: 'migration', label: 'Migration', path: '/app/migration' },
		{ id: 'about', label: 'About', path: '/about' }
	];
</script>

<div class="min-h-screen bg-white text-neutral-900">
	<Header {tabs} />
	<main class="mx-auto max-w-4xl px-4 py-4 sm:px-0">
		{@render props.children?.()}
	</main>
	<Footer />
</div>
