import PocketBase from 'pocketbase';
import CommonHelper from '$lib/CommonHelper';

export async function getSessionHosts(id) {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        return await pb.collection('session_hosts').getOne(id);
    } catch (error) {
        console.error('Error connecting to PocketBase:', error);
    }
}