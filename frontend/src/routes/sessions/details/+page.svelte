<script>
    import CommonHelper from '$lib/CommonHelper.js';
    import TabbedSessionDetails from '@/components/sessions/TabbedSessionDetails.svelte';
    import { httpxdata, httpxdataLoading } from '@/stores/httpxdata';
    import { onMount, onDestroy, afterUpdate, tick } from 'svelte';
    import { selectedSession } from '@/stores/app';

    $: totalEndpoints = $httpxdata !== null && $httpxdata !== undefined ? $httpxdata.length : 0;

    let scrollTop = 0;

    function handleScroll() {
      scrollTop = window.scrollY;
      sessionStorage.setItem('scrollPosition', scrollTop);
    }

    function formatDateTime(dateTime) {
      const formatted = CommonHelper.formatToLocalDate(dateTime, "ff");
      if (formatted === 'Invalid DateTime') {
        return ''
      }
      return formatted;
    }

    onMount(async () => {
      window.addEventListener('scroll', handleScroll);
    });

    afterUpdate(async () => {
        const storedScrollPosition = sessionStorage.getItem('scrollPosition');
        await tick();
        window.scrollTo(0, parseInt(storedScrollPosition));
    });

    onDestroy(() => {
      window.removeEventListener('scroll', handleScroll);
    });
</script>

<main class="bg-black text-white cursor-default">
  <div class="p-6">
    <div class="breadcrumbs text-sm mb-6">
        <ul>
          <li><a href="/sessions">Sessions</a></li>
          <li>Session Details</li>
        </ul>
      </div>
    <header class="mb-12">
      <div class="flex flex-row justify-between">
        <div class=" font-thin align-middle text-sm">
          <div class="text-neutral-600 mb-2">Started</div>
          <div>{formatDateTime($selectedSession.started_ts)}</div>
        </div>
        <div class=" font-thin align-middle text-sm">
          <div class="text-neutral-600 mb-2">Completed</div>
          <div>{formatDateTime($selectedSession.completed_ts)}</div>
        </div>
        <div class=" font-thin align-middle text-sm">
          <div class="text-neutral-600 mb-2">Run Time</div>
          <div>{CommonHelper.getDurationFromDates($selectedSession.started_ts, $selectedSession.completed_ts)}</div>
        </div>
        <div class=" font-thin align-middle text-sm">
          <div class="text-neutral-600 mb-2">Endpoints</div>
          <div>{totalEndpoints}</div>
        </div>
      </div>
    </header>
      {#if $httpxdataLoading}
        <div class="hero min-h-screen">
          <div class="hero-content text-center">
            <div class="max-w-md">
              <span class="loading loading-ring loading-lg"></span>
            </div>
          </div>
        </div>
      {:else}
        {#if totalEndpoints > 0}
        <TabbedSessionDetails items={$httpxdata} />
        {:else}
        <p>No live endpoints found.</p>
        {/if}
      {/if}
  </div>
</main>