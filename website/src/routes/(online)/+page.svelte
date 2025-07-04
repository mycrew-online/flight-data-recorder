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

// Helper: Convert seconds since midnight to HH:MM (no seconds)
function formatTimeHHMM(seconds: number | undefined | null): string {
  if (typeof seconds !== 'number' || isNaN(seconds)) return '';
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
}

// Helper: Convert Zulu time (seconds since midnight) to local time (seconds since midnight), then format as HH:MM
function zuluToLocalTime(zuluSeconds: number | undefined | null, offset: number | undefined | null): string {
  if (typeof zuluSeconds !== 'number' || isNaN(zuluSeconds) || typeof offset !== 'number' || isNaN(offset)) return '';
  // offset is in seconds, positive east of GMT, negative west
  let localSeconds = zuluSeconds - offset;
  // wrap around 0-86399
  if (localSeconds < 0) localSeconds += 86400;
  if (localSeconds >= 86400) localSeconds -= 86400;
  return formatTimeHHMM(localSeconds);
}

// Helper: Format local date as DD.MM.YYYY
function formatLocalDate(year: number | undefined, month: number | undefined, day: number | undefined): string {
  if (!year || !month || !day) return '';
  return `${day.toString().padStart(2, '0')}.${month
    .toString()
    .padStart(2, '0')}.${year.toString().padStart(4, '0')}`;
}
</script>
<div class="lg:flex lg:items-start lg:justify-between">
  <div class="min-w-0 flex-1">
    <h2 class="text-2xl/7 font-bold text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">{$airplaneState?.title}</h2>
    <div class="mt-1 flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap sm:space-x-6">
    <!-- Combined Sim Time, Date, and Local Time -->
    <div class="mt-2 flex items-center text-sm text-gray-500">
      <!--<svg class="mr-1.5 size-5 shrink-0 text-gray-400" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <title>Simulated time (in-game)</title>
        <path fill-rule="evenodd" d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm0-14.5A6.5 6.5 0 1 1 3.5 10 6.508 6.508 0 0 1 10 3.5Zm.75 3.75a.75.75 0 0 0-1.5 0v3.25a.75.75 0 0 0 .22.53l2.25 2.25a.75.75 0 1 0 1.06-1.06l-2.03-2.03V7.25Z" clip-rule="evenodd" />
      </svg>-->
      <svg class="mr-1.5 size-5 shrink-0 text-gray-400" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <title>Local date (simulator)</title>
        <path fill-rule="evenodd" d="M6 2a1 1 0 0 0-1 1v1H5A3 3 0 0 0 2 7v7a3 3 0 0 0 3 3h10a3 3 0 0 0 3-3V7a3 3 0 0 0-3-3h-.002V3a1 1 0 1 0-2 0v1H7V3a1 1 0 0 0-1-1Zm10 5v1H4V7a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1Zm-1 3v4a1 1 0 0 1-1 1H6a1 1 0 0 1-1-1v-4h10Z" clip-rule="evenodd" />
      </svg>
      {formatLocalDate($environmentState?.local_year, $environmentState?.local_month, $environmentState?.local_day)}
      <span>&nbsp;{formatLocalTime($environmentState?.local_time)}</span>
      <svg class="ml-4 mr-1.5 size-5 shrink-0 text-gray-400" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <title>Timezone Offset</title>
        <path fill-rule="evenodd" d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm0-14.5A6.5 6.5 0 1 1 3.5 10 6.508 6.508 0 0 1 10 3.5Zm.75 3.75a.75.75 0 0 0-1.5 0v3.25a.75.75 0 0 0 .22.53l2.25 2.25a.75.75 0 1 0 1.06-1.06l-2.03-2.03V7.25Z" clip-rule="evenodd" />
      </svg>
      <span>GMT{($environmentState?.time_zone_offset != null ? ((Math.round(-$environmentState.time_zone_offset / 3600) >= 0 ? '+' : '') + Math.round(-$environmentState.time_zone_offset / 3600)) : '')}</span>
    </div>
    <!-- Timezone and Zulu Sun Times -->
      <div class="mt-2 flex items-center text-sm text-gray-500">
        {#if $environmentState && typeof $environmentState.time_of_day === 'number'}
          <span class="flex items-center justify-end mr-1.5 mt-[-1.5px]" title={(() => {
            switch ($environmentState.time_of_day) {
              case 0: return 'Dawn';
              case 1: return 'Day';
              case 2: return 'Dusk';
              case 3: return 'Night';
              default: return 'Unknown';
            }
          })()}>
            {#if $environmentState.time_of_day === 0}
              <!-- Dawn Icon -->
              <svg class="size-5 text-amber-300" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><title>Dawn</title><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v3m0 12v3m9-9h-3M6 12H3m15.364-6.364l-2.121 2.121M6.757 17.243l-2.121 2.121M17.243 17.243l2.121 2.121M6.757 6.757L4.636 4.636M12 7a5 5 0 1 1 0 10a5 5 0 0 1 0-10Z"/></svg>
            {:else if $environmentState.time_of_day === 1}
              <!-- Day Icon -->
              <svg class="size-5 text-yellow-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><title>Day</title><circle cx="12" cy="12" r="5" fill="currentColor"/><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2m0 14v2m9-9h-2M5 12H3m15.364-6.364l-1.414 1.414M6.05 17.95l-1.414 1.414M17.95 17.95l1.414 1.414M6.05 6.05L4.636 4.636"/></svg>
            {:else if $environmentState.time_of_day === 2}
              <!-- Dusk Icon -->
              <svg class="size-5 text-orange-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><title>Dusk</title><path stroke-linecap="round" stroke-linejoin="round" d="M3 17h18M4.5 21h15M12 3v2m0 14v2m9-9h-2M5 12H3m15.364-6.364l-1.414 1.414M6.05 17.95l-1.414 1.414M17.95 17.95l1.414 1.414M6.05 6.05L4.636 4.636"/><circle cx="12" cy="14" r="5" fill="currentColor"/></svg>
            {:else if $environmentState.time_of_day === 3}
              <!-- Night Icon -->
              <svg class="size-5 text-blue-900" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><title>Night</title><path stroke-linecap="round" stroke-linejoin="round" d="M21 12.79A9 9 0 1 1 11.21 3a7 7 0 1 0 9.79 9.79Z"/></svg>
            {:else}
              <!-- Unknown Icon -->
              <svg class="size-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><title>Unknown</title><circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5" fill="none"/><text x="12" y="16" text-anchor="middle" font-size="8" fill="currentColor">?</text></svg>
            {/if}
          </span>
        {/if}
        {formatSimTime($environmentState?.sim_time)}
        <!-- Removed sunrise/sunset from subtitle -->
  <!-- Add Zulu Sunrise/Sunset to AircraftPanel -->
      </div>
    </div>
  </div>
  <div class="mt-5 flex lg:mt-0 lg:ml-4">
      <span>
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

<div class="relative border-b border-gray-200 pb-5 sm:pb-0">
  <div class="mt-12">
    <div class="grid grid-cols-1 sm:hidden">
      <!-- Use an "onChange" listener to redirect the user to the selected tab URL. -->
      <select aria-label="Select a tab" class="col-start-1 row-start-1 w-full appearance-none rounded-md bg-white py-2 pr-8 pl-3 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600">
        <option>Applied</option>
        <option>Phone Screening</option>
        <option selected>Interview</option>
        <option>Offer</option>
        <option>Hired</option>
      </select>
      <svg class="pointer-events-none col-start-1 row-start-1 mr-2 size-5 self-center justify-self-end fill-gray-500" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true" data-slot="icon">
        <path fill-rule="evenodd" d="M4.22 6.22a.75.75 0 0 1 1.06 0L8 8.94l2.72-2.72a.75.75 0 1 1 1.06 1.06l-3.25 3.25a.75.75 0 0 1-1.06 0L4.22 7.28a.75.75 0 0 1 0-1.06Z" clip-rule="evenodd" />
      </svg>
    </div>
    <!-- Tabs at small breakpoint and up -->
    <div class="hidden sm:block">
      <nav class="-mb-px flex space-x-8">
        <!-- Current: "border-indigo-500 text-indigo-600", Default: "border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700" -->
        <a href="#" class="border-b-2 border-indigo-500 px-1 pb-4 text-sm font-medium whitespace-nowrap text-indigo-600" aria-current="page">General</a>
        <!-- border-b-2 border-indigo-500 px-1 pb-4 text-sm font-medium whitespace-nowrap text-indigo-600 -->
        <!-- border-b-2 border-transparent px-1 pb-4 text-sm font-medium whitespace-nowrap text-gray-500 hover:border-gray-300 hover:text-gray-700 -->
      </nav>
    </div>
  </div>
</div>


<div class="mt-8">
  <WeatherPanel />
  <AircraftPanel />
</div>


