import type { LayoutServerLoad } from './$types';

export const prerender = true;
export const ssr = true;

export const load: LayoutServerLoad = async ({ fetch }) => {
	const INITIAL_APP_VERSION = '0.0.0-local-no-api';

	try {
		const res = await fetch('/api/version');
		if (!res.ok) {
			return { appVersion: INITIAL_APP_VERSION };
		}

		const version = (await res.text())?.trim() ?? INITIAL_APP_VERSION;
		return {
			appVersion: version || INITIAL_APP_VERSION
		};
	} catch {
		return { appVersion: INITIAL_APP_VERSION };
	}
};
