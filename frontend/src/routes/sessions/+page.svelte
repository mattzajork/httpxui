<script>
  import CommonHelper from '$lib/CommonHelper.js';
  import { getSessionHosts } from '@/stores/views';
  import { goto } from '$app/navigation';
  import { loadHttpxData, deleteHttpxDataBySessionId } from '@/stores/httpxdata';
  import { updateSessionData, deleteSession, sessions, sessionsLoading, runningSessions } from '@/stores/sessions';
  import { removeAllToasts, addSuccessToast } from "@/stores/toasts";
  import { selectedSession } from '@/stores/app';
  import { selectedTab } from '@/stores/app';

  let selectedForDelete = {};
  let isDeleting = false;

  $: sessionLength = $sessions ? $sessions.length : 0;

  function formatDateTime(dateTime) {
    const formatted = CommonHelper.formatToLocalDate(dateTime, "fff");
    if (formatted === 'Invalid DateTime') {
      return ''
    }
    return formatted;
  }

  sessions.subscribe((s) => {
    let unfinishedSessions = $sessions.filter((session) => session.completed_ts === '');
    unfinishedSessions.forEach(async function(item) {
        if ($runningSessions.find(obj => obj.id === item.id) === undefined) {
          item.completed_ts = new Date();
          await updateSessionData(item);
          sessions.update((ss) => ss.map((x) => x.id === item.id ? {...x, completed_ts: item.completed_ts} : x));
        }
      });
  });

  async function handleViewSession(session) {
    selectedSession.update(() => session);
    loadHttpxData(session.id);
    goto('/sessions/details');
  }

  async function deleteSelected() {
    isDeleting = true;
    removeAllToasts();
    const sessionId = selectedForDelete.id;
    await deleteHttpxDataBySessionId(sessionId);
    await deleteSession(selectedForDelete);

    sessions.update((sessions) => sessions.filter(s => s.id !== selectedForDelete.id));
    selectedForDelete = {};
    document.getElementById('delete_modal').close();
    addSuccessToast("Successfully deleted session");
    isDeleting = false;
  }

  async function handleDeleteSession(session) {
    selectedForDelete = session;
    document.getElementById('delete_modal').showModal();
  }

  function handleStartNewSession() {
    selectedTab.update(() => 'run');
    goto('/run');
  }
</script>

{#if $sessionsLoading === false && sessionLength > 0}
<main class="bg-black text-white">
  <div class="p-6">
    <header class="mb-6">
      <div class="text-3xl font-bold tracking-wide">
        <span class="bg-clip-text text-transparent bg-gradient-to-r from-pdorange to-pdyellow">Sessions</span>
      </div>
    </header>
    <div class="divide-y divide-base-200 cursor-default">
      {#each $sessions as session}
      <div class="flex flex-row p-6 align-middle">
        <div class="basis-1/4 font-thin align-middle text-sm mr-4">
          <div class="text-neutral-600 mb-2 mr-2">Started</div>
          <div>{formatDateTime(session.started_ts)}</div>
        </div>
        <div class="basis-1/4 font-thin align-middle text-sm my-auto mr-4">
          {#if session.completed_ts !== ''}
            <div class="text-neutral-600 mb-2">Completed</div>
            <div>{formatDateTime(session.completed_ts)}</div>
          {:else}
            <div class="text-pdcontentmuted">
              <span class="loading loading-ring loading-sm align-middle mr-1"></span>
              <span class="uppercase tracking-widest text-xs">In Progress</span>
            </div>
          {/if}
        </div>
        <div class="basis-1/4 font-thin align-middle text-sm">
          {#if session.completed_ts !== ''}
          <div class="text-neutral-600 mb-2">Endpoints</div>
          {#await getSessionHosts(session.id)}
          <span class="loading loading-spinner loading-xs text-pdcontentmuted"></span>
          {:then result} 
          <div>{result.total_hosts}</div>
          {/await}
          {/if}
        </div>
        <div class="basis-1/4 m-auto">
          <div class="join">
            <button class="join-item btn btn-primary btn-sm text-xs min-w-20" disabled={session.completed_ts ? '' : 'disabled'} on:click={handleViewSession(session)}>View</button>
            <button class="join-item btn btn-error btn-sm text-xs min-w-20" disabled={session.completed_ts ? '' : 'disabled'} on:click={handleDeleteSession(session)}>Delete</button>
          </div>
        </div>
      </div>
      {/each}
    </div>
  </div>
</main>
{:else if $sessionsLoading === true && sessionLength === 0}
<div class="hero min-h-screen">
  <div class="hero-content text-center">
    <div class="max-w-md">
      <span class="loading loading-ring loading-lg"></span>
    </div>
  </div>
</div>
{:else}
<main class="bg-black text-white">
    <div class="divide-y divide-base-200">
      <div class="hero bg-base-200 min-h-screen">
        <div class="hero-content text-center">
          <div class="max-w-md">
            <h1 class="text-5xl font-bold">No Sessions</h1>
            <p class="py-6">
              Start a Run to create a new session.
            </p>
            <button on:click={() => handleStartNewSession()} class="btn btn-primary">Start a new session</button>
          </div>
        </div>
      </div>
    </div>
</main>
{/if}

<dialog id="delete_modal" class="modal">
  <div class="modal-box border border-pdgrayoutline">
    <p class="py-4">Are you sure you want to delete this session?</p>
    <div class="modal-action">
      <button on:click={() => document.getElementById('delete_modal').close()} class="btn">Cancel</button>
      <button on:click={() => deleteSelected()} class="btn btn-error" disabled={isDeleting ? 'disabled' : ''}>
        {#if isDeleting}
        <span class="loading loading-spinner loading-xs"></span>
        {/if}
        Delete
      </button>
    </div>
  </div>
</dialog>