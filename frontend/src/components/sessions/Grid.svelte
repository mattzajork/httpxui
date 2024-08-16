<script>
    import Icon from '@/components/base/Icon.svelte';
    import ScreenshotImage from '@/components/base/ScreenshotImage.svelte';
    import StatusCodeBadge from '@/components/base/StatusCodeBadge.svelte';
    import { goto } from '$app/navigation';
    import { httpxdata, ignoreVisuallySimilar, unignoreAll } from '@/stores/httpxdata.js';
    import { removeAllToasts, addSuccessToast, addWarningToast, addErrorToast } from "@/stores/toasts";
    import { selectedHost } from '@/stores/app';
    import { orderBy, filter } from 'lodash-es';

    export let items;
    let processingIgnore = false;

    $: filter_s = '';
    $: items_s = sortFilter(items, filter_s);

    function sortFilter(items_sf, filter_sp) {
        if (filter_sp) {
            let ordered = orderBy(searchObjects(items_sf, filter_sp), ['knowledgebase_phash'], ['desc']);
            return filter(ordered, function(o) { return o.status !== 'Ignore'});
        } 

        let ordered = orderBy(items_sf, ['knowledgebase_phash'], ['desc'])
        return filter(ordered, function(o) { return o.status !== 'Ignore'});
    }

    function applyFilter(term) {
        filter_s = term;
    }

    function searchObjects(data, filter_sp) {
        const searchableProps = ['url', 'status', 'title', 'webserver', 'method', 'status_code', 'webserver', 'content_type', 'tech'];
        return data.filter(obj => {
            for (const prop of searchableProps) {
                if (obj[prop] && !Array.isArray(obj[prop]) && obj[prop].toString().toLowerCase().includes(filter_sp.toLowerCase())) {
                    return true;
                }
                if (obj[prop] && Array.isArray(obj[prop])) {
                    if (obj[prop].some(x => x.toString().toLowerCase().includes(filter_sp.toString().toLowerCase()))) {
                        return true;
                    }
                    return false;
                }
            }
            return false;
        });
    }   

    async function handleHostClick(host) {
        selectedHost.update(() => host);
        goto('/sessions/details/host');
    }

    async function handleUnignoreAll() {
        removeAllToasts();

        const firstItem = items[0];
        const results = await unignoreAll(firstItem.session);

        let updatedData = $httpxdata.map((element) => {
            if (results.find((obj) => obj.session === firstItem.session)) {
                element.status = '';
            }
            return element;
        });

        httpxdata.set(updatedData);

        addSuccessToast("Unignored all");
    }

    async function handleIgnoreSimilar(host) {
        processingIgnore = true;
        removeAllToasts();
        try {
            const results = await ignoreVisuallySimilar(host, host.session);

            let updatedData = $httpxdata.map((element) => {
                if (results.find((obj) => obj.id === element.id)) {
                    element.status = 'Ignore';
                }
                return element;
            });

            httpxdata.set(updatedData);

            if (results.length === 0) {
                addWarningToast("No endpoints ignored"); 
                return;
            } else {
                addSuccessToast("Ignored " + results.length + " visually similar endpoints");
            }
        } catch(error) {
            addErrorToast('An error occurred'); 
        } finally {
            processingIgnore = false;
        }
    }
</script>

<div class="mt-6 mb-4">
    <div class="flex flex-row justify-between items-center">
        <div class="text-pdcontentmuted text-sm">
            <div>
                Showing <span class="font-semibold">{items_s.length}</span> of <span class="font-semibold">{items.length}</span> endpoints
                &lpar;<a href="" on:click={() => handleUnignoreAll()} class="link link-primary">Unignore All</a>&rpar;
            </div>
        </div>
        <div>
            <div class="flex flex-row justify-end items-center">
                <div>
                    <label class="input input-bordered flex items-center gap-2">
                        <input 
                            on:keyup={(e) => applyFilter(e.target.value)}
                            type="text" class="grow" style="border: none; box-shadow: none;" placeholder="Search" />
                    </label>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="flex flex-wrap gap-2 justify-between">
    {#each items_s as host }
    <div class="bg-pdgray p-3 border border-pdgrayoutline rounded-xl max-w-72">
        <button class="w-full" on:click={() => handleHostClick(host)}>
            <div class="tooltip" data-tip="{host.url}">
            <figure class="">
                <ScreenshotImage fullView={false} src={'http://127.0.0.1:8081/api/files/data/'+host.id+'/'+host.screenshot} />
            </figure>
            </div>
        </button>
        <div class="text-xs truncate">{host.url}</div>
        <div class="mt-1">
            <div class="flex flex-row justify-between items-center">
                <div class="text-sm font-mono tracking-wider">
                    <span class="text-pdcontentmuted">{host.method}</span>&nbsp;<StatusCodeBadge statusCode={host.status_code} />
                </div>
                <div>
                    <div class="tooltip tooltip-bottom" data-tip="Ignore visually similar endpoints">
                        <button on:click={() => handleIgnoreSimilar(host)} class="btn btn-sm btn-ghost" disabled={processingIgnore ? 'disabled': ''}>
                            {#if processingIgnore}
                            <span class="loading loading-spinner loading-xs"></span>
                            {:else}
                            <Icon name="eye-off" class="fill-white w-5 h-5"/>
                            {/if}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {/each}
</div>