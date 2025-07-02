import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';

export interface SimulatorState {
  Sim: number;
  Pause: number;
  Crashed: number;
  View: number;
  AircraftLoaded: string;
  FlightLoaded: string;
  FlightPlan: string;
}

export const simulatorState = writable<SimulatorState | null>(null);

EventsOn('simulator::state', (state: SimulatorState) => {
  simulatorState.set(state);
});
