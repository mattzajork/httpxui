<script>
    import Icon from '@/components/base/Icon.svelte';
    import { sortDir, sortProp } from '@/stores/sessions';

    let classes = "";
    export { classes as class }; // export reserved keyword

    const states = ['', 'asc', 'desc'];
    let currentStateIndex = 0;

    function nextState() {
        currentStateIndex = (currentStateIndex + 1) % states.length;
    }

    export let name;

    function toggleSort() {
        nextState();
        if (currentStateIndex > 0) {
            sortProp.set(name);
            sortDir.set(states[currentStateIndex]);
        } else {
            currentStateIndex = 0;
            sortProp.set('');
            sortDir.set('');
        }
    }
</script>

<th
    tabindex="0"
    title={name}
    class="cursor-pointer select-none {classes}"
    on:click={() => toggleSort()}
    on:keydown={(e) => {
        if (e.code === "Enter" || e.code === "Space") {
            e.preventDefault();
            toggleSort();
        }
    }}>
    <div class="flex flex-row">
        <div class="basis-1/2">
            <slot/>
        </div>
        <div class="basis-1/2 flex justify-end">
            {#if $sortDir === 'asc' && $sortProp === name}
            <Icon name="caret-up" x="0" y="150" class="fill-pdcontentmuted w-5 h-5 pr-1"/>
            {:else if $sortDir === 'desc' && $sortProp === name}
            <Icon name="caret-down" x="0" y="150" class="fill-pdcontentmuted w-5 h-5 pr-1"/>
            {:else}
            <Icon name="none" class="w-5 h-5 pr-1"/>
            {/if}
        </div>
    </div>
</th>