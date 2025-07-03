<script lang="ts">
import { airplaneState } from '$lib/stores/airplaneState';
import { BrowserOpenURL } from '$lib/wailsjs/runtime/runtime';

// Helper to split a number into integer and decimal parts as strings
function splitDecimal(val: number, digits = 5): [string, string] {
  const [intPart, decPart] = val.toFixed(digits).split('.')
  return [intPart, decPart];
}

function openInGoogleMaps(lat: number, long: number) {
  const url = `https://maps.google.com/?q=${lat},${long}`;
  BrowserOpenURL(url);
}
</script>

<dl class="mt-5 grid grid-cols-1 divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow-sm md:grid-cols-4 md:divide-x md:divide-y-0">
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Latitude / Longitude</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.latitude === 'number' && typeof $airplaneState?.longitude === 'number'}
          <button class="cursor-pointer focus:outline-none text-indigo-600" title="Open in Google Maps" on:click={() => openInGoogleMaps($airplaneState.latitude, $airplaneState.longitude)}>
            {#await Promise.resolve(splitDecimal($airplaneState.latitude)) then latParts}
              {#await Promise.resolve(splitDecimal($airplaneState.longitude)) then lonParts}
                <span>{latParts[0]}<span class="relative text-[.55em]" style="top:0">.{latParts[1]}</span></span>
                <span class="mx-1">/</span>
                <span>{lonParts[0]}<span class="relative text-[.55em]" style="top:0">.{lonParts[1]}</span></span>
              {/await}
            {/await}
          </button>
        {:else}
          <span>-/-</span>
        {/if}
      </div>
    </dd>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Altitude</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.altitude === 'number'}
          <span>{$airplaneState.altitude.toLocaleString(undefined, {maximumFractionDigits: 0})}<span class="ml-1">ft</span></span>
          <span class="ml-2 text-sm font-medium text-gray-500">
            {#if $airplaneState.altitude * 0.3048 >= 1000}
              {($airplaneState.altitude * 0.3048 / 1000).toFixed(2)} km
            {:else}
              {($airplaneState.altitude * 0.3048).toFixed(0)} m
            {/if}
          </span>
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Heading</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.heading === 'number'}
          <span>{$airplaneState.heading.toFixed(0)}</span>
          <span class="ml-1">°</span>
          {#if typeof $airplaneState.heading_magnetic === 'number'}
            <span class="ml-2 text-sm font-medium text-gray-500">{$airplaneState.heading_magnetic.toFixed(0)}° MAG</span>
          {/if}
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Airspeed</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.airspeed === 'number'}
          <span>{$airplaneState.airspeed.toFixed(0)}<span class="ml-1">kt</span></span>
          <span class="ml-2 text-sm font-medium text-gray-500">{($airplaneState.airspeed * 1.852).toFixed(0)} km/h</span>
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
</dl>

<dl class="mt-5 grid grid-cols-1 divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow-sm md:grid-cols-4 md:divide-x md:divide-y-0">
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Bank</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.bank === 'number'}
          <span>{$airplaneState.bank.toFixed(1)}</span>
          <span class="ml-1">°</span>
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Pitch</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.pitch === 'number'}
          <span>{$airplaneState.pitch.toFixed(1)}</span>
          <span class="ml-1">°</span>
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Vertical Speed</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.vertical_speed === 'number'}
          <span>{$airplaneState.vertical_speed.toFixed(0)}<span class="ml-1">ft/min</span></span>
          <span class="ml-2 text-sm font-medium text-gray-500">{($airplaneState.vertical_speed * 0.3048).toFixed(1)} m/min</span>
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
  <div class="px-4 py-5 sm:p-6">
    <dt class="text-base font-normal text-gray-900">Alt Above Ground</dt>
    <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
      <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
        {#if typeof $airplaneState?.alt_above_ground === 'number'}
          <span>{$airplaneState.alt_above_ground.toLocaleString(undefined, {maximumFractionDigits: 0})}<span class="ml-1">ft</span></span>
          <span class="ml-2 text-sm font-medium text-gray-500">
            {#if $airplaneState.alt_above_ground * 0.3048 >= 1000}
              {($airplaneState.alt_above_ground * 0.3048 / 1000).toFixed(2)} km
            {:else}
              {($airplaneState.alt_above_ground * 0.3048).toFixed(0)} m
            {/if}
          </span>
        {:else}
          -
        {/if}
      </div>
    </dd>
  </div>
</dl>
