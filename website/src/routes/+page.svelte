<script lang="ts">
import { simStatus } from '$lib/stores/simStatus';
import { airplaneState } from '$lib/stores/airplaneState';
import { environmentState } from '$lib/stores/environmentState';
import SimulatorStatePanel from '$lib/components/SimulatorStatePanel.svelte';
</script>

<div class="min-h-full isolate relative">
{#if !$simStatus}
    <img src="https://images.unsplash.com/photo-1545972154-9bb223aac798?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=3050&q=80&exp=8&con=-15&sat=-75" alt="" class="fixed inset-0 -z-10 w-full h-full object-cover object-top" />
    <div class="w-full px-6 py-32 text-center sm:py-40 lg:px-8">
        <span class="inline-block rounded-full bg-rose-100/80 px-4 py-1 text-base font-semibold text-rose-700 shadow-md mb-4">No Connection</span>
        <h1 class="mt-4 text-5xl font-extrabold tracking-tight text-balance text-white drop-shadow-lg sm:text-7xl">Simulator Not Connected</h1>
        <p class="mt-6 text-lg font-medium text-pretty text-slate-200/90 sm:text-xl/8 drop-shadow">The application is not connected to the simulator.<br>Start the simulator and ensure SimConnect is available.</p>
    </div>
    {:else}
    <!-- Main app content goes here -->
    <section class="mx-auto max-w-3xl px-0 py-0">
      <SimulatorStatePanel />
    </section>
    <section class="mx-auto max-w-3xl px-0 py-0">
        <h2 class="text-2xl font-bold mb-4 text-blue-700 dark:text-blue-300">Airplane State</h2>
        <pre class="rounded-lg bg-gradient-to-br from-blue-900/80 to-blue-700/80 text-blue-100 p-0 overflow-x-auto shadow-lg border border-blue-400/30">{JSON.stringify($airplaneState, null, 2)}</pre>
    </section>
    <section class="mx-auto max-w-3xl px-0 py-0">
        <h2 class="text-2xl font-bold mb-4 text-emerald-700 dark:text-emerald-300">Environment State</h2>
        <pre class="rounded-lg bg-gradient-to-br from-emerald-900/80 to-emerald-700/80 text-emerald-100 p-0 overflow-x-auto shadow-lg border border-emerald-400/30">{JSON.stringify($environmentState, null, 2)}</pre>
    </section>
{/if}
</div>
