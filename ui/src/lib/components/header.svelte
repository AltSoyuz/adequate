<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';

	type Tab = { id: string; label: string; path: string };
	type Props = { tabs: Tab[] };

	let { tabs = [] }: Props = $props();

	let active = $derived(
		tabs.find((t) => page.url.pathname === t.path)?.id ?? tabs[0]?.id ?? 'home'
	);

	function navigateTo(id: string) {
		const path = tabs.find((t) => t.id === id)?.path ?? '/';
		goto(resolve(path));
	}
</script>

<!-- MOBILE: bottom bar -->
<nav aria-label="Primary" class="fixed inset-x-0 bottom-0 z-50 md:hidden">
	<div class="flex justify-center px-3 pb-3">
		<ul
			class="flex items-center justify-center gap-2 rounded-3xl border border-neutral-200 bg-white p-1"
		>
			{#each tabs as tab (tab.id)}
				<li>
					<button
						type="button"
						class="flex flex-col items-center justify-center gap-1 rounded-3xl px-4 py-3 text-neutral-600"
						class:bg-neutral-100={tab.id === active}
						class:text-neutral-900={tab.id === active}
						aria-current={tab.id === active ? 'page' : undefined}
						onclick={() => navigateTo(tab.id)}
					>
						<span class="text-sm leading-none font-medium">{tab.label}</span>
					</button>
				</li>
			{/each}
		</ul>
	</div>
</nav>

<!-- DESKTOP: top bar (same tabs, same logic) -->
<header class="sticky top-0 z-50 bg-white/80 backdrop-blur-md max-md:hidden">
	<div class="mx-auto flex h-14 max-w-4xl items-center justify-between">
		<a href={resolve('/')} class="text-sm font-semibold text-neutral-900">Adequate</a>

		<nav aria-label="Primary">
			<ul class="flex items-center gap-1 text-sm">
				{#each tabs as tab (tab.id)}
					<li>
						<button
							type="button"
							class="inline-flex items-center gap-2 rounded-lg px-3 py-2 text-neutral-700"
							class:bg-neutral-100={tab.id === active}
							class:text-neutral-900={tab.id === active}
							aria-current={tab.id === active ? 'page' : undefined}
							onclick={() => navigateTo(tab.id)}
						>
							<span class="font-medium">{tab.label}</span>
						</button>
					</li>
				{/each}
			</ul>
		</nav>
	</div>
</header>
