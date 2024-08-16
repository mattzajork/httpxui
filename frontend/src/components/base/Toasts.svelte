<script>
    import { fade, slide } from "svelte/transition";
    import { flip } from "svelte/animate";
    import { toasts, removeToast } from "@/stores/toasts";
</script>

<div class="toast toast-center">
    {#each $toasts as toast (toast.message)}
        <div
            class="alert shadow-md"
            class:alert-info={toast.type == "info"}
            class:alert-success={toast.type == "success"}
            class:alert-error={toast.type == "error"}
            class:alert-warning={toast.type == "warning"}
            in:slide={{ duration: 150 }}
            out:fade={{ duration: 150 }}
            animate:flip={{ duration: 150 }}
        >
            <div class="icon">
                {#if toast.type === "info"}
                    <i class="ri-information-line" />
                {:else if toast.type === "success"}
                    <i class="ri-checkbox-circle-line" />
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6 shrink-0 stroke-current"
                        fill="none"
                        viewBox="0 0 24 24">
                        <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                {:else if toast.type === "warning"}
                    <i class="ri-error-warning-line" />
                {:else}
                    <i class="ri-alert-line" />
                {/if}
            </div>
            <div class="content">{toast.message}</div>
            <button type="button" class="close" on:click|preventDefault={() => removeToast(toast)}>
                <i class="ri-close-line" />
            </button>
        </div>
    {/each}
</div>
