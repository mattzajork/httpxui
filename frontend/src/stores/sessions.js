import PocketBase from 'pocketbase';
import { writable } from 'svelte/store';
import CommonHelper from '$lib/CommonHelper';

export const sessions = writable([]);
export const sessionsLoading = writable(false);
export const runningSessions = writable([]);
export const selectedTab = writable('grid');
export const sortProp = writable('');
export const sortDir = writable('');
export const filter = writable('');
export const currentPage = writable(1);
export const perPage = writable(15);
export const urlColumnEnabled = writable(true);
export const statusColumnEnabled = writable(true);
export const titleColumnEnabled = writable(true);
export const methodColumnEnabled = writable(true);
export const responseCodeColumnEnabled = writable(true);
export const webserverColumnEnabled = writable(true);
export const responseTimeColumnEnabled = writable(true);
export const contentTypeColumnEnabled = writable(true);
export const contentLengthColumnEnabled = writable(true);
export const wordsColumnEnabled = writable(true);
export const linesColumnEnabled = writable(true);
export const techColumnEnabled = writable(true);

export async function loadSessionsData() {
    sessionsLoading.set(true);
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        const results = await pb.collection('sessions').getFullList({
            sort: "-created"
        });
        sessions.set(results);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
        sessions.set([]);
    } finally {
        sessionsLoading.set(false);
    }
}

export async function saveSessionData(session) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        return await pb.collection('sessions').create(session);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function updateSessionData(session) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        return await pb.collection('sessions').update(session.id, session);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function deleteSession(session) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        await pb.collection('sessions').delete(session.id);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}