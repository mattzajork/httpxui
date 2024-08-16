<script>
	import Toasts from "@/components/base/Toasts.svelte";
  import "../app.css";
  import Navigation from "@/components/base/Navigation.svelte";
  import SessionStatus from "@/components/base/SessionStatus.svelte";
  import logo from '$lib/assets/httpx-logo.png';
  import { EventsOn } from '$lib/wailsjs/runtime/runtime.js';
  import { goto } from '$app/navigation';
  import { selectedTab } from '@/stores/app.js';
  import { sessions, updateSessionData, runningSessions } from "@/stores/sessions";

  function handlePdButtonClick() {
    selectedTab.set('');
    goto('/');
  }

  // EVENT HANDLING

  EventsOn("runCompleted", function (sessionId) {
    const foundSession = $runningSessions.find((rs) => rs.id === sessionId);
    if (foundSession) {
      foundSession.completed_ts = new Date();
      updateSessionData(foundSession).then((updatedSession) => {
          sessions.update((s) => s.map((s) => s.id === sessionId ? updatedSession : s));
          runningSessions.update((s) => s.filter(s => s.id !== sessionId));
      });
    }
  });
</script>

<div class="flex min-h-screen overflow-y-auto select-none">
  <div class="flex-shrink">
    <div class="relative w-48 h-full">
      <div class="fixed h-screen flex flex-col w-48">
        <div class="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto bg-pdgray font-light">
          <Navigation />
        </div>
      </div>
    </div>
  </div>
  <div class="flex-shrink w-full overflow-x-hidden bg-black cursor-default">
    <slot/>
  </div>  
</div>
<footer class="sticky bottom-0 bg-pdgray text-white text-xs select-none">
  <Toasts />
  <div class="flex">
    <div class="flex-none bg-pdorange w-48 p-2 min-w-max flex justify-center items-center bg-opacity-75">
      <button on:click={() => handlePdButtonClick()}>
        <img alt="The project logo" class="h-3" src={logo} />
      </button>
    </div>
    <div class="flex-auto py-2 min-w-max">
      <div class="align-middle">
        <SessionStatus />
      </div>
    </div>
    <div class="flex-auto select-none cursor-default text-right p-2 min-w-max text-pdcontentmuted">By 0xm1, v1.0</div>
  </div>
</footer>