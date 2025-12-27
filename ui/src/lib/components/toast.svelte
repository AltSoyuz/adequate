<script lang="ts">
	import { toasts } from '$lib/components/toast-state.svelte';
	import { fly } from 'svelte/transition';
	import { flip } from 'svelte/animate';
	import { CircleX, X, Check } from 'lucide-svelte';

	let items = $state(toasts.get());

	toasts.subscribe(() => {
		items = toasts.get();
	});
</script>

<div
	class="pointer-events-none fixed inset-0 z-50 flex flex-col items-center px-4 py-4 pt-16 sm:pt-4"
	aria-live="polite"
	aria-relevant="additions"
>
	<div class="flex w-full flex-col items-center space-y-3">
		{#each items as t (t.id)}
			{@const bg =
				t.type === 'success' ? 'bg-emerald-50' : t.type === 'error' ? 'bg-red-50' : 'bg-slate-50'}
			{@const iconColor =
				t.type === 'success'
					? 'text-emerald-600'
					: t.type === 'error'
						? 'text-rose-600'
						: 'text-slate-600'}
			{@const IconComponent = t.type === 'success' ? Check : t.type === 'error' ? CircleX : X}

			<div
				animate:flip={{ duration: 300 }}
				transition:fly={{ y: -100, duration: 250 }}
				class="pointer-events-auto w-full max-w-sm overflow-hidden rounded"
			>
				<div class={'border border-slate-200 shadow-lg ' + bg}>
					<div class="p-4">
						<div class="flex items-start gap-3">
							<IconComponent class={'h-5 w-5 ' + iconColor} />

							<div class="min-w-0 flex-1">
								{#if t.title}
									<p class="text-sm font-semibold text-slate-900">{t.title}</p>
								{/if}
								<p class={'text-sm text-slate-700 ' + (t.title ? 'mt-0.5' : '')}>{t.message}</p>
							</div>

							<button
								type="button"
								class="flex-none rounded p-1 transition-colors hover:bg-black/10 focus:ring-2 focus:ring-slate-400 focus:ring-offset-2 focus:outline-none"
								onclick={() => toasts.dismiss(t.id)}
								aria-label="Dismiss"
							>
								<X class="h-4 w-4" />
							</button>
						</div>
					</div>
				</div>
			</div>
		{/each}
	</div>
</div>
