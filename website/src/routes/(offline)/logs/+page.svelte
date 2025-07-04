<script lang="ts">
import { simStatus } from '$lib/stores/simStatus';
import { airplaneState } from '$lib/stores/airplaneState';
import { environmentState } from '$lib/stores/environmentState';
import { simulatorState } from '$lib/stores/simulatorState';

console.log('Sim Status:', $simStatus);
console.log('Airplane State:', $airplaneState);
console.log('Environment State:', $environmentState);
console.log('Simulator State:', $simulatorState);
</script>

<div class="min-h-full isolate relative p-8">
    <a href="/settings" aria-label="Go to settings" class="absolute top-4 left-4 z-10 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
        <svg class="size-6" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path fill-rule="evenodd" d="M11.78 5.22a.75.75 0 0 1 0 1.06L8.06 10l3.72 3.72a.75.75 0 1 1-1.06 1.06l-4.25-4.25a.75.75 0 0 1 0-1.06l4.25-4.25a.75.75 0 0 1 1.06 0Z" clip-rule="evenodd" />
        </svg>
    </a>
    {#if !$simStatus}
        <img src="/hero-image.jpg" alt="" class="fixed inset-0 -z-10 w-full h-full object-cover object-top" />
        <div class="w-full px-6 py-32 text-center sm:py-40 lg:px-8">
            <span class="inline-block rounded-full bg-rose-100/80 px-4 py-1 text-base font-semibold text-rose-700 shadow-md mb-4">No Connection</span>
            <h1 class="mt-4 text-5xl font-extrabold tracking-tight text-balance text-white drop-shadow-lg sm:text-7xl">Simulator Not Connected</h1>
            <p class="mt-6 text-lg font-medium text-pretty text-slate-200/90 sm:text-xl/8 drop-shadow">The application is not connected to the simulator.<br>Start the simulator and ensure SimConnect is available.</p>
        </div>
    {:else}
        <section class="mx-auto max-w-3xl px-0 py-0">
            <h2 class="text-2xl font-bold mb-4 text-blue-700 dark:text-blue-300 pt-6">Sim Status</h2>
            <pre class="rounded-lg bg-gradient-to-br from-blue-900/80 to-blue-700/80 text-blue-100 p-0 overflow-x-auto shadow-lg border border-blue-400/30">{JSON.stringify($simStatus, null, 2)}</pre>
        </section>
        <section class="mx-auto max-w-3xl px-0 py-0">
            <h2 class="text-2xl font-bold mb-4 text-blue-700 dark:text-blue-300 pt-6">Airplane State</h2>
            <pre class="rounded-lg bg-gradient-to-br from-blue-900/80 to-blue-700/80 text-blue-100 p-0 overflow-x-auto shadow-lg border border-blue-400/30">{JSON.stringify($airplaneState, null, 2)}</pre>
        </section>
        <section class="mx-auto max-w-3xl px-0 py-0">
            <h2 class="text-2xl font-bold mb-4 text-emerald-700 dark:text-emerald-300 pt-6">Environment State</h2>
            <pre class="rounded-lg bg-gradient-to-br from-emerald-900/80 to-emerald-700/80 text-emerald-100 p-0 overflow-x-auto shadow-lg border border-emerald-400/30">{JSON.stringify($environmentState, null, 2)}</pre>
        </section>
        <section class="mx-auto max-w-3xl px-0 py-0">
            <h2 class="text-2xl font-bold mb-4 text-purple-700 dark:text-purple-300 pt-6">Simulator State</h2>
            <pre class="rounded-lg bg-gradient-to-br from-purple-900/80 to-purple-700/80 text-purple-100 p-0 overflow-x-auto shadow-lg border border-purple-400/30">{JSON.stringify($simulatorState, null, 2)}</pre>
        </section>
    {/if}
</div>
