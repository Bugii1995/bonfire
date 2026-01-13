<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Navbar from "$lib/components/Navbar.svelte";
	import Footer from '$lib/components/Footer.svelte';
	let { children } = $props();
</script>


<Navbar />

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}

<Footer />


