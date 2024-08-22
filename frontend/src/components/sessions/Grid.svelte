<script>
    import Icon from '@/components/base/Icon.svelte';
    import ScreenshotImage from '@/components/base/ScreenshotImage.svelte';
    import StatusCodeBadge from '@/components/base/StatusCodeBadge.svelte';
    import { goto } from '$app/navigation';
    import { httpxdata, ignoreVisuallySimilar, unignoreAll } from '@/stores/httpxdata.js';
    import { orderBy, filter } from 'lodash-es';
    import { removeAllToasts, addSuccessToast, addWarningToast, addErrorToast } from "@/stores/toasts";
    import { selectedHost } from '@/stores/app';
    import { selectedSession } from '@/stores/app';
    import { onMount, tick } from 'svelte';

    export let restoreScrollCallback;

    let processingIgnore = false;
    let timer = null;

    $: filter_s = '';
    $: items_s = sortFilter(filter_s, $httpxdata);

    onMount(async () => {
        await tick();
        restoreScrollCallback();
    });

    function sortFilter(filter_sp, data) {
        if (filter_sp && data) {
            let ordered = orderBy(searchObjects(filter_sp), ['knowledgebase_phash'], ['desc']);
            return filter(ordered, function(o) { return o.status !== 'Ignore'});
        } 

        let ordered = orderBy($httpxdata, ['knowledgebase_phash'], ['desc'])
        return filter(ordered, function(o) { return o.status !== 'Ignore'});
    }

    function searchObjects(filter_sp) {
        const searchableProps = ['url', 'status', 'title', 'webserver', 'method', 'status_code', 'webserver', 'content_type', 'tech'];
        return $httpxdata.filter(obj => {
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
        selectedHost.set(host);
        goto('/sessions/details/host');
    }

    async function handleUnignoreAll() {
        removeAllToasts();

        unignoreAll($selectedSession.id).then((results) => {
            addSuccessToast("Unignored all");

            let updatedData = $httpxdata.map((element) => {
                if (results.find((obj) => obj.session === $selectedSession.id)) {
                    element.status = '';
                }
                return element;
            });

            httpxdata.set(updatedData);
        });
    }

    async function handleIgnoreSimilar(host) {
        processingIgnore = true;
        removeAllToasts();
        try {
            ignoreVisuallySimilar(host, host.session).then((results) => {
                if (results.length === 0) {
                    addWarningToast("No endpoints ignored"); 
                    return;
                } else {
                    addSuccessToast("Ignored " + results.length + " visually similar endpoints");
                }

                let updatedData = $httpxdata.map((element) => {
                    if (results.find((obj) => obj.id === element.id)) {
                        element.status = 'Ignore';
                    }
                    return element;
                });

                httpxdata.set(updatedData);
            });
        } catch(error) {
            addErrorToast('An error occurred'); 
        } finally {
            processingIgnore = false;
        }
    }

    function handleKeyUp(e) {
        clearTimeout(timer);
        timer = setTimeout(() => {
            filter_s = e.target.value;
        }, 750);
    }
</script>

<div class="mt-6 mb-4">
    <div class="flex flex-row justify-between items-center">
        <div class="text-pdcontentmuted text-sm">
            <div>
                Showing <span class="font-semibold">{items_s.length}</span> of <span class="font-semibold">{$httpxdata.length}</span> endpoints
                &lpar;<a href="" on:click={() => handleUnignoreAll()} class="link link-primary">Unignore All</a>&rpar;
            </div>
        </div>
        <div>
            <div class="flex flex-row justify-end items-center">
                <div>
                    <label class="input input-bordered flex items-center gap-2">
                        <input 
                            on:keyup={handleKeyUp}
                            type="text" class="grow" style="border: none; box-shadow: none;" placeholder="Search" />
                    </label>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="flex flex-wrap gap-2 justify-between">
    {#each items_s as endpoint }
    <div class="bg-pdgray p-3 border border-pdgrayoutline rounded-xl max-w-72">
        <button class="w-full" on:click={() => handleHostClick(endpoint)}>
            <div class="tooltip" data-tip="{endpoint.url}">
            <figure class="">
                <ScreenshotImage fullView={false} src={'http://127.0.0.1:8081/api/files/data/'+endpoint.id+'/'+endpoint.screenshot} />
            </figure>
            </div>
        </button>
        <div class="text-xs truncate">{endpoint.url}</div>
        <div class="mt-1">
            <div class="flex flex-row justify-between items-center">
                <div class="text-sm font-mono tracking-wider">
                    <span class="text-pdcontentmuted">{endpoint.method}</span>&nbsp;<StatusCodeBadge statusCode={endpoint.status_code} />
                </div>
                <div>
                    <div class="tooltip tooltip-bottom" data-tip="Ignore visually similar endpoints">
                        <button on:click={() => handleIgnoreSimilar(endpoint)} class="btn btn-sm btn-ghost" disabled={processingIgnore ? 'disabled': ''}>
                            <Icon name="eye-off" class="fill-white w-5 h-5"/>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
    {/each}
</div>