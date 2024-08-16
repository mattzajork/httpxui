<script>
    export let getSelected = true;
    export let postSelected = false;
    export let deleteSelected = false;
    export let putSelected = false;
    export let patchSelected = false;
    export let optionsSelected = false;
    export let headSelected = false;
    export let traceSelected = false;
    export let connectSelected = false;
    export let allSelected = false;

    function selectAll() {
        getSelected = true;
        postSelected = true;
        deleteSelected = true;
        putSelected = true;
        patchSelected = true;
        optionsSelected = true;
        headSelected = true;
        traceSelected = true;
        connectSelected = true;
    }
</script>

<h2 class="uppercase tracking-wide text-sm font-semibold text-neutral-700">Request Methods</h2>
<div class="flex flex-row sm:gap-0 lg:gap-48">
    <div class="basis-1/2">
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">GET</span>
                <input bind:checked={getSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">POST</span>
                <input bind:checked={postSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">DELETE</span>
                <input bind:checked={deleteSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">PUT</span>
                <input bind:checked={putSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">PATCH</span>
                <input bind:checked={patchSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
    </div>
    <div class="basis-1/2">
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">OPTIONS</span>
                <input bind:checked={optionsSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">HEAD</span>
                <input bind:checked={headSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">TRACE</span>
                <input bind:checked={traceSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">CONNECT</span>
                <input bind:checked={connectSelected} disabled={allSelected} type="checkbox" class="checkbox checkbox-accent" />
            </label>
        </div>
        <div class="form-control">
            <label class="label cursor-pointer">
                <span class="label-text">ALL</span>
                <input 
                    bind:checked={allSelected}
                    on:change={() => selectAll()}
                    type="checkbox" 
                    class="checkbox checkbox-accent" />
            </label>
        </div>
    </div>
</div>