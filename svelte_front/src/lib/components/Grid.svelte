<script lang="ts">
	import { T } from '@threlte/core';
	import type { Grid } from '$lib/types';
	import { InstancedMesh, createTransition } from '@threlte/extras';
	import Cell from './Cell.svelte';
	import { Material } from 'three';
	import { cubicOut } from 'svelte/easing';

	export let grid: Grid;

	const cellFade = createTransition<Material>((ref) => {
		if (!ref.transparent) ref.transparent = true;
		return {
			tick(t) {
				ref.opacity = t;
			},
			easing: cubicOut,
			duration: 200
		};
	});
</script>

<T.Mesh position={[grid.width / 2, grid.height / 2, 0]}>
	<T.MeshBasicMaterial color={'black'} in={cellFade} out={cellFade} />
	<T.BoxGeometry args={[grid.width + 2, grid.height + 2, 0]} />
</T.Mesh>

{#if !grid.paused && grid.width != 0 && grid.height != 0}
	<InstancedMesh limit={grid.cells.length} range={grid.cells.length}>
		<T.MeshBasicMaterial in={cellFade} out={cellFade} />
		<T.BoxGeometry args={[1, 1, 1]} />
		{#each grid.cells as cell, idx}
			{@const y = idx / grid.width}
			{@const x = idx % grid.width}
			<Cell {x} {y} {cell} />
		{/each}
	</InstancedMesh>
{/if}
