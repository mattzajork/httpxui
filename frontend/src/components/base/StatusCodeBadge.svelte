<script>
    export let statusCode;

    $: badgeStateClass = 
          statusCode.toString().startsWith('2') ? 'badge-success badge-outline bg-opacity-15 bg-success' 
        : statusCode.toString().startsWith('3') ? 'badge-info badge-outline bg-opacity-15 bg-info' 
        : statusCode.toString().startsWith('4') ? 'badge-warning badge-outline bg-opacity-15 bg-warning' 
        : statusCode.toString().startsWith('5') ? 'badge-error badge-outline bg-opacity-15 bg-error' 
        : 'badge-info';
</script>

<div class="badge gap-2 font-mono select-none {badgeStateClass}">{statusCode}</div>