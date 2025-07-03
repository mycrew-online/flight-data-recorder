<script lang="ts">
import { environmentState } from '$lib/stores/environmentState';
import { derived } from 'svelte/store';

// Helper to format numbers with units
function formatValue(val: number | undefined | null, unit: string, digits = 2) {
  if (val === undefined || val === null || isNaN(val)) return '-';
  return `${val.toFixed(digits)}${unit}`;
}

// Helper to split a number into integer and decimal parts as strings
function splitDecimal(val: number, digits = 2): [string, string] {
  const [intPart, decPart] = val.toFixed(digits).split('.')
  return [intPart, decPart];
}

// Helper for integer only (no decimals)
function splitInt(val: number): [string, string] {
  return [val.toFixed(0), ''];
}

// Convert inHg to hPa (1 inHg = 33.8639 hPa)
function inHgTohPa(inHg: number): number {
  return inHg * 33.8639;
}

// Convert C to F
function cToF(c: number): number {
  return c * 9/5 + 32;
}

// Helper to format visibility
function formatVisibility(val: number | undefined | null): string {
  if (val === undefined || val === null || isNaN(val)) return '-';
  if (val > 1000) {
    return `${(val / 1000).toFixed(2)} km`;
  }
  return `${val.toFixed(0)} m`;
}

// Derived store for easier destructuring
const env = derived(environmentState, ($env) => $env ?? {});

$: seaLevelPressureInHg = $environmentState?.sea_level_pressure;
$: seaLevelPressureInHgParts = typeof seaLevelPressureInHg === 'number' ? splitDecimal(seaLevelPressureInHg, 2) : null;
$: seaLevelPressure = typeof seaLevelPressureInHg === 'number' ? inHgTohPa(seaLevelPressureInHg) : null;
$: seaLevelPressureParts = typeof seaLevelPressure === 'number' ? [seaLevelPressure.toFixed(0), ''] : null;

$: ambientTemp = $env.ambient_temperature;
$: ambientTempParts = typeof ambientTemp === 'number' ? splitInt(ambientTemp) : null;
$: ambientTempF = typeof ambientTemp === 'number' ? cToF(ambientTemp) : null;
$: ambientTempFParts = typeof ambientTempF === 'number' ? splitInt(ambientTempF) : null;
</script>

<div>
  <dl class="mt-5 grid grid-cols-1 divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow-sm md:grid-cols-4 md:divide-x md:divide-y-0">
    <!-- Sea Level Pressure -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Sea Level Pressure</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {#if seaLevelPressureParts}
            <span>{seaLevelPressureParts[0]}<span class="ml-1">hPa</span></span>
            {#if seaLevelPressureInHgParts}
              <span class="ml-2 text-sm font-medium text-gray-500">{seaLevelPressureInHgParts[0]}.{seaLevelPressureInHgParts[1]} inHg</span>
            {/if}
          {:else}
            -
          {/if}
        </div>
      </dd>
    </div>
    <!-- Ambient Temperature -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Ambient Temperature</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {#if ambientTempParts}
            <span>{ambientTempParts[0]}<span class="ml-1">°C</span></span>
            {#if ambientTempFParts}
              <span class="ml-2 text-sm font-medium text-gray-500">{ambientTempFParts[0]} °F</span>
            {/if}
          {:else}
            -
          {/if}
        </div>
      </dd>
    </div>
    <!-- Wind -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Wind</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {formatValue($env.ambient_wind_velocity, ' kt', 1)}
          <span class="ml-2 text-sm font-medium text-gray-500">{formatValue($env.ambient_wind_direction, '°', 0)}</span>
        </div>
      </dd>
    </div>
    <!-- Visibility -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Visibility</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {formatVisibility($env.ambient_visibility)}
        </div>
      </dd>
    </div>
  </dl>
</div>
