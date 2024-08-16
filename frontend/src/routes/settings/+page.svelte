<script>
    import { loadSettings, saveSettings, httpProxyEnabled, httpProxyUrl, customResolvers, customHeaders } from '@/stores/settings.js';
    import { removeAllToasts, addSuccessToast, addErrorToast } from "@/stores/toasts";

    loadSettings();

    let hasChanges = false;

    // save button
    $: saveButtonEnabled = hasChanges && httpProxyAddressValid && sResolversValid && sCustomHeadersValid;

    // http proxy
    let isHttpProxyEnabled = $httpProxyEnabled;
    let httpProxyAddress = $httpProxyUrl;
    $: httpProxyAddressValid = validateProxyAddress(httpProxyAddress);

    // resolvers
    let sResolvers = $customResolvers.toString().replace(/,/g, '\n');
    $: sResolversChanged = false;
    $: sResolversValid = validateResolvers(sResolvers);

    // custom headers
    let sCustomHeaders = $customHeaders.toString().replace(/,/g, '\n');
    $: sCustomHeadersChanged = false;
    $: sCustomHeadersValid = validateCustomHeaders(sCustomHeaders);

    function validateCustomHeaders(value) {
        const pattern = /^([^:]+:\s{0,1}[^\n]+\n{0,1})*$/;
        return pattern.test(value);
    }

    function validateResolvers(value) {
        const pattern = /^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\n{0,1})*$/;
        return pattern.test(value);
    }

    function validateProxyAddress(value) {
        const pattern = /^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}\n*)+$/;
        return pattern.test(value);
    }

    function resolversChanged(value) {
        hasChanges = true;
        sResolversChanged = true;
        validateResolvers(value);
    }

    function headersChanged(value) {
        hasChanges = true;
        sCustomHeadersChanged = true;
        validateCustomHeaders(value);
    }

    function toggleHttpProxy() {
        hasChanges = true;
        isHttpProxyEnabled = !isHttpProxyEnabled;
    }

    async function saveChanges() {
        if (hasChanges) {
            try {
                let newSettings = {
                    'http_proxy': httpProxyAddress,
                    'http_proxy_enabled': isHttpProxyEnabled,
                    'resolvers': sResolvers,
                    'custom_headers': sCustomHeaders
                }
                await saveSettings(newSettings);
                httpProxyEnabled.set(isHttpProxyEnabled);
                httpProxyUrl.set(httpProxyAddress);
                if (sResolvers === "") {
                    sResolvers = [];
                }
                if (sCustomHeaders === "") {
                    sCustomHeaders = [];
                }
                let newResolvers;
                let newHeaders;
                if (Array.isArray(sResolvers) === false) {
                    newResolvers = sResolvers.split('\n');
                } else {
                    newResolvers = sResolvers;
                }
                if (Array.isArray(sCustomHeaders) === false) {
                    newHeaders = sCustomHeaders.split('\n');
                } else {
                    newHeaders = sCustomHeaders;
                }
                customResolvers.set(newResolvers);
                customHeaders.set(newHeaders);
                hasChanges = false;
                removeAllToasts();
                addSuccessToast("Settings saved")
            } catch(error) {
                removeAllToasts();                
                addErrorToast(error);
            }
        }
    }
</script>

<main class="bg-black text-white">
  <div class="p-6">
    <header>
        <div class="flex flex-row gap-6 mb-8">
            <div class="basis-1/2">
                <div class="flex flex-row justify-between">
                    <div class="text-3xl font-bold tracking-wide">
                        <span class="bg-clip-text text-transparent bg-gradient-to-r from-pdorange to-pdyellow">Settings</span>
                    </div>
                </div>
            </div>
            <div class="basis-1/2 flex justify-end">
                <button 
                    on:click={() => saveChanges()}
                    disabled="{saveButtonEnabled ? '' : 'disabled'}"
                    class="btn btn-primary">
                    Save
                </button>
            </div>
        </div>
    </header>
       <div class="flex flex-row gap-6 mb-8">
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label">
                        <span class="label-text">HTTP Proxy</span>
                    </div>
                    <div class="flex flex-row gap-4">
                        <div class="grow">
                            <input 
                                on:change={(e) => validateProxyAddress(e.target.value)}
                                bind:value={httpProxyAddress}
                                type="text" 
                                placeholder="127.0.0.1:8080" 
                                disabled="{isHttpProxyEnabled ? '' : 'disabled'}"
                                class="input input-bordered w-full {httpProxyAddressValid ? '' : 'input-error'}" />
                        </div>
                        <div class="my-auto text-right">
                            <input 
                                on:click={() => toggleHttpProxy()}
                                type="checkbox" 
                                class="toggle" 
                                checked="{isHttpProxyEnabled ? 'checked' : ''}" />
                        </div>
                    </div>
                </label>
            </div>
            <div class="basis-1/2"></div>
       </div>
       <div class="flex flex-row gap-6 mb-6">
            <div class="basis-1/2">
                <div class="basis-1/2">
                    <label class="form-control">
                        <div class="label">
                            <span class="label-text">Custom Headers</span>
                            <span class="label-text-alt text-neutral-600"><span class="font-mono">\n</span> delineated</span>
                        </div>
                        <textarea 
                            bind:value={sCustomHeaders} 
                            on:keypress={(e) => headersChanged(e.target.value)}
                            on:keydown={(e) => e.key === 'Delete' || e.key === 'Backspace' ? headersChanged(e.target.value) : null}
                            class="textarea textarea-bordered h-32 {sCustomHeadersValid ? '' : 'textarea-error'}" 
                            placeholder="X-Bugbounty: h4ck3r"></textarea>
                    </label>
                </div>
            </div>
            <div class="basis-1/2">
                <label class="form-control">
                    <div class="label">
                        <span class="label-text">Custom Resolvers</span>
                        <span class="label-text-alt text-neutral-600"><span class="font-mono">\n</span> delineated</span>
                    </div>
                    <textarea 
                        bind:value={sResolvers} 
                        on:keypress={(e) => resolversChanged(e.target.value)}
                        on:keydown={(e) => e.key === 'Delete' || e.key === 'Backspace' ? resolversChanged(e.target.value) : null}
                        class="textarea textarea-bordered h-32 {sResolversValid ? '' : 'textarea-error'}" 
                        placeholder="8.8.8.8"></textarea>
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