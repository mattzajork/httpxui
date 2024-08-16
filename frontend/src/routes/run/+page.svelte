<script>
    import { RunHTTPX } from '$lib/wailsjs/go/main/App.js';
    import { saveSessionData, runningSessions } from "@/stores/sessions";
    import { httpProxyUrl, httpProxyEnabled, customResolvers, customHeaders } from "@/stores/settings";
    import { removeAllToasts, addSuccessToast } from "@/stores/toasts";
    import RequestMethodBoxes from '@/components/run/RequestMethodBoxes.svelte';

    let targets = '';
    let paths = '';
    let httpPorts = '';
    let httpsPorts = '';
    let followRedirects = false;
    let followHostRedirects = false;
    let maxRedirects = 0;
    let delay = 0;
    let timeout = 10;
    let threads = 50;
    let requestsPerSecond = 150;
    let postBody = '';
    let allSelected = false;
    let postSelected = false;
    let getSelected = false;
    let deleteSelected = false;
    let putSelected = false;
    let patchSelected = false;
    let optionsSelected = false;
    let headSelected = false;
    let traceSelected = false;
    let connectSelected = false;

    $: hasTargets = targets.length > 0;
    $: targetsValid = validateTargets(targets);
    $: pathsValid = validatePaths(paths);
    $: httpPortsValid = validatePorts(httpPorts);
    $: httpsPortsValid = validatePorts(httpsPorts);

    function validateTargets(value) {
        const patternNewLines = /^([^\n]+\n{0,1})*$/;
        const patternCommas = /^([^,\n]+,\s{0,1})*$/;
        return patternNewLines.test(value) || patternCommas.test(value);
    }

    function validatePaths(value) {
        const patternNewLines = /^([^\n]+\n{0,1})*$/;
        const patternCommas = /^([^,\n]+,\s{0,1})*$/;
        return patternNewLines.test(value) || patternCommas.test(value);
    }

    function validatePorts(value) {
        const pattern = /^((\d{1,5}(,\s{0,1}\d{1,5})*)|)$/;
        return pattern.test(value);
    }

    function formatCommaOrNewline(items) {
        if (items.includes('\n')) {
            return items.split('\n').map(item => item.trim()).join(',');
        }  else {
            return items.split(',').map(item => item.trim()).join(',');
        }
    }

    function requestMethods() {
        let methods = '';
        if (allSelected) {
            return 'ALL';
        }
        if (getSelected) {
            methods = 'GET';
        }
        if (postSelected) {
            methods = methods + ',POST';
        }
        if (deleteSelected) {
            methods = methods + ',DELETE';
        }
        if (putSelected) {
            methods = methods + ',PUT';
        }
        if (patchSelected) {
            methods = methods + ',PATCH';
        }
        if (optionsSelected) {
            methods = methods + ',OPTIONS';
        }
        if (headSelected) {
            methods = methods + ',HEAD';
        }
        if (traceSelected) {
            methods = methods + ',TRACE';
        }
        if (connectSelected) {
            methods = methods + ',CONNECT';
        }
        return methods;
    }

    function formatCustomPorts() {
        let toJoin = [];
        if (httpPorts !== '') {
            toJoin.push('http:' + httpPorts.replace(/\s+/g, ''));
        }
        if (httpsPorts !== '') {
            toJoin.push('https:' + httpsPorts.replace(/\s+/g, ''));
        }
        if (httpPorts === '' && httpsPorts === '') {
            toJoin.push('http:' + '80');
            toJoin.push('https:' + '443');
        }
        return toJoin.join(',');
    }

    async function handleStart() {
        removeAllToasts();

        const newsession = {
            started_ts: new Date(),
            completed_ts: null
        };

        const savedsession = await saveSessionData(newsession);

        runningSessions.update((currentSessions) => [...currentSessions, savedsession]);

        RunHTTPX(savedsession.id, 
            formatCommaOrNewline(targets), 
            formatCommaOrNewline(paths), 
            formatCustomPorts(), 
            threads, 
            requestsPerSecond, 
            delay,
            requestMethods(),
            Boolean(followRedirects),
            maxRedirects,
            Boolean(followHostRedirects),
            postBody,
            Boolean($httpProxyEnabled), 
            $httpProxyUrl, 
            $customHeaders, 
            $customResolvers);
        
        addSuccessToast("Session started");
    }
</script>

<main class="bg-black text-white">
   <div class="p-6">
        <header>
            <div class="flex flex-row gap-6 mb-10">
                <div class="basis-1/2">
                    <div class="text-3xl font-bold tracking-wide">
                        <span class="bg-clip-text text-transparent bg-gradient-to-r from-pdorange to-pdyellow">Run</span>
                      </div>
                </div>
                <div class="basis-1/2 flex justify-end">
                    <button 
                        on:click={handleStart}
                        disabled="{hasTargets ? '' : 'disabled'}"
                        class="btn btn-primary">
                        Start
                    </button>
                </div>
            </div>
        </header>

        <div class="flex flex-row gap-6 mb-10">
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label">
                        <span class="label-text">Targets</span>
                        <span class="label-text-alt text-neutral-600">Comma or <span class="font-mono">\n</span> delineated</span>
                    </div>
                    <textarea
                        bind:value={targets}
                        class="textarea textarea-bordered h-24 {targetsValid ? '' : 'textarea-error'}" 
                        placeholder="www.example.com, example.com"></textarea>
                </label>
            </div>
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label">
                        <span class="label-text">Paths</span>
                        <span class="label-text-alt text-neutral-600">Comma or <span class="font-mono">\n</span> delineated</span>
                    </div>
                    <textarea 
                        bind:value={paths}
                        class="textarea textarea-bordered h-24 {pathsValid ? '' : 'textarea-error'}" 
                        placeholder="/api/v1, /robots.txt, index.jsp"></textarea>
                </label>
            </div>
        </div>

        <div class="flex flex-row gap-6 mb-10">
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label">
                        <span class="label-text">HTTP Ports</span>
                        <span class="label-text-alt text-neutral-600">Comma delineated</span>
                    </div>
                    <input 
                        bind:value={httpPorts}
                        type="text" 
                        placeholder="80, 8080" 
                        class="input input-bordered w-full {httpPortsValid ? '' : 'input-error'}" />
                </label>
            </div>
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label">
                        <span class="label-text">HTTPS Ports</span>
                        <span class="label-text-alt text-neutral-600">Comma delineated</span>
                    </div>
                    <input 
                        bind:value={httpsPorts} 
                        type="text" 
                        placeholder="443, 8443"
                        class="input input-bordered w-full {httpsPortsValid ? '' : 'input-error'}" />
                </label>
            </div>
        </div>

        <div class="flex flex-row gap-6 mb-10">
            <div class="basis-1/2">
                <RequestMethodBoxes
                    bind:allSelected={allSelected}
                    bind:getSelected={getSelected}
                    bind:postSelected={postSelected}
                    bind:putSelected={putSelected}
                    bind:deleteSelected={deleteSelected}
                    bind:patchSelected={patchSelected}
                    bind:optionsSelected={optionsSelected}
                    bind:headSelected={headSelected}
                    bind:traceSelected={traceSelected}
                    bind:connectSelected={connectSelected} />
            </div>
            <div class="basis-1/2 my-auto">
                <label class="form-control">
                    <div class="label">
                      <span class="label-text">Post Body</span>
                    </div>
                    <textarea
                        bind:value={postBody}
                        disabled={postSelected ? '' : 'disabled'}
                        class="textarea textarea-bordered" 
                        placeholder="x=1&y=2"></textarea>
                </label>
            </div>
        </div>


        <div class="flex flex-row gap-6 mb-10">
            <div class="basis-1/2">
                <div class="flex flex-row">
                    <div class="basis-1/3 mr-6 my-auto">
                        <div class="form-control">
                            <label class="label cursor-pointer">
                                <span class="label-text mr-2">Follow Redirects</span>
                                <input 
                                    bind:value={followRedirects}
                                    type="checkbox" 
                                    class="toggle" 
                                    checked="{followRedirects ? 'checked' : ''}" />
                            </label>
                        </div>
                    </div>
                    <div class="basis-2/3">
                        <label class="form-control">
                            <div class="label mb-2">
                                <span class="label-text">Maximum Redirects to Follow</span>
                                <span class="label-text-alt text-neutral-600">{maxRedirects}</span>
                            </div>
                            <input 
                                bind:value={maxRedirects} 
                                type="range" 
                                min="1" 
                                max="10" 
                                class="range range-xs range-accent" />
                        </label>
                    </div>
                </div>
            </div>
            <div class="basis-1/2">
                <div class="mr-6">
                    <div class="form-control mb-6">
                        <label class="label cursor-pointer">
                            <span class="label-text">Follow Host Redirects</span>
                            <input 
                                bind:value={followHostRedirects}
                                type="checkbox" 
                                class="toggle" 
                                checked="{followHostRedirects ? 'checked' : ''}" />
                        </label>
                    </div>
                </div>
            </div>
        </div>

        <div class="flex flex-row gap-6 mb-10">
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label mb-2">
                        <span class="label-text">Delay</span>
                        <span class="label-text-alt text-neutral-600">{delay} miliseconds between requests</span>
                    </div>
                    <input 
                        bind:value={delay} 
                        type="range" 
                        min="1" 
                        max="5000" 
                        class="range range-xs range-accent" />
                </label>
            </div>
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label mb-2">
                        <span class="label-text">Timeout</span>
                        <span class="label-text-alt text-neutral-600">{timeout} seconds</span>
                    </div>
                    <input 
                        bind:value={timeout} 
                        type="range" 
                        min="1" 
                        max="60" 
                        class="range range-xs range-accent" />
                </label>
            </div>
        </div>
        <div class="flex flex-row gap-6">
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label mb-2">
                        <span class="label-text">Threads</span>
                        <span class="label-text-alt text-neutral-600">{threads} Total Threads</span>
                    </div>
                    <input 
                        bind:value={threads} 
                        type="range" 
                        min="1" 
                        max="100" 
                        class="range range-xs range-accent" />
                </label>
            </div>
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label mb-2">
                        <span class="label-text">Maximum Requests Per Second</span>
                        <span class="label-text-alt text-neutral-600">{requestsPerSecond} RPS</span>
                    </div>
                    <input 
                        bind:value={requestsPerSecond} 
                        type="range" 
                        min="1" 
                        max="1000" 
                        class="range range-xs range-accent" />
                </label>
            </div>
        </div>
    </div>
</main>

<style>
/* Prevent resizing of textarea */
textarea {
    resize: none;
}
</style>