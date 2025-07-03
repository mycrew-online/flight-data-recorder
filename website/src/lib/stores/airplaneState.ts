import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetAirplaneState } from '$lib/wailsjs/go/internal/App';

export interface AirplaneState {
  title: string;
  latitude: number;
  longitude: number;
  altitude: number;
  heading: number;
  airspeed: number;
  // Add other properties as needed from your backend
}

export const airplaneState = writable<AirplaneState>();

// Initialize with backend state
GetAirplaneState().then(airplaneState.set);

EventsOn('airplane::state', (state: AirplaneState) => {
  airplaneState.set(state);
});
