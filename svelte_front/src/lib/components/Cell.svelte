<script lang="ts">
	import type { Cell } from '$lib/types';
	import { Color, type Mesh } from 'three';
	import { Instance, createTransition } from '@threlte/extras';
	import { cubicOut } from 'svelte/easing';

	export let cell: Cell;
	export let x: number;
	export let y: number;

	const colors = ['#34667A', '#49798A', '#D69700'];

	const pop = (delay: number) =>
		createTransition<Mesh>((ref) => {
			return {
				tick(t) {
					ref.scale.setScalar(t);
				},
				easing: cubicOut,
				duration: 500,
				delay
			};
		});
</script>

{#if cell.alive}
	<Instance
		color={colors[cell.turns] ?? colors[colors.length - 1]}
		position.y={y}
		position.x={x}
		in={pop(200)}
		out={pop(0)}
	/>
{/if}
