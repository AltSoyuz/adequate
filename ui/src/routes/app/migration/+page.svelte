<script lang="ts">
	import { onMount } from 'svelte';
	import { toasts } from '$lib/components/toast-state.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	onMount(() => {
		if (data.error) {
			toasts.error(data.error);
		}
	});
</script>

<h1 class="text-xl leading-tight font-semibold">Migration</h1>
{#if data.result}
	<p class="my-4 text-sm leading-relaxed text-neutral-700">
		Migration Version: {data.result.version}
	</p>
{:else}
	<p class="my-4 text-sm leading-relaxed text-neutral-700">Migration Version: unknown</p>
{/if}
