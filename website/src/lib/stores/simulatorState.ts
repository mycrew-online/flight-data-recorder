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
}

export const simulatorState = writable<SimulatorState>();

GetSimulatorState().then(simulatorState.set);

EventsOn('simulator::state', (state: SimulatorState) => {
  simulatorState.set(state);
});
