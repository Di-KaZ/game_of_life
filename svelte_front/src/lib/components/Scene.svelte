<script lang="ts">
	import { T } from '@threlte/core';
	import type { Grid as GridObj } from '$lib/types';
	import { OrbitControls, interactivity, transitions } from '@threlte/extras';
	import Grid from './Grid.svelte';
	import { Grid as GridHelper } from '@threlte/extras';
	interactivity();
	transitions();

	const ws = new WebSocket('ws://127.0.0.1:8080/state');
	let grid: GridObj = { width: 0, height: 0, cells: [], paused: true };

	ws.onmessage = ({ data }) => (grid = JSON.parse(data));
</script>

<T.PerspectiveCamera
	makeDefault
	on:create={({ ref }) => {
		ref.position.set(0, 3, 10);
	}}
>
	<OrbitControls autoRotateSpeed={2} enableDamping target.y={2} />
</T.PerspectiveCamera>

<GridHelper />
<Grid {grid} />
