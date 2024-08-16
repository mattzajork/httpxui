<script>
    import { loadSessionsData } from '@/stores/sessions';
    import { loadSettings } from '@/stores/settings';
    import { goto } from '$app/navigation';
    import { selectedTab } from '@/stores/app';

    loadSessionsData();
    loadSettings();

    function handleStartNewSession() {
      selectedTab.update(currentTab => 'run');
      goto('/run');
    }
</script>

<main class="bg-black text-white select-none cursor-default">
  <div class="hero bg-base-200 min-h-screen">
      <div class="hero-content text-center">
        <div class="max-w-md">
          <div class="text-5xl font-bold tracking-wide">
            <span class="bg-clip-text text-transparent bg-gradient-to-r from-pdorange to-pdyellow">Welcome</span>
          </div>
          <p class="py-6">
            An httpx-based flyover reconnaissance tool.
          </p>
          <button on:click={() => handleStartNewSession()} class="btn text-black bg-gradient-to-r from-pdorange to-pdyellow">Get Started</button>
        </div>
      </div>
    </div>
</main>