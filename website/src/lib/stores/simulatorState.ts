import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetSimulatorState } from '$lib/wailsjs/go/internal/App';

export interface SimulatorState {
  sim: number;
  pause: number;
  crashed: number;
  view: number;
  aircraft_loaded: string;
  flight_loaded: string;
  flight_plan: string;
  simulation_rate: number;
  realism: number;
  surface_condition: number;
  surface_info_valid: boolean; // Changed to boolean
  surface_type: number;
  on_any_runway: number;
  in_parking_state: number;
  on_ground: boolean; // Changed to boolean
}

export const simulatorState = writable<SimulatorState>();

GetSimulatorState().then(simulatorState.set);

EventsOn('simulator::state', (state: SimulatorState) => {
  simulatorState.set(state);
});
