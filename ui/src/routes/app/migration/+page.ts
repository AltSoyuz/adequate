import { getMigrationVersion } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return getMigrationVersion(fetch);
};
