import PocketBase from 'pocketbase';
import CommonHelper from '$lib/CommonHelper';
import { writable } from 'svelte/store';

export const httpxdata = writable([]);
export const httpxdataLoading = writable(false);

export async function loadHttpxData(sessionId) {
    httpxdataLoading.set(true);
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        
        const resultList = await pb.collection('data').getList(1, 9999999999, {
            filter: pb.filter('session = {:session}', 
                {session: sessionId})
        });

        httpxdata.set(resultList.items);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
        httpxdata.set([]);
    } finally {
        httpxdataLoading.set(false);
    }
}

export async function deleteHttpxDataBySessionId(sessionId) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);

        const resultList = await pb.collection('data').getList(1, 9999999999, {
            filter: pb.filter('session = {:session}', 
                {session: sessionId})
        });

        resultList.items.forEach((item) => {
            pb.collection('data').delete(item.id);
        });
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function saveStatusChange(data, newStatus) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        data.status = newStatus;
        return await pb.collection('data').update(data.id, data);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function saveNotes(data, newNote) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        data.notes = newNote;
        return await pb.collection('data').update(data.id, data);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function ignoreVisuallySimilar(host) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);

        const resultList = await pb.collection('data').getList(1, 9999999999, {
            filter: pb.filter('knowledgebase_phash = {:knowledgebase_phash} && session = {:session} && status = {:status}', 
                {knowledgebase_phash: host.knowledgebase_phash, session: host.session, status: ''})
        });

        resultList.items.forEach(async (i) => {
            i.status = 'Ignore';
            await pb.collection('data').update(i.id, i);
        });

        return resultList.items;
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function unignoreAll(sessionId) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);

        const resultList = await pb.collection('data').getList(1, 9999999999, {
            filter: pb.filter('session = {:session} && status = {:status}', 
                {session: sessionId, status: 'Ignore'})
        });
        
        resultList.items.forEach(async (i) => {
            i.status = '';
            await pb.collection('data').update(i.id, i);
        });

        return resultList.items;
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}