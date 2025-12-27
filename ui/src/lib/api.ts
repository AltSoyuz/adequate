const serverURL = '/api';

export type ApiResult<T> = { result: T } | { error: string; status: number };

async function readApiError(res: Response): Promise<{ error: string; status: number }> {
	const body = (await res.json()) as { error?: string };
	return {
		status: res.status,
		error: body.error ?? res.statusText
	};
}

export async function getMigrationVersion(
	fetch: typeof window.fetch
): Promise<ApiResult<{ version: string }>> {
	const res = await fetch(serverURL + '/migrations/version');

	if (!res.ok) return readApiError(res);

	const result = await res.json();
	return { result };
}
