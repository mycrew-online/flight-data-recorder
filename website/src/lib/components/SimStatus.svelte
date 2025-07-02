<script lang="ts">
import { simStatus } from '../stores/simStatus';
</script>

<span class="inline-flex items-center" title={ $simStatus ? 'SimConnect is connected to the simulator.' : 'SimConnect is not connected. No data will be received.' }>
  {#if $simStatus}
    <svg class="w-7 h-7 text-green-600" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
      <circle cx="10" cy="10" r="8" />
      <path d="M7 10.5l2 2 4-4" stroke="#fff" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" />
    </svg>
  {:else}
    <svg class="w-7 h-7 text-red-600" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
      <circle cx="10" cy="10" r="8" />
      <line x1="7" y1="7" x2="13" y2="13" stroke="#fff" stroke-width="2" stroke-linecap="round" />
      <line x1="13" y1="7" x2="7" y2="13" stroke="#fff" stroke-width="2" stroke-linecap="round" />
    </svg>
  {/if}
</span>
