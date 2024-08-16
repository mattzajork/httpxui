<script>
    export let src;
    export let fullView;

    let image = {};
</script>

{#if image}
    <!-- svelte-ignore a11y-missing-attribute -->
     {#if fullView}
    <img
        src={src}
        class="border border-pdgrayoutline"
        on:error={() => image = undefined}
    >
    {:else}
    <img
        class="object-cover object-top aspect-video border border-pdgrayoutline" 
        src={src}
        on:error={() => image = undefined}
    >
    {/if}
{:else}
    <div class="border border-pdgrayoutline aspect-video">
        <svg class="max-w-full object-cover"></svg>
    </div>
{/if}