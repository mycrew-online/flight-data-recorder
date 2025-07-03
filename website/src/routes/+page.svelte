<script lang="ts">
import { simStatus } from '$lib/stores/simStatus';
import { airplaneState } from '$lib/stores/airplaneState';
import { environmentState } from '$lib/stores/environmentState';
import { recordingState } from '$lib/stores/recordingState';
import WeatherPanel from '$lib/components/WeatherPanel.svelte';
import AircraftPanel from '$lib/components/AircraftPanel.svelte';

// Helper: Convert sim_time (seconds) to HH:MM:SS
function formatSimTime(simTime: number | undefined | null): string {
  if (typeof simTime !== 'number' || isNaN(simTime)) return '';
  const hours = Math.floor(simTime / 3600);
  const minutes = Math.floor((simTime % 3600) / 60);
  const seconds = Math.floor(simTime % 60);
  return `${hours.toString().padStart(2, '0')}:${minutes
    .toString()
    .padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}

// Helper: Convert local_time (seconds since midnight) to HH:MM:SS
function formatLocalTime(localTime: number | undefined | null): string {
  if (typeof localTime !== 'number' || isNaN(localTime)) return '';
  const hours = Math.floor(localTime / 3600);
  const minutes = Math.floor((localTime % 3600) / 60);
  const seconds = Math.floor(localTime % 60);
  return `${hours.toString().padStart(2, '0')}:${minutes
    .toString()
    .padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}

// Helper: Format local date as DD.MM.YYYY
function formatLocalDate(year: number | undefined, month: number | undefined, day: number | undefined): string {
  if (!year || !month || !day) return '';
  return `${day.toString().padStart(2, '0')}.${month
    .toString()
    .padStart(2, '0')}.${year.toString().padStart(4, '0')}`;
}
</script>

<div class="min-h-full isolate relative p-6">
{#if !$simStatus}
    <img src="/hero-image.jpg" alt="" class="fixed inset-0 -z-10 w-full h-full object-cover object-top" />
    <div class="w-full px-6 py-32 text-center sm:py-40 lg:px-8">
        <span class="inline-block rounded-full bg-rose-100/80 px-4 py-1 text-base font-semibold text-rose-700 shadow-md mb-4">Is your simulator running?</span>
        <h1 class="mt-4 text-5xl font-extrabold tracking-tight text-balance text-white drop-shadow-lg sm:text-7xl">Not Connected</h1>
        <p class="mt-6 text-lg font-medium text-pretty text-slate-200/90 sm:text-xl/8 drop-shadow">The application is not connected to the simulator.<br>Start the simulator and ensure SimConnect is available.</p>
    </div>
    {:else}
    <div class="lg:flex lg:items-center lg:justify-between">
        <div class="min-w-0 flex-1">
            <h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">{$airplaneState?.title}</h2>
            <div class="mt-1 flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap sm:space-x-6">
            <!-- Combined Sim Time, Date, and Local Time -->
            <div class="mt-2 flex items-center text-sm text-gray-500">
                <svg class="mr-1.5 size-5 shrink-0 text-gray-400" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <title>Simulated time (in-game)</title>
                  <path fill-rule="evenodd" d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm0-14.5A6.5 6.5 0 1 1 3.5 10 6.508 6.508 0 0 1 10 3.5Zm.75 3.75a.75.75 0 0 0-1.5 0v3.25a.75.75 0 0 0 .22.53l2.25 2.25a.75.75 0 1 0 1.06-1.06l-2.03-2.03V7.25Z" clip-rule="evenodd" />
                </svg>
                {formatSimTime($environmentState?.sim_time)}
                <svg class="ml-4 mr-1.5 size-5 shrink-0 text-gray-400" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <title>Local date (simulator)</title>
                  <path fill-rule="evenodd" d="M6 2a1 1 0 0 0-1 1v1H5A3 3 0 0 0 2 7v7a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3V7a3 3 0 0 0-3-3h-.002V3a1 1 0 1 0-2 0v1H7V3a1 1 0 0 0-1-1Zm10 5v1H4V7a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1Zm-1 3v4a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1v-4h10Z" clip-rule="evenodd" />
                </svg>
                {formatLocalDate($environmentState?.local_year, $environmentState?.local_month, $environmentState?.local_day)}
                <span>&nbsp;{formatLocalTime($environmentState?.local_time)}</span>
            </div>

            </div>
        </div>
        <div class="mt-5 flex lg:mt-0 lg:ml-4">
            <span class="hidden sm:block">
            <button
                type="button"
                class="inline-flex items-center rounded-md px-3 py-2 text-sm font-semibold shadow-xs ring-1 ring-inset focus:outline-none transition-colors duration-150
                    bg-green-600 text-white hover:bg-green-700 ring-green-500"
                aria-pressed={$recordingState === "recording" ? 'true' : 'false'}
    
            >
                {#if $recordingState === "recording"}
                    <svg class="mr-1.5 -ml-0.5 size-5 text-red-400 animate-pulse" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <title>Stop Recording</title>
                        <rect x="6" y="6" width="8" height="8" rx="2" />
                    </svg>
                    Stop Recording
                {:else}
                    <svg class="mr-1.5 -ml-0.5 size-5 text-green-200" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                        <title>Start Recording</title>
                        <circle cx="10" cy="10" r="6" />
                    </svg>
                    Start Recording
                {/if}
            </button>
            </span>
        </div>
        </div>

        <div class="mt-8">
          <WeatherPanel />
          <AircraftPanel />
        </div>
{/if}
</div>


