<script lang="ts">
  import { page } from "$app/state";
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from "flowbite-svelte";
  let activeUrl = $derived(page.url.pathname);
  let activeClass = "text-white bg-green-700 md:bg-transparent md:text-green-700 md:dark:text-white dark:bg-green-600 md:dark:bg-transparent";
  let nonActiveClass =
    "text-gray-700 hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-green-700 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent";
</script>

<Navbar>
  <NavBrand href="/">
    <span class="self-center text-xl font-semibold whitespace-nowrap dark:text-white">🍫 Chocolingo</span>
  </NavBrand>
  <NavHamburger />
  <NavUl {activeUrl} classes={{ active: activeClass, nonActive: nonActiveClass }}>
    <NavLi href="/">Home</NavLi>
    <NavLi href="/contact">Study</NavLi>
    <NavLi href="/dashboard">Dashboard</NavLi>
    <NavLi href="/leaderboard">Leaderboard</NavLi>
    <NavLi href="/account">Account</NavLi>
  </NavUl>
</Navbar>
