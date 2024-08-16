import { writable } from 'svelte/store';

export const targets = writable([]);
export const paths = writable([]);
export const httpPorts = writable([]);
export const httpsPorts = writable([]);
export const threads = writable(1);
export const maximumRequestsPerSecond = writable(1);
export const maximumRequestsPerMinute = writable(60);