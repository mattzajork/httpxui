import { writable } from 'svelte/store';

export const selectedHost = writable({});
export const selectedSession = writable(null);
export const selectedTab = writable('');