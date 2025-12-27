export type ToastItem = {
	id: string;
	type: 'info' | 'success' | 'error';
	title?: string;
	message: string;
};

class ToastState {
	private items: ToastItem[] = [];
	private listeners = new Set<() => void>();

	subscribe(listener: () => void) {
		this.listeners.add(listener);
		return () => this.listeners.delete(listener);
	}

	private notify() {
		this.listeners.forEach((fn) => fn());
	}

	get() {
		return this.items;
	}

	show(toast: Omit<ToastItem, 'id'>) {
		const id = crypto.randomUUID();
		this.items = [...this.items, { ...toast, id }];
		this.notify();

		setTimeout(() => this.dismiss(id), 5000);
	}

	success(message: string, title?: string) {
		this.show({ type: 'success', message, title });
	}

	error(message: string, title?: string) {
		this.show({ type: 'error', message, title });
	}

	info(message: string, title?: string) {
		this.show({ type: 'info', message, title });
	}

	dismiss(id: string) {
		this.items = this.items.filter((t) => t.id !== id);
		this.notify();
	}

	clear() {
		this.items = [];
		this.notify();
	}
}

export const toasts = new ToastState();
