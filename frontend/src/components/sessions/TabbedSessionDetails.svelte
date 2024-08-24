<script>
    import Grid from '@/components/sessions/Grid.svelte';
    import Table from '@/components/sessions/Table.svelte';
    import { selectedTab } from '@/stores/sessions';
 
    export let restoreScrollCallback;
</script>

<div role="tablist" class="tabs tabs-bordered">
    <a 
      on:click={() => $selectedTab = 'grid'}
      href=""
      role="tab" 
      class="tab {$selectedTab === 'grid' ? 'tab-active' : ''}">Grid</a>
    <a 
      on:click={() => $selectedTab = 'table'}
      href=""
      role="tab" 
      class="tab {$selectedTab === 'table' ? 'tab-active' : ''}">Table</a>
</div>

{#if $selectedTab === 'grid'}
    <Grid restoreScrollCallback={restoreScrollCallback} />
{:else}
    <Table />
{/if}