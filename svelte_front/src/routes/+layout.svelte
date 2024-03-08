<script lang="ts">
	import { AppShell, Avatar } from '@skeletonlabs/skeleton';
	import '../app.postcss';
	import { Canvas } from '@threlte/core';
	import Scene from '$lib/components/Scene.svelte';
	import { onMount } from 'svelte';

	let currentMessage: string;

	let messages: { name: string; content: string }[] = [];

	let ws: WebSocket;

	onMount(() => {
		ws = new WebSocket('ws://127.0.0.1:8080/chat');

		ws.onmessage = ({ data }) => (messages = [...messages, JSON.parse(data)]);
	});

	const onSend = () => {
		ws.send(currentMessage);
		currentMessage = '';
	};
</script>

<AppShell>
	<svelte:fragment slot="header">Header</svelte:fragment>
	<svelte:fragment slot="sidebarRight">
		<div class="h-full overflow-y-hidden">
			<section class="p-4 overflow-y-auto space-y-4">
				{#each messages as item}
					<div
						class={`grid ${
							item.name == '[SERVER]' ? 'flex flex-row-reverse' : 'flex flex-row'
						} gap-2`}
					>
						<Avatar src="" width="w-12" />
						<div class={`card p-4 variant-soft rounded-tl-none space-y-2`}>
							<header class="flex justify-between items-center">
								<p class="font-bold">{item.name}</p>
							</header>
							<p>{item.content}</p>
						</div>
					</div>
				{/each}
			</section>
			<div
				class="input-group input-group-divider grid-cols-[auto_1fr_auto] rounded-container-token self-end"
			>
				<button class="input-group-shim">+</button>
				<textarea
					bind:value={currentMessage}
					class="bg-transparent border-0 ring-0"
					name="prompt"
					id="prompt"
					placeholder="Write a message..."
					rows="1"
				/>
				<button class="variant-filled-primary" on:click={onSend}>Send</button>
			</div>
		</div>
	</svelte:fragment>
	<!-- (pageHeader) -->
	<!-- Router Slot -->
	<slot />
	<!-- ---- / ---- -->
	<!-- (pageFooter) -->
	<!-- (footer) -->
</AppShell>
<!-- <AppShell> -->
<!-- 	<svelte:fragment slot="header">Header</svelte:fragment> -->
<!-- 	<slot /> -->
<!-- 	<svelte:fragment slot="sidebarRight"> -->
<!-- 			{#each messages as item} -->
<!-- 				<div class="grid grid-cols-[auto_1fr] gap-2"> -->
<!-- 					<Avatar src="" width="w-12" /> -->
<!-- 					<div class="card p-4 variant-soft rounded-tl-none space-y-2"> -->
<!-- 						<header class="flex justify-between items-center"> -->
<!-- 							<p class="font-bold">{item.name}</p> -->
<!-- 						</header> -->
<!-- 						<p>{item.content}</p> -->
<!-- 					</div> -->
<!-- 				</div> -->
<!-- 			{/each} -->
<!-- 	</svelte:fragment> -->
<!-- </AppShell> -->
