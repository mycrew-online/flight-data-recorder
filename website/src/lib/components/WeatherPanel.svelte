<script lang="ts">
import { environmentState } from '$lib/stores/environmentState';
import { simulatorState } from '$lib/stores/simulatorState';
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

// Surface condition enum map
const surfaceConditionMap: Record<number, string> = {
  0: 'Normal',
  1: 'Wet',
  2: 'Icy',
  3: 'Snow',
};

// Surface type enum map
const surfaceTypeMap: Record<number, string> = {
  0: 'Concrete',
  1: 'Grass',
  2: 'Water',
  3: 'Grass (bumpy)',
  4: 'Asphalt',
  5: 'Short grass',
  6: 'Long grass',
  7: 'Hard turf',
  8: 'Snow',
  9: 'Ice',
  10: 'Urban',
  11: 'Forest',
  12: 'Dirt',
  13: 'Coral',
  14: 'Gravel',
  15: 'Oil treated',
  16: 'Steel mats',
  17: 'Bituminus',
  18: 'Brick',
  19: 'Macadam',
  20: 'Planks',
  21: 'Sand',
  22: 'Shale',
  23: 'Tarmac',
  24: 'Wright flyer track',
};

function formatSurfaceCondition(val: number | undefined | null): string {
  if (val === undefined || val === null || isNaN(val)) return '-';
  return surfaceConditionMap[val] ?? `Unknown (${val})`;
}

function formatSurfaceType(val: number | undefined | null): string {
  if (val === undefined || val === null || isNaN(val)) return '-';
  return surfaceTypeMap[val] ?? `Unknown (${val})`;
}

// Helper: Convert Zulu time (seconds since midnight) to local time (seconds since midnight), then format as HH:MM
function zuluToLocalTime(zuluSeconds: number | undefined | null, offset: number | undefined | null): string {
  if (typeof zuluSeconds !== 'number' || isNaN(zuluSeconds)) return '-';
  let localSeconds = zuluSeconds - (typeof offset === 'number' && !isNaN(offset) ? offset : 0);
  if (localSeconds < 0) localSeconds += 86400;
  if (localSeconds >= 86400) localSeconds -= 86400;
  const hours = Math.floor(localSeconds / 3600);
  const minutes = Math.floor((localSeconds % 3600) / 60);
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
}


// Derived stores for easier destructuring
const env = derived(environmentState, ($env) => $env ?? {});
const sim = derived(simulatorState, ($sim) => $sim ?? {});

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
    <!-- Surface Condition (from simulator state) -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Surface Condition</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {formatSurfaceCondition($sim.surface_condition)}
        </div>
      </dd>
    </div>
    <!-- Surface Type (from simulator state) -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Surface Type</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {formatSurfaceType($sim.surface_type)}
        </div>
      </dd>
    </div>
    <!-- Sunrise (local time) -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Sunrise</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          <svg class="mr-2 size-5 text-amber-300" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true"><title>Sunrise</title><path stroke-linecap="round" stroke-linejoin="round" d="M12 3v3m0 12v3m9-9h-3M6 12H3m15.364-6.364l-2.121 2.121M6.757 17.243l-2.121 2.121M17.243 17.243l2.121 2.121M6.757 6.757L4.636 4.636M12 7a5 5 0 1 1 0 10a5 5 0 0 1 0-10Z"/></svg>
          {zuluToLocalTime($env.zulu_sunrise_time, $env.time_zone_offset)}
        </div>
      </dd>
    </div>
    <!-- Sunset (local time) -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Sunset</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          <svg class="mr-2 size-5 text-orange-400" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true"><title>Sunset</title><path stroke-linecap="round" stroke-linejoin="round" d="M3 17h18M4.5 21h15M12 3v2m0 14v2m9-9h-2M5 12H3m15.364-6.364l-1.414 1.414M6.05 17.95l-1.414 1.414M17.95 17.95l1.414 1.414M6.05 6.05L4.636 4.636"/><circle cx="12" cy="14" r="5" fill="currentColor"/></svg>
          {zuluToLocalTime($env.zulu_sunset_time, $env.time_zone_offset)}
        </div>
      </dd>
    </div>
  </dl>
</div>
