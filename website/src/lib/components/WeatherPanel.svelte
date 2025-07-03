<script lang="ts">
import { environmentState } from '$lib/stores/environmentState';
import { derived } from 'svelte/store';

// Helper to format numbers with units
function formatValue(val: number | undefined | null, unit: string, digits = 2) {
  if (val === undefined || val === null || isNaN(val)) return '-';
  return `${val.toFixed(digits)}${unit}`;
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
</script>

<div>
  <h3 class="text-base font-semibold text-gray-900">Weather Information</h3>
  <dl class="mt-5 grid grid-cols-1 divide-y divide-gray-200 overflow-hidden rounded-lg bg-white shadow-sm md:grid-cols-4 md:divide-x md:divide-y-0">
    <!-- Sea Level Pressure -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Sea Level Pressure</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {formatValue($environmentState.sea_level_pressure, ' hPa')}
        </div>
      </dd>
    </div>
    <!-- Ambient Temperature -->
    <div class="px-4 py-5 sm:p-6">
      <dt class="text-base font-normal text-gray-900">Ambient Temperature</dt>
      <dd class="mt-1 flex items-baseline justify-between md:block lg:flex">
        <div class="flex items-baseline text-2xl font-semibold text-indigo-600">
          {formatValue($env.ambient_temperature, ' °C')}
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
