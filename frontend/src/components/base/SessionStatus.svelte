<script>
    import { runningSessions } from '@/stores/sessions.js'
    $: numRunningSessions = $runningSessions.length;
    $: sessionSingularOrPlural = numRunningSessions > 1 ? 'sessions' : 'session';
</script>

<div class="pl-2 select-none cursor-default">
    {#if numRunningSessions > 0}
        <span class="loading loading-ring loading-xs align-middle"></span>
        <span class=" text-pdcontentmuted">{numRunningSessions} running {sessionSingularOrPlural}</span> 
    {:else}
        <span class=" text-pdcontentmuted">No running sessions</span> 
    {/if}
</div>