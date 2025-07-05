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
  heading_magnetic: number; // Added magnetic heading
  bank: number;
  alt_above_ground: number;
  pitch: number;
  vertical_speed: number; // Changed to number for consistency
  ground_velocity: number; // Changed to number for consistency
  airpeed_true: number; // Changed to number for consistency
  angle_of_attack: number; // Changed to number for consistency
  // Add other properties as needed from your backend
}

export const airplaneState = writable<AirplaneState>();

// Initialize with backend state
GetAirplaneState().then(airplaneState.set);

EventsOn('airplane::state', (state: AirplaneState) => {
  airplaneState.set(state);
});
