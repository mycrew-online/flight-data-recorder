import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetSimulatorState } from '$lib/wailsjs/go/internal/App';

export interface SimulatorState {
  Sim: number;
  Pause: number;
  Crashed: number;
  View: number;
  AircraftLoaded: string;
  FlightLoaded: string;
  FlightPlan: string;
}

export const simulatorState = writable<SimulatorState>();

GetSimulatorState().then(simulatorState.set);

EventsOn('simulator::state', (state: SimulatorState) => {
  simulatorState.set(state);
});
