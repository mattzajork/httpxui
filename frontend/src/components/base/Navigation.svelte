<script>
    import { selectedTab } from '@/stores/app';
    import { loadSessionsData } from '@/stores/sessions';

    selectedTab.subscribe((tabName) => {
        handleNavClick(tabName);
    }); 

    async function handleNavClick(tabName) {
        if (tabName === 'sessions') {
            loadSessionsData();
        }
        selectedTab.set(tabName);
    }
</script>

<ul class="menu bg-pdgray">
    <li><a href="/run"
            class="{$selectedTab === 'run' ? 'bg-neutral-800' : ''}"
            on:click={(event) => {
            handleNavClick('run');
          }}>
          Run</a>
    </li>
    <li>
        <a href="/sessions"
            class="{$selectedTab === 'sessions' ? 'bg-neutral-800' : ''}"
            on:click={(event) => {
            handleNavClick('sessions');
          }}>
          Sessions</a>
    </li>
    <li><a href="/settings"
            class="{$selectedTab === 'settings' ? 'bg-neutral-800' : ''}"
            on:click={(event) => {
            handleNavClick('settings');
          }}>
          Settings</a>
    </li>
</ul>