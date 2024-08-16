<script>
    import { selectedHost } from '@/stores/app';
    import { saveStatusChange, saveNotes } from '@/stores/httpxdata';
    import { BrowserOpenURL } from '$lib/wailsjs/runtime/runtime.js';
    import CommonHelper from '$lib/CommonHelper.js';
    import { DateTime } from "luxon";
    import StatusCodeBadge from '@/components/base/StatusCodeBadge.svelte';
    import ScreenshotImage from '@/components/base/ScreenshotImage.svelte';
    import Icon from '@/components/base/Icon.svelte';

    $: asnRange = $selectedHost.asn_range ? $selectedHost.asn_range : [];
    $: aRecords = $selectedHost.a ? $selectedHost.a : [];
    $: tech = $selectedHost.tech ? $selectedHost.tech : [];
    $: tlsSubjectAN = $selectedHost.tls_subject_an ? $selectedHost.tls_subject_an : [];

    let isSaving = false;
    let notes = $selectedHost.notes;

    function formatDateTime(dateTime) {
      const formatted = CommonHelper.formatToLocalDate(dateTime, "fff");
      if (formatted === 'Invalid DateTime') {
        return ''
      }
      return formatted;
    }

    function formatHeaderDateTime(dateTime) {
      return dateTime;
    }

    function formatCertificateDateTime(dateTime) {
      const formatted = DateTime.fromISO(dateTime).toLocaleString(DateTime.DATETIME_SHORT_WITH_SECONDS);
      if (formatted === 'Invalid DateTime') {
        return ''
      }
      return formatted;
    }

    function handleAddEditNotes() {
      document.getElementById('notes_modal').showModal();
    }

    async function handleSaveNote() {
      isSaving = true;
      await saveNotes($selectedHost, notes);
      isSaving = false;
      document.getElementById('notes_modal').close();
    }

    function handleStatusChange(target) {
      saveStatusChange($selectedHost, target.value);
    }
</script>

<main class="bg-black text-white">
    <div class="p-6">
        <div class="breadcrumbs text-sm mb-6">
          <ul>
            <li><a href="/sessions">Sessions</a></li>
            <li><a href="/sessions/details">Session Details</a></li>
            <li>Host Details</li>
          </ul>
        </div>
        <header class="mb-6">
          <div class="flex flex-row justify-between">
            <div>
              <div class="flex flex-row space-x-4">
                <span class="text-sm font-mono font-semibold tracking-wider truncate">{$selectedHost.method}</span>
                <div class="font-bold text-sm text-pdcontentmuted tracking-wider mb-2">
                  <a href="" class="link" on:click={BrowserOpenURL($selectedHost.url)}>{$selectedHost.url}</a>
                </div>
                <StatusCodeBadge statusCode={$selectedHost.status_code} />
              </div>
              <div class="flex flex-row">
                {#if $selectedHost.location}
                <div>
                  <Icon name="arrow-right" x="50" y="50" class="fill-white w-6 h-6 pr-1"/>
                </div>
                <div class="text-pdcontentmuted text-sm">
                  <a href="" class="link" on:click={BrowserOpenURL($selectedHost.location)}>{$selectedHost.location}</a>
                </div>
                {/if}
              </div>
            </div>
            <div>
              <div class="text-sm mb-2 text-right">
                {formatDateTime($selectedHost.timestamp)}
              </div>
              <div class="flex flex-row">
                <div class="text-xs tracking-wider text-neutral-700 uppercase mr-2 m-auto">Status</div>
                <select 
                  bind:value={$selectedHost.status}
                  on:change={(e) => handleStatusChange(e.target)}
                  class="select select-sm select-bordered w-full max-w-xs text-xs mr-4">
                  <option></option>
                  <option>Follow Up</option>
                  <option>Investigating</option>
                  <option>Complete</option>
                  <option>Ignore</option>
                </select>
                <button class="btn btn-sm btn-primary text-xs" on:click={handleAddEditNotes}>Notes</button>
              </div>
            </div>
          </div>
        </header>
        <div class="flex flex-row gap-6">
          <!-- col A -->
          <div class="basis-1/2 min-w-96">
            <!-- Overview -->
            <div class="bg-pdgray p-4 mb-6 border border-neutral-700 rounded-xl">
              <figure class="mb-4">
                <ScreenshotImage fullView={true} src={'http://127.0.0.1:8081/api/files/data/'+$selectedHost.id+'/'+$selectedHost.screenshot} />
              </figure>
              <div class="flex flex-row justify-between text-sm text-pdcontentmuted ">
                <div>
                  <div class="mb-2 select-text cursor-auto">{$selectedHost.url}</div>
                  <div class="select-text cursor-auto">{$selectedHost.title}</div>
                </div>
                <div class="font-mono text-sm select-text cursor-auto">
                  {$selectedHost.time}<span class="font-sans">&nbsp;ms</span>
                </div>
              </div>
            </div>
            <!-- page content -->
            <div class="bg-pdgray p-4 mb-6 border border-pdgrayoutline rounded-xl">
              <h2 class="mb-2 font-semibold text-lg">Response Details</h2>
              <div class="flex flex-row text-sm text-pdcontentmuted justify-between">
                <div>
                  <h3 class="text-neutral-700 mb-2">Content Length</h3>
                  <span class="select-text cursor-auto">{$selectedHost.content_length}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Words</h3>
                  <span class="text-sm select-text cursor-auto">{$selectedHost.words}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Lines</h3>
                  <span class="text-sm select-text cursor-auto">{$selectedHost.lines}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Content Type</h3>
                  <span class="font-mono select-text cursor-auto">{$selectedHost.content_type}</span>
                </div>
              </div>
            </div>
            <!--  response headers -->
            <div class="bg-pdgray p-4 mb-6 border border-pdgrayoutline rounded-xl">
              <h2 class="mb-2 font-semibold text-lg">Response Headers</h2>
              <div class="flex flex-row gap-2 text-sm text-pdcontentmuted justify-between">
                <div class="basis-1/3 truncate">
                  <h3 class="text-neutral-700 mb-2">Date</h3>
                  <span class="select-text cursor-auto">{formatHeaderDateTime($selectedHost.header_date)}</span>
                </div>
                <div class="basis-1/3 truncate">
                  <h3 class="text-neutral-700 mb-2">Etag</h3>
                  <span class="select-text cursor-auto">{$selectedHost.header_etag}</span>
                </div>
                <div class="basis-1/3 truncate">
                  <h3 class="text-neutral-700 mb-2">Last Modified</h3>
                    <span class="select-text cursor-auto">{formatHeaderDateTime($selectedHost.header_last_modified)}</span>
                </div>
              </div>
            </div>
          </div>
          <!-- col B-->
          <div class="basis-1/2 ">
            <!--  ASN -->
            <div class="bg-pdgray p-4 mb-6 border border-neutral-700 rounded-xl">
              <h2 class="mb-2 font-semibold text-lg">IP &amp; ASN</h2>
              <h3 class="text-neutral-700 text-sm mb-2">IP Addresses</h3>
              <div class="text-pdcontentmuted mb-2">
                <ul>
                  {#each aRecords as arecord }
                  <li class="font-mono text-sm select-text cursor-auto">{arecord}</li>
                  {/each}
                </ul>
              </div>
              <div class="divider"/>
              <div class="flex flex-row justify-between text-sm text-pdcontentmuted">
                <div>
                  <h3 class="text-neutral-700 mb-2">ASN</h3>
                  <span class="uppercase tracking-wider select-text cursor-auto">{$selectedHost.asn_number ? $selectedHost.asn_number : 'N/A'}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">ASN Name</h3>
                  <span class="uppercase tracking-wider select-text cursor-auto">{$selectedHost.asn_name ? $selectedHost.asn_name: 'N/A'}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">ASN Range</h3>
                  <ul class="font-mono text-sm">
                    {#each asnRange as asn_range}
                    <li class="select-text cursor-auto">{asn_range}</li>
                    {/each}
                  </ul>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Country</h3>
                  <span class="tracking-wider select-text cursor-auto">{$selectedHost.asn_country ? $selectedHost.asn_country : 'N/A'}</span>
                </div>
              </div>
            </div>
            <!-- Stack -->
            {#if tech.length > 0}
            <div class="bg-pdgray p-4 mb-6 border border-neutral-700 rounded-xl">
              <h2 class="mb-2 font-semibold text-lg">Stack</h2>
              <div class="flex flex-row justify-between text-sm text-pdcontentmuted">
                <div>
                  <h3 class="text-neutral-700 mb-2">Technologies</h3>
                  <ul class="select-text cursor-auto">
                    {#each tech as techItem}
                    <li class="uppercase tracking-wider pb-1">{techItem}</li>
                    {/each}
                  </ul>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Webserver</h3>
                  <span class="uppercase tracking-wider select-text cursor-auto">{$selectedHost.webserver}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">CDN</h3>
                  <div class="flex flex-row justify-end items-center">
                    {#if $selectedHost.cdn_name}
                    <div class="uppercase tracking-wider select-text cursor-auto">{$selectedHost.cdn_name}&nbsp;</div>
                    {/if}
                    {#if $selectedHost.cdn_type}
                    <div class="badge badge-sm badge-accent badge-outline uppercase bg-opacity-15 bg-pdyellow text-xs tracking-wider">{$selectedHost.cdn_type}</div>
                    {/if}
                  </div>
                </div>
              </div>
            </div>
            {/if}
            <!--  TLS -->
            {#if $selectedHost.tls_probe_status === true}
            <div class="bg-pdgray p-4 mb-6 border border-pdgrayoutline rounded-xl">
              <h2 class="mb-2 font-semibold text-lg">TLS</h2>
              <div class="flex flex-row gap-2 text-sm text-pdcontentmuted justify-between">
                <div>
                  <h3 class="text-neutral-700 mb-2">SNI</h3>
                  <span class="select-text cursor-auto">{$selectedHost.tls_sni}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Cipher</h3>
                  <span class="text-sm select-text cursor-auto">{$selectedHost.tls_cipher}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">TLS Version</h3>
                  <span class="uppercase tracking-wider select-text cursor-auto">{$selectedHost.tls_version}</span>
                </div>
              </div>
            </div>
            <!-- TLS Certificate -->
            <div class="bg-pdgray p-4 mb-6 border border-pdgrayoutline rounded-xl">
              <h2 class="mb-2 font-semibold text-lg">Certificate</h2>
              <div class="flex flex-row justify-between text-sm text-pdcontentmuted">
                <div>
                  <h3 class="text-neutral-700 mb-2">Not Before</h3>
                  <span class="text-sm select-text cursor-auto">{formatCertificateDateTime($selectedHost.tls_not_before)}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Not After</h3>
                  <span class="text-sm select-text cursor-auto">{formatCertificateDateTime($selectedHost.tls_not_after)}</span>
                </div>
                <div>
                  <h3 class="text-neutral-700 mb-2">Issuer</h3>
                  <span class="select-text cursor-auto">{$selectedHost.tls_issuer_org}</span>
                </div>
              </div>
            </div>
            <!-- TLS Certificate Subject AN -->
            <div class="bg-pdgray p-4 mb-6 border border-pdgrayoutline rounded-xl">
              <h2 class="mb-2 font-semibold text-lg select-none cursor-default">Certificate Subject Alternate Names</h2>
              <div class="grid gap-y-1 auto-cols-min grid-cols-2 text-sm text-pdcontentmute">
                {#each tlsSubjectAN as record}
                <div class="truncate text-pdcontentmuted select-text cursor-auto">{record}</div>
                {/each}
              </div>
            </div>
          {/if}
        </div>
     </div>
</main>

<dialog id="notes_modal" class="modal">
  <div class="modal-box border border-pdgrayoutline">
    <h2 class="mb-4 font-semibold text-lg select-none cursor-default">Notes for {$selectedHost.url}</h2>
    <div class="">
      <textarea bind:value={notes} class="textarea textarea-bordered w-full" placeholder="Notes"></textarea>
    </div>
    <div class="modal-action">
      <button on:click={() => document.getElementById('notes_modal').close()} class="btn" disabled={isSaving ? 'disabled' : ''}>Cancel</button>
      <button on:click={() => handleSaveNote()} class="btn btn-primary " disabled={isSaving ? 'disabled' : ''}>
        {#if isSaving}
        <span class="loading loading-spinner loading-xs"></span>
        {/if}
       Save 
      </button>
    </div>
  </div>
</dialog>