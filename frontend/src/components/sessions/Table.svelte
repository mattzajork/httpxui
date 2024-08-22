<script>
    import SortHeader from '@/components/sessions/SortHeader.svelte';
    import StatusCodeBadge from '@/components/base/StatusCodeBadge.svelte';
    import { goto } from '$app/navigation';
    import { httpxdata } from '@/stores/httpxdata.js';
    import { onDestroy, onMount } from 'svelte';
    import { orderBy } from 'lodash-es';
    import { selectedHost } from '@/stores/app';
    import { sortProp, sortDir, currentPage, perPage, urlColumnEnabled, statusColumnEnabled, titleColumnEnabled, methodColumnEnabled, responseCodeColumnEnabled, webserverColumnEnabled, responseTimeColumnEnabled, contentTypeColumnEnabled, contentLengthColumnEnabled, wordsColumnEnabled, linesColumnEnabled, techColumnEnabled } from '@/stores/sessions';

    $: filter_s = '';
    $: items_s = sortFilter(filter_s, $sortProp, $sortDir);
    $: totalPages = Math.ceil(items_s.length / $perPage);
    $: totalItems = items_s.length;
    $: startIndexForPage = ($currentPage - 1) * $perPage;
    $: endIndexForPage = Math.min(startIndexForPage + $perPage, items_s.length);
    $: displayItems = items_s.slice(startIndexForPage, endIndexForPage);

    function changePageSize(e) {
        currentPage.set(1);
        perPage.set(parseInt(e.target.value));
    }

    let timer = null;

    function handleKeyUp(e) {
        clearTimeout(timer);
        timer = setTimeout(() => {
            filter_s = e.target.value;
        }, 750);
    }

    function searchObjects(filter) {
        currentPage.set(1);
        const searchableProps = ['url', 'status', 'title', 'method', 'status_code', 'time', 'webserver', 'content_type', 'content_length', 'words', 'lines', 'tech'];
        return $httpxdata.filter(obj => {
            for (const prop of searchableProps) {
                if (obj[prop] && !Array.isArray(obj[prop]) && obj[prop].toString().toLowerCase().includes(filter.toLowerCase())) {
                    return true;
                }
                if (obj[prop] && Array.isArray(obj[prop])) {
                    if (obj[prop].some(x => x.toString().toLowerCase().includes(filter.toString().toLowerCase()))) {
                        return true;
                    }
                    return false;
                }
            }
            return false;
        });
    }

    function sortFilter(filter, sortProp, sortDir) {
        let temp = [];

        if (filter) {
            temp = searchObjects(filter);
        }  else {
            temp = $httpxdata;
        }

        if (sortProp && sortDir) {
            return orderBy(temp, [sortProp], [sortDir]);
        } 

        return temp;
    }

    function incrementPage() {
        if (parseInt($currentPage) + 1 <= totalPages) {
            currentPage.set(parseInt($currentPage) + 1); 
        }
    }

    function decrementPage() {
        if (parseInt($currentPage) - 1 > 0) {
            currentPage.set(parseInt($currentPage) - 1);
        }
    }

    function techItems(techItems) {
        if (techItems !== null && techItems !== undefined) {
            return techItems;
        }
        return [];
    }

    async function handleHostClick(host) {
        selectedHost.update(() => host);
        goto('/sessions/details/host');
    }

    onDestroy(() => {
        sortProp.set('');
        sortDir.set('');
	});

    onMount(() => {
        if ($currentPage > totalPages) {
            currentPage.set(1);
            perPage.set(15);
        }
    });
</script>

<div class="mt-6 mb-4">
    <div class="flex flex-row justify-between items-center">
        <div class="dropdown">
            <div tabindex="0" role="button" class="btn btn-ghost mr-1">Columns</div>
            <!-- svelte-ignore a11y-no-noninteractive-tabindex -->
            <div tabindex="0" class="dropdown-content card card-compact bg-pdgray text-primary-content z-[1] w-64 p-2 shadow">
              <div class="card-body">
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">URL</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$urlColumnEnabled ? 'checked' : ''} on:change={(e) => urlColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Status</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$statusColumnEnabled ? 'checked' : ''} on:change={(e) => statusColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Title</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$titleColumnEnabled ? 'checked' : ''} on:change={(e) => titleColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Method</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$methodColumnEnabled ? 'checked' : ''} on:change={(e) => methodColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Response Code</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$responseCodeColumnEnabled ? 'checked' : ''} on:change={(e) => responseCodeColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Webserver</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$webserverColumnEnabled ? 'checked' : ''} on:change={(e) => webserverColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Response Time</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$responseTimeColumnEnabled ? 'checked' : ''} on:change={(e) => responseTimeColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Content Type</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$contentTypeColumnEnabled ? 'checked' : ''} on:change={(e) => contentTypeColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Content Length</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$contentLengthColumnEnabled ? 'checked' : ''} on:change={(e) => contentLengthColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Words</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$wordsColumnEnabled ? 'checked' : ''} on:change={(e) => wordsColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Lines</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$linesColumnEnabled ? 'checked' : ''} on:change={(e) => linesColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
                <div class="form-control">
                    <label class="label cursor-pointer">
                        <span class="label-text text-xs">Tech</span>
                        <input type="checkbox" class="toggle toggle-sm" checked={$techColumnEnabled ? 'checked' : ''} on:change={(e) => techColumnEnabled.set(e.target.checked)} />
                    </label>
                </div>
              </div>
            </div>
        </div>
        <div>
            <label class="input input-bordered flex items-center gap-2">
                <input 
                    on:keyup={handleKeyUp}
                    type="text" class="grow" style="border: none; box-shadow: none;" placeholder="Search" />
            </label>
        </div>
    </div>
</div>

<div class="overflow-x-scroll scrollbar-hide">
    <table class="table select-auto">
        <thead>
            <tr>
                <SortHeader class={$urlColumnEnabled ? '' : 'hidden'} name="url">URL</SortHeader>
                <SortHeader class={$statusColumnEnabled ? '' : 'hidden'} name="status">Status</SortHeader>
                <SortHeader class={$titleColumnEnabled ? '' : 'hidden'} name="title">Title</SortHeader>
                <SortHeader class={$methodColumnEnabled ? '' : 'hidden'} name="method">Method</SortHeader>
                <SortHeader class={$responseCodeColumnEnabled ? '' : 'hidden'} name="status_code">Response Code</SortHeader>
                <SortHeader class={$webserverColumnEnabled ? '' : 'hidden'} name="webserver">Webserver</SortHeader>
                <SortHeader class={$responseTimeColumnEnabled ? '' : 'hidden'} name="time">Response Time (ms)</SortHeader>
                <SortHeader class={$contentTypeColumnEnabled ? '' : 'hidden'} name="content_type">Content Type</SortHeader>
                <SortHeader class={$contentLengthColumnEnabled ? '' : 'hidden'} name="content_length">Content Length</SortHeader>
                <SortHeader class={$wordsColumnEnabled ? '' : 'hidden'} name="words">Words</SortHeader>
                <SortHeader class={$linesColumnEnabled ? '' : 'hidden'} name="lines">Lines</SortHeader>
                <th class={$techColumnEnabled ? '' : 'hidden'}>Tech</th>
            </tr>
        </thead>
        <tbody class="text-pdcontentmuted text-xs">
            {#each displayItems as host }
            <tr>
                <td class={$urlColumnEnabled ? '' : 'hidden'}>
                    <a
                        on:click|preventDefault={handleHostClick(host)}
                        class="link"
                        href=""
                    >{host.url}</a>
                </td>
                <td class={$statusColumnEnabled ? '' : 'hidden'}>
                    {#if host.status}
                    <div class="badge badge-sm badge-neutral truncate text-xs">{host.status}</div>
                    {/if}
                </td>
                <td class={$titleColumnEnabled ? '' : 'hidden'}>
                    <div class="tooltip" data-tip="{host.title}">
                        <div class="max-w-24 truncate">{host.title}</div>
                    </div>
                </td>
                <td class="{$methodColumnEnabled ? '' : 'hidden'} font-mono font-semibold">{host.method}</td>
                <td class={$responseCodeColumnEnabled ? '' : 'hidden'}><StatusCodeBadge statusCode={host.status_code} /></td>
                <td class={$webserverColumnEnabled ? '' : 'hidden'}>{host.webserver}</td>
                <td class={$responseTimeColumnEnabled ? '' : 'hidden'}>{host.time}</td>
                <td class="{$contentTypeColumnEnabled ? '' : 'hidden'} font-mono">{host.content_type}</td>
                <td class={$contentLengthColumnEnabled ? '' : 'hidden'}>{host.content_length}</td>
                <td class={$wordsColumnEnabled ? '' : 'hidden'}>{host.words}</td>
                <td class={$linesColumnEnabled ? '' : 'hidden'}>{host.lines}</td>
                <td class={$techColumnEnabled ? '' : 'hidden'}>
                    {#each techItems(host.tech) as tech}
                        <div class="badge badge-neutral my-1 mr-1 badge-sm truncate">{tech}</div>
                    {/each}
                </td>
            </tr>
            {/each}
        </tbody>
    </table>
</div>

<div class="flex flex-row justify-end space-x-12 mt-4 items-center">
    <div class="flex flex-row space-x-3 items-center">
        <div class="text-sm text-pdcontentmuted">
            Page Size
        </div>
        <div>
            <select 
                on:change={(e) => changePageSize(e)}
                class="select select-bordered select-sm text-xs w-full max-w-xs">
                <option selected={$perPage === 15 ? 'selected' : ''}>15</option>
                <option selected={$perPage === 50 ? 'selected' : ''}>50</option>
                <option selected={$perPage === 100 ? 'selected' : ''}>100</option>
            </select>
        </div>
    </div>
    <div class="text-sm text-pdcontentmuted">
        {(startIndexForPage + 1).toLocaleString()} <span class="text-neutral-600">to</span> {(endIndexForPage).toLocaleString()} <span class="text-neutral-600">of</span> {totalItems.toLocaleString()}
    </div>
    <div class="flex flex-row align-middle items-center gap-x-2 text-sm text-pdcontentmuted">
        <div>
            <button on:click={() => decrementPage()} class="btn btn-square btn-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="fill-white w-6 h-6" viewBox="25 0 512 512"><polyline points="328 112 184 256 328 400"/></svg>
            </button>
        </div>
        <div>
            Page {$currentPage} <span class="text-neutral-600">of</span> {totalPages}
        </div>
        <div>
            <button on:click={() => incrementPage()} class="btn btn-square btn-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="fill-white w-6 h-6" viewBox="-25 0 512 512"><polyline points="184 112 328 256 184 400"/></svg>
            </button>
        </div>
    </div>
</div>

<style>
/* For Webkit-based browsers (Chrome, Safari and Opera) */
.scrollbar-hide::-webkit-scrollbar {
    display: none;
}

/* For IE, Edge and Firefox */
.scrollbar-hide {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
}
</style>