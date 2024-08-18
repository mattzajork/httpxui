import PocketBase from 'pocketbase';
import CommonHelper from '$lib/CommonHelper';
import { writable } from 'svelte/store';

export const httpxdata = writable([]);
export const httpxdataLoading = writable(false);

export async function loadHttpxData(sessionId) {
    httpxdataLoading.set(true);
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        
        let firstResultList = await pb.collection('data').getList(1, 500, {
            filter: pb.filter('session = {:session}', 
                {session: sessionId})
        });

        if (firstResultList.totalPages === 1) {
            httpxdata.set(firstResultList.items);
        } else {
            let currentPage = 2;
            while(currentPage <= firstResultList.totalPages) {
                let nextResultList = await pb.collection('data').getList(currentPage, 500, {
                    filter: pb.filter('session = {:session}', 
                        {session: sessionId})
                }); 
                firstResultList.items = firstResultList.items.concat(nextResultList.items);
                currentPage++;
            }
            firstResultList.perPage = firstResultList.totalItems;
            httpxdata.set(firstResultList.items);
        }
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

        let firstResultList = await pb.collection('data').getList(1, 500, {
            filter: pb.filter('session = {:session}', 
                {session: sessionId})
        });

        firstResultList.items.forEach((item) => {
            pb.collection('data').delete(item.id);
        });

        if (firstResultList.totalPages > 1) {
            let currentPage = 2;
            while(currentPage <= firstResultList.totalPages) {
                let nextResultList = await pb.collection('data').getList(currentPage, 500, {
                    filter: pb.filter('session = {:session}', 
                        {session: sessionId})
                }); 
                nextResultList.items.forEach((item) => {
                    pb.collection('data').delete(item.id);
                });
                currentPage++;
            }
        }
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

        const firstResultList = await pb.collection('data').getList(1, 500, {
            filter: pb.filter('knowledgebase_phash = {:knowledgebase_phash} && session = {:session} && status = {:status}', 
                {knowledgebase_phash: host.knowledgebase_phash, session: host.session, status: ''})
        });

        firstResultList.items.forEach(async (i) => {
            i.status = 'Ignore';
            await pb.collection('data').update(i.id, i);
        });

        if (firstResultList.totalPages > 1) {
            let currentPage = 2;
            while(currentPage <= firstResultList.totalPages) {
                let nextResultList = await pb.collection('data').getList(currentPage, 500, {
                    filter: pb.filter('knowledgebase_phash = {:knowledgebase_phash} && session = {:session} && status = {:status}', 
                        {knowledgebase_phash: host.knowledgebase_phash, session: host.session, status: ''})
                });
                nextResultList.items.forEach(async (i) => {
                    i.status = 'Ignore';
                    await pb.collection('data').update(i.id, i);
                });
                firstResultList.items = firstResultList.items.concat(nextResultList.items);
                currentPage++;
            }
            firstResultList.perPage = firstResultList.totalItems;
        }

        return firstResultList.items;
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}

export async function unignoreAll(sessionId) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);

        const firstResultList = await pb.collection('data').getList(1, 500, {
            filter: pb.filter('session = {:session} && status = {:status}', 
                {session: sessionId, status: 'Ignore'})
        });

        firstResultList.items.forEach(async (i) => {
            i.status = '';
            await pb.collection('data').update(i.id, i);
        });

        if (firstResultList.totalPages > 1) {
            let currentPage = 2;
            while(currentPage <= firstResultList.totalPages) {
                let nextResultList = await pb.collection('data').getList(currentPage, 500, {
                    filter: pb.filter('session = {:session} && status = {:status}', 
                        {session: sessionId, status: 'Ignore'})
                });

                nextResultList.items.forEach(async (i) => {
                    i.status = '';
                    await pb.collection('data').update(i.id, i);
                });

                firstResultList.items = firstResultList.items.concat(nextResultList.items);
                currentPage++;
            }
            firstResultList.perPage = firstResultList.totalItems;
        }

        return firstResultList.items;
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}