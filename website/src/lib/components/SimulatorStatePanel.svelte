<script lang="ts">
  import { simulatorState } from '$lib/stores/simulatorState';
</script>

{#if $simulatorState}
  <div class="p-4 bg-white rounded shadow mt-4">
    <h2 class="text-lg font-bold mb-2">Simulator State</h2>
    <ul class="text-sm">
      <li><b>Sim:</b> {$simulatorState.Sim}</li>
      <li><b>Pause:</b> {$simulatorState.Pause}</li>
      <li><b>Crashed:</b> {$simulatorState.Crashed}</li>
      <li><b>View:</b> {$simulatorState.View}</li>
      <li><b>AircraftLoaded:</b> {$simulatorState.AircraftLoaded}</li>
      <li><b>FlightLoaded:</b> {$simulatorState.FlightLoaded}</li>
      <li><b>FlightPlan:</b> {$simulatorState.FlightPlan}</li>
    </ul>
  </div>
{:else}
  <div class="p-4 bg-white rounded shadow mt-4 text-gray-400">No simulator state available.</div>
{/if}
