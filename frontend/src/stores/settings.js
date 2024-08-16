import PocketBase from 'pocketbase';
import CommonHelper from '$lib/CommonHelper';
import { writable } from 'svelte/store';

export const httpProxyEnabled = writable(false);
export const httpProxyUrl = writable('');
export const customResolvers = writable([]);
export const customHeaders = writable([]);

export async function saveSettings(settings) {
    // check if a settings row exists
    const pb = new PocketBase(CommonHelper.PB_SERVER);
    let result = await pb.collection("appsettings").getFullList();
    if (result.length === 1) {
        // we have settings, so update them
        let theResult = result[0];
        theResult.http_proxy_enabled = settings.http_proxy_enabled;
        theResult.http_proxy = settings.http_proxy;
        theResult.resolvers = settings.resolvers.split('\n');
        theResult.custom_headers = settings.custom_headers.split('\n');

        if (theResult.resolvers[0] === "") {
            theResult.resolvers = [];
        }
        if (theResult.custom_headers[0] === "") {
            theResult.custom_headers = [];
        }

        return await pb.collection("appsettings").update(theResult.id, theResult);
    } else {
        // no existing settings, create them
        let newSettings = {
            'http_proxy': settings.http_proxy,
            'http_proxy_enabled': settings.http_proxy_enabled,
            'resolvers': settings.resolvers.split('\n'),
            'custom_headers': settings.custom_headers.split('\n')
        }

        if (newSettings.resolvers[0] === "") {
            newSettings.resolvers = [];
        }
        if (newSettings.custom_headers[0] === "") {
            newSettings.custom_headers = [];
        }
        return await pb.collection("appsettings").create(newSettings);
    }
}

export async function loadSettings() {
    try {
        const pb = new PocketBase(CommonHelper.PB_SERVER);
        let result = await pb.collection("appsettings").getFullList();
        if (result.length === 1) {
            httpProxyUrl.set(result[0].http_proxy);
            httpProxyEnabled.set(Boolean(result[0].http_proxy_enabled));
            customResolvers.set(result[0].resolvers);
            customHeaders.set(result[0].custom_headers);
        } 
    } catch(error) {
        console.error('Error: ', error);
    }
}