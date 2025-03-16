import { writable } from "svelte/store";

export const isConfigOpen = writable(true);

const savedWidth = localStorage.getItem("sidebarWidth");
export const sidebarWidth = writable(savedWidth ? parseInt(savedWidth, 10) : 200);