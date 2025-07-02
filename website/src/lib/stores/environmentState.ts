import { writable } from 'svelte/store';
import { EventsOn } from '$lib/wailsjs/runtime/runtime';
import { GetEnvironmentState } from '$lib/wailsjs/go/internal/App';

export const environmentState = writable({});

// Initialize with backend state
GetEnvironmentState().then(environmentState.set);

EventsOn('environment::state', (state) => {
  environmentState.set(state);
});
