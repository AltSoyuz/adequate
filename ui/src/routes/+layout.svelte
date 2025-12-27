<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Toast from '$lib/components/toast.svelte';
	import BaseLayout from '$lib/layouts/BaseLayout.svelte';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} title="Adequate" /></svelte:head>
<BaseLayout>
	<main class="min-h-screen pb-20 md:pb-0">
		{@render children()}
	</main>
</BaseLayout>
<Toast />
