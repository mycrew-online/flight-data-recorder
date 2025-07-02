import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetAirplaneState } from '$lib/wailsjs/go/internal/App';

export const airplaneState = writable({});

// Initialize with backend state
GetAirplaneState().then(airplaneState.set);

EventsOn('airplane::state', (state) => {
  airplaneState.set(state);
});
