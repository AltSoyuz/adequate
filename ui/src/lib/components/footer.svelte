<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';

	const appVersion = $derived(page.data.appVersion);
</script>

<footer class="border-t border-neutral-200 bg-white pb-14 md:pb-0">
	<div class="px-4 py-6 md:px-6">
		<div
			class="mx-auto flex max-w-4xl flex-col gap-3 text-sm text-neutral-600 md:max-w-4xl md:flex-row md:items-center md:justify-between"
		>
			<p class="text-neutral-500">© {new Date().getFullYear()} adequate v{appVersion}</p>

			<nav aria-label="Footer" class="flex flex-wrap items-center gap-x-4 gap-y-2">
				<a href={resolve('/')} class="text-neutral-600 hover:text-neutral-900">Home</a>
				<a href={resolve('/app/migration')} class="text-neutral-600 hover:text-neutral-900"
					>Migration</a
				>
				<a href={resolve('/about')} class="text-neutral-600 hover:text-neutral-900">About</a>
			</nav>
		</div>
	</div>
</footer>
